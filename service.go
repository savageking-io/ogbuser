package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/savageking-io/ogbcommon/dbinit"
	restproto "github.com/savageking-io/ogbrest/proto"
	"github.com/savageking-io/ogbrest/restlib"
	"github.com/savageking-io/ogbuser/group"
	"github.com/savageking-io/ogbuser/proto"
	"github.com/savageking-io/ogbuser/schema"
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
	db     *sqlx.DB
	groups *group.GroupsData
	users  *user.UsersData

	proto.UnimplementedUserServiceServer
}

func NewService(config *ServiceConfig) *Service {
	return &Service{
		config: config,
		users:  user.NewUsersData(),
		groups: group.NewGroupsData(),
	}
}

func (s *Service) Init() error {
	log.Infof("Initializing service")

	return s.ConnectToDatabase()
}

func (s *Service) ConnectToDatabase() error {
	sslMode := "disable"
	if s.config.Postgres.SslMode {
		sslMode = "require"
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.config.Postgres.Hostname,
		s.config.Postgres.Port,
		s.config.Postgres.Username,
		s.config.Postgres.Password,
		s.config.Postgres.Database,
		sslMode)

	log.Debugf("Connecting to database: %s", connStr)

	var err error
	s.db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	s.db.SetMaxOpenConns(s.config.Postgres.MaxOpenConns)
	s.db.SetMaxIdleConns(s.config.Postgres.MaxIdleConns)
	s.db.SetConnMaxLifetime(time.Duration(s.config.Postgres.ConnMaxLifetime) * time.Minute)

	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Infof("Connected to database %s", s.config.Postgres.Database)

	// Populate db if it's empty
	err = dbinit.VerifyTables([]string{"users", "platforms", "groups"}, s.db)
	if err != nil {
		log.Warnf("Table verification failed: %s", err.Error())
		if s.config.Postgres.SourceFile == "" {
			log.Errorf("No source file provided. Cannot populate database")
			return err
		}
		if err := dbinit.Populate(s.db, s.config.Postgres.SourceFile); err != nil {
			log.Errorf("Failed to populate database: %s", err.Error())
			return err
		}
	}

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

func (s *Service) InitGroups() error {
	if s.db == nil {
		return fmt.Errorf("postgres db is not initialized")
	}

	log.Infof("Initializing groups")

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	query := `
		SELECT id FROM groups WHERE deleted_at IS NULL
	`

	var groups []schema.GroupSchema
	err = tx.SelectContext(context.Background(), &groups, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no groups found")
		}
		return err
	}

	for _, userGroup := range groups {
		g := group.NewGroup(s.db)
		if err := g.Init(context.Background(), userGroup.Id); err != nil {
			log.Errorf("Failed to initialize group %d: %v", userGroup.Id, err)
			continue
		}
		s.groups.Add(*g)
	}

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
		if errors.Is(err, user.ErrUserNotFound) {
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

	if err := u.LoadGroups(ctx); err != nil {
		log.Errorf("failed to load groups: %v", err)
		return &restproto.RestApiResponse{
			HttpCode: 401,
			Code:     12007,
			Error:    err.Error(),
		}, fmt.Errorf("failed to load user: %w", err)
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

	// We keep users cached until they logout or we didn't receive anything from them for a long period of time
	// @TODO: Handle timeout
	// @TODO: Handle cleanup of duplicates
	log.Debugf("User [%s] authenticated and cached", u.GetUsername())
	s.users.Add(*u)

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

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return &proto.ValidateTokenResponse{
			Code:  1,
			Error: err.Error(),
		}, err
	}

	sessionData := &schema.UserSessionSchema{}

	query := `
		SELECT user_id, created_at, updated_at
		FROM user_sessions
		WHERE token = $1 AND deleted_at IS NULL`

	log.Tracef("Query: %s", query)

	err = tx.GetContext(ctx, sessionData, query, in.Token)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.ValidateTokenResponse{
				Code:    0,
				IsValid: false,
			}, nil
		}
		log.Errorf("[Service::ValidateToken] Failed to get session data: %s", err.Error())
		return &proto.ValidateTokenResponse{
			Code:  4,
			Error: err.Error(),
		}, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[Service::ValidateToken] Failed to commit transaction: %s", err.Error())
		return &proto.ValidateTokenResponse{
			Code:    3,
			Error:   err.Error(),
			IsValid: false,
		}, err
	}

	if time.Since(sessionData.CreatedAt) > time.Duration(AppConfig.Crypto.JWT.Expiry)*time.Minute {
		log.Debugf("Token expired: %s. Duration: %+v", in.Token, time.Duration(AppConfig.Crypto.JWT.Expiry))
		return &proto.ValidateTokenResponse{
			Code:    0,
			IsValid: false,
			UserId:  int32(sessionData.UserId),
		}, nil
	}

	log.Debugf("Token is valid: %s", in.Token)
	return &proto.ValidateTokenResponse{
		Code:    0,
		IsValid: true,
		UserId:  int32(sessionData.UserId),
	}, nil
}

func (s *Service) HasPermission(ctx context.Context, in *proto.HasPermissionRequest) (*proto.HasPermissionResponse, error) {
	log.Tracef("HasPermission")

	if in.UserId == 0 {
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, fmt.Errorf("invalid user id")
	}

	user, err := s.users.GetById(in.UserId)
	if err != nil {
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, fmt.Errorf("user not found")
	}

	permission, err := user.GetPermission(in.Permission, in.Domain)
	if err != nil {
		return &proto.HasPermissionResponse{
			Read:   0,
			Write:  0,
			Delete: 0,
		}, err
	}

	return &proto.HasPermissionResponse{
		Read:   boolToInt32(permission.Read),
		Write:  boolToInt32(permission.Write),
		Delete: boolToInt32(permission.Delete),
	}, nil
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
