package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	restproto "github.com/savageking-io/ogbrest/proto"
	"github.com/savageking-io/ogbrest/restlib"
	"github.com/savageking-io/ogbuser/db"
	"github.com/savageking-io/ogbuser/group"
	"github.com/savageking-io/ogbuser/proto"
	"github.com/savageking-io/ogbuser/user"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"time"
)

type Service struct {
	rest   *restlib.RestInterServiceServer
	config *ServiceConfig
	db     *db.Database
	groups *group.GroupsData
	users  *user.UsersData

	proto.UnimplementedUserServiceServer
}

func NewService(config *ServiceConfig) *Service {
	return &Service{
		config: config,
		users:  user.NewUsersData(nil),
		groups: group.NewGroupsData(),
	}
}

func (s *Service) Init() error {
	log.Infof("Initializing service")

	if err := s.ConnectToDatabase(); err != nil {
		return err
	}

	s.users.SetDb(s.db)

	if err := s.InitGroups(); err != nil {
		log.Errorf("Failed to initialize groups: %v", err)
		return err
	}

	if err := s.InitializeRest(s.config.Rest); err != nil {
		log.Errorf("Failed to initialize REST server: %v", err)
		return err
	}

	return nil
}

func (s *Service) ConnectToDatabase() error {
	s.db = new(db.Database)
	if err := s.db.Init(&s.config.Postgres); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := s.db.Connect(); err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Infof("Connected to database %s", s.config.Postgres.Database)
	if err := s.db.TryPopulate(); err != nil {
		log.Errorf("Failed to populate database: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) InitGroups() error {
	if s.db == nil {
		return fmt.Errorf("postgres db is not initialized")
	}

	log.Infof("Initializing groups")

	groups, err := s.db.LoadGroups(context.Background())
	if err != nil {
		return fmt.Errorf("failed to load groups: %w", err)
	}

	log.Infof("Found %d groups", len(groups))

	for _, userGroup := range groups {
		g := group.NewGroupFromSchema(s.db, &userGroup)
		if err := g.Init(context.Background()); err != nil {
			log.Errorf("Failed to initialize group %d: %v", userGroup.Id, err)
			continue
		}
		s.groups.Add(g)
	}

	log.Infof("Groups initialized. Total: %d", len(s.groups.GetAll()))

	return nil
}

func (s *Service) InitializeRest(config restlib.RestInterServiceConfig) error {
	log.Infof("Initializing REST service")
	s.rest = restlib.NewRestInterServiceServer(config)
	if err := s.rest.Init(); err != nil {
		return err
	}

	log.Infof("Registering handlers")
	if err := s.rest.RegisterHandler("/auth/credentials", "POST", s.HandleAuthCredentialsRequest, true); err != nil {
		log.Warnf("Failed to register handler for /auth/credentials: %v", err)
	}
	if err := s.rest.RegisterHandler("/auth/platform", "POST", s.HandleAuthPlatformRequest, true); err != nil {
		log.Warnf("Failed to register handler for /auth/platform: %v", err)
	}
	if err := s.rest.RegisterHandler("/auth/server", "POST", s.HandleAuthServerRequest, true); err != nil {
		log.Warnf("Failed to register handler for /auth/server: %v", err)
	}
	if err := s.rest.RegisterHandler("/token", "POST", s.HandleVerifyTokenRequest, false); err != nil {
		log.Warnf("Failed to register handler for /token: %v", err)
	}

	for _, key := range s.rest.GetRegisteredHandlerKeys() {
		log.Infof("Registered handler: %s", key)
	}

	return nil
}

func (s *Service) Start() error {
	if s.rest == nil {
		log.Errorf("REST server is not initialized")
		return nil
	}

	log.Infof("Starting service")

	restErrChan := make(chan error, 1)
	rpcErrChan := make(chan error, 1)

	go func() {
		log.Infof("Starting REST server: %s:%d", AppConfig.Rest.Hostname, AppConfig.Rest.Port)
		err := s.rest.Start()
		restErrChan <- err
	}()

	go func() {
		log.Infof("Starting user service")

		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", AppConfig.Rpc.Hostname, AppConfig.Rpc.Port))
		if err != nil {
			log.Errorf("failed to listen: %v", err)
			rpcErrChan <- err
			return
		}
		grpcServer := grpc.NewServer()
		proto.RegisterUserServiceServer(grpcServer, s)
		if err := grpcServer.Serve(lis); err != nil {
			rpcErrChan <- err
			return
		}
		return
	}()

	select {
	case err := <-restErrChan:
		if err != nil {
			log.Errorf("REST channel reported error: %v", err)
		}
	}

	return nil
}

// HandleAuthCredentialsRequest will authenticate users using username and password
// Username might be actual username or email. This method handles both cases
func (s *Service) HandleAuthCredentialsRequest(ctx context.Context, in *restproto.RestApiRequest) (*restproto.RestApiResponse, error) {
	log.Tracef("HandleAuthCredentialsRequest")
	log.Tracef("Request: %+v", in)

	credentials := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal([]byte(in.Body), &credentials)
	if err != nil {
		return &restproto.RestApiResponse{
			Code:     12000,
			HttpCode: 400,
			Error:    "failed to parse request",
		}, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	username := ""
	if credentials.Username == "" && credentials.Email != "" {
		username = credentials.Email
	} else if credentials.Username != "" && credentials.Email == "" {
		username = credentials.Username
	} else {
		return &restproto.RestApiResponse{
			Code:     12001,
			HttpCode: 400,
			Error:    "empty username",
		}, fmt.Errorf("username and email are both empty")
	}

	if credentials.Password == "" {
		return &restproto.RestApiResponse{
			Code:     12002,
			HttpCode: 400,
			Error:    "empty password",
		}, fmt.Errorf("empty password")
	}

	u := user.NewUser(s.db, nil)
	if err := u.LoadByUsername(ctx, username); err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return &restproto.RestApiResponse{
				Code:     12003,
				HttpCode: 404,
				Error:    err.Error(),
			}, fmt.Errorf("user not found")
		}
		return &restproto.RestApiResponse{
			HttpCode: 401,
			Code:     12004,
			Error:    err.Error(),
		}, fmt.Errorf("failed to load user: %w", err)
	}

	groupIds, err := u.LoadGroups(ctx)
	if err != nil {
		log.Errorf("failed to load groups: %v", err)
		return &restproto.RestApiResponse{
			HttpCode: 401,
			Code:     12007,
			Error:    err.Error(),
		}, fmt.Errorf("failed to load user: %w", err)
	}

	log.Infof("Loaded %d groups", len(groupIds))
	for _, groupId := range groupIds {
		userGroup, exists := s.groups.Get(groupId)
		if !exists {
			log.Errorf("Service::HandleAuthCredentialsRequest: Group %d not found for user %d", groupId, u.GetId())
			continue
		}
		u.AddGroup(userGroup)
	}

	ok, err := VerifyPassword(credentials.Password, u.GetPassword())
	if err != nil {
		return &restproto.RestApiResponse{
			HttpCode: 500,
			Code:     12005,
			Error:    err.Error(),
		}, fmt.Errorf("failed to verify password")
	}

	if !ok {
		return &restproto.RestApiResponse{
			Code:     12003,
			HttpCode: 404,
			Error:    "user not found",
		}, fmt.Errorf("user not found")
	}

	// We keep users cached until they log out or we didn't receive anything from them for a long period of time
	// @TODO: Handle timeout
	// @TODO: Handle cleanup of duplicates
	log.Debugf("User [%s] authenticated and cached", u.GetUsername())
	if err := s.users.Add(u); err != nil {
		log.Errorf("failed to add user to cache: %v", err)
		return &restproto.RestApiResponse{
			Code:     12009,
			HttpCode: 500,
			Error:    "user load error",
		}, fmt.Errorf("user load error")
	}

	// @TOOO: Properly determine user platform
	session, err := u.InitializeSession(ctx, "web")
	if err != nil {
		return &restproto.RestApiResponse{
			Code:     12006,
			HttpCode: 500,
			Error:    err.Error(),
		}, fmt.Errorf("failed to initialize session: %w", err)
	}

	return &restproto.RestApiResponse{
		Code:     0,
		HttpCode: 200,
		Body:     fmt.Sprintf(`{"id": %d, "username": "%s", "email": "%s", "token": "%s"}`, u.GetId(), u.GetUsername(), u.GetEmail(), session.Token),
	}, nil
}

func (s *Service) HandleAuthPlatformRequest(ctx context.Context, in *restproto.RestApiRequest) (*restproto.RestApiResponse, error) {
	log.Tracef("HandleAuthPlatformRequest")
	return nil, nil
}

func (s *Service) HandleAuthServerRequest(ctx context.Context, in *restproto.RestApiRequest) (*restproto.RestApiResponse, error) {
	log.Tracef("HandleAuthServerRequest")
	return nil, nil
}

func (s *Service) HandleVerifyTokenRequest(ctx context.Context, in *restproto.RestApiRequest) (*restproto.RestApiResponse, error) {
	log.Tracef("HandleVerifyTokenRequest")

	return &restproto.RestApiResponse{
		Code:     0,
		HttpCode: 200,
	}, nil
}

// UserService
func (s *Service) Ping(ctx context.Context, in *proto.PingMessage) (*proto.PingMessage, error) {
	in.RepliedAt = timestamppb.New(time.Now())
	return in, nil
}

func (s *Service) ValidateToken(ctx context.Context, in *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	log.Tracef("ValidateToken")
	log.Tracef("Request: %+v", in)

	if s.db == nil {
		return &proto.ValidateTokenResponse{
			Code:  2,
			Error: "database is not initialized",
		}, nil
	}

	session, err := s.db.GetUserSessionByToken(ctx, in.Token)
	if err != nil {
		if errors.Is(err, db.ErrSessionNotFound) {
			return &proto.ValidateTokenResponse{
				Code:    0,
				IsValid: false,
				UserId:  0,
			}, nil
		}
		return &proto.ValidateTokenResponse{
			Code:  1,
			Error: err.Error(),
		}, err
	}

	if session == nil {
		return &proto.ValidateTokenResponse{
			Code:    0,
			IsValid: false,
			UserId:  0,
		}, nil
	}

	if time.Since(session.CreatedAt) > time.Duration(AppConfig.Crypto.JWT.Expiry)*time.Minute {
		log.Debugf("Token expired: %s. Duration: %+v", in.Token, time.Duration(AppConfig.Crypto.JWT.Expiry))
		return &proto.ValidateTokenResponse{
			Code:    0,
			IsValid: false,
			UserId:  int32(session.UserId),
		}, nil
	}

	log.Debugf("Token is valid: %s", in.Token)
	return &proto.ValidateTokenResponse{
		Code:    0,
		IsValid: true,
		UserId:  int32(session.UserId),
	}, nil
}

// HasPermission will check if a specified user has specific permission
func (s *Service) HasPermission(ctx context.Context, in *proto.HasPermissionRequest) (*proto.HasPermissionResponse, error) {
	log.Tracef("HasPermission")

	if in.UserId == 0 {
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, fmt.Errorf("invalid user id")
	}

	u, err := s.users.GetById(in.UserId)
	if err != nil {
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, fmt.Errorf("user not found")
	}

	if u == nil {
		log.Errorf("Service::HasPermission: received nil user struct without an error")
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, fmt.Errorf("user not found")
	}

	permission, err := u.HasPermission(ctx, in.Permission, in.Domain)
	if err != nil {
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, err
	}

	if permission == nil {
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, fmt.Errorf("received nil permission without an error")
	}

	return &proto.HasPermissionResponse{
		Read:   permission.Read,
		Write:  permission.Write,
		Delete: permission.Delete,
	}, nil
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
