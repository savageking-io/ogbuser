package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	restproto "github.com/savageking-io/ogbrest/proto"
	"github.com/savageking-io/ogbrest/restlib"
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
	groups []*Group
	users  []*user.User

	proto.UnimplementedUserServiceServer
}

func NewService(config *ServiceConfig) *Service {
	return &Service{
		config: config,
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

	for _, group := range groups {
		g := NewGroup(s.db)
		if err := g.Init(context.Background(), group.Id); err != nil {
			log.Errorf("Failed to initialize group %d: %v", group.Id, err)
			continue
		}
		s.groups = append(s.groups, g)
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
	if err := s.rest.RegisterHandler("/auth/credentials", "POST", s.HandleAuthCredentialsRequest); err != nil {
		log.Warnf("Failed to register handler for /auth/credentials: %v", err)
	}
	if err := s.rest.RegisterHandler("/auth/platform", "POST", s.HandleAuthPlatformRequest); err != nil {
		log.Warnf("Failed to register handler for /auth/platform: %v", err)
	}
	if err := s.rest.RegisterHandler("/auth/server", "POST", s.HandleAuthServerRequest); err != nil {
		log.Warnf("Failed to register handler for /auth/server: %v", err)
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
	s.users = append(s.users, u)

	session, err := u.InitializeSession(ctx)
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

// UserService
func (s *Service) Ping(ctx context.Context, in *proto.PingMessage) (*proto.PingMessage, error) {
	in.RepliedAt = timestamppb.New(time.Now())
	return in, nil
}
