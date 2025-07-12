package main

import (
	"context"
	restpb "github.com/savageking-io/ogbrest/proto"
	"github.com/savageking-io/ogbrest/restlib"
	pb "github.com/savageking-io/ogbuser/proto"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	token string
}

// AuthService authenticates a service using the provided token
func (s *Server) AuthService(ctx context.Context, req *restpb.AuthenticateServiceRequest) (*restpb.AuthenticateServiceResponse, error) {
	restlib.HandleIncomingAuthRequest(req)
}

// RequestRESTData provides information about REST endpoints for User service
func (s *Server) RequestRESTData(ctx context.Context, req *restpb.RestDataRequest) (*restpb.RestDataResponse, error) {
	log.Debugf("Received REST data request for version: %s", req.Version)

	var endpointsList []*restpb.RestEndpoint
	for _, endpoint := range AppConfig.Rest.Endpoints {
		endpointsList = append(endpointsList, &restpb.RestEndpoint{
			Path:   endpoint.Endpoint,
			Method: endpoint.Method,
		})
	}

	resp := &restpb.RestDataResponse{
		Code:         200,
		EndpointsNum: int32(len(AppConfig.Rest.Endpoints)),
		Root:         AppConfig.Rest.Root,
		Endpoints:    endpointsList,
	}

	log.Debug("REST data request processed successfully")
	return resp, nil
}
