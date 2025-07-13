package main

import (
	"context"
	restproto "github.com/savageking-io/ogbrest/proto"
	"github.com/savageking-io/ogbrest/restlib"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	rest *restlib.RestInterServiceServer
}

func (s *Service) Init() error {
	return nil
}

func (s *Service) InitializeRest(config restlib.RestInterServiceConfig) error {
	log.Infof("Initializing REST service")
	s.rest = restlib.NewRestInterServiceServer(config)
	if err := s.rest.Init(); err != nil {
		return err
	}
	if err := s.rest.RegisterHandler("/auth/credentials", "POST", s.HandleAuthCredentialsRequest); err != nil {
		log.Warnf("Failed to register handler for /auth/credentials: %v", err)
	}
	if err := s.rest.RegisterHandler("/auth/platform", "POST", s.HandleAuthPlatformRequest); err != nil {
		log.Warnf("Failed to register handler for /auth/platform: %v", err)
	}
	if err := s.rest.RegisterHandler("/auth/server", "POST", s.HandleAuthServerRequest); err != nil {
		log.Warnf("Failed to register handler for /auth/server: %v", err)
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

	go func() {
		log.Infof("Starting REST server: %s:%d", AppConfig.Rest.Hostname, AppConfig.Rest.Port)
		err := s.rest.Start()
		restErrChan <- err
	}()

	select {
	case err := <-restErrChan:
		if err != nil {
			log.Errorf("REST channel reported error: %v", err)
		}
	}

	return nil
}

func (s *Service) HandleAuthCredentialsRequest(ctx context.Context, in *restproto.RestApiRequest) (*restproto.RestApiResponse, error) {
	return nil, nil
}

func (s *Service) HandleAuthPlatformRequest(ctx context.Context, in *restproto.RestApiRequest) (*restproto.RestApiResponse, error) {
	return nil, nil
}

func (s *Service) HandleAuthServerRequest(ctx context.Context, in *restproto.RestApiRequest) (*restproto.RestApiResponse, error) {
	return nil, nil
}
