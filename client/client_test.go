package client

import (
	"context"
	"github.com/savageking-io/ogbuser/proto"
	"google.golang.org/grpc"
	"reflect"
	"testing"
)

func TestClient_Connect(t *testing.T) {
	type fields struct {
		hostname  string
		port      uint16
		ErrorChan chan error
		conn      *grpc.ClientConn
		client    proto.UserServiceClient
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				hostname:  tt.fields.hostname,
				port:      tt.fields.port,
				ErrorChan: tt.fields.ErrorChan,
				conn:      tt.fields.conn,
				client:    tt.fields.client,
			}
			if err := c.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ConnectWithRetry(t *testing.T) {
	type fields struct {
		hostname  string
		port      uint16
		ErrorChan chan error
		conn      *grpc.ClientConn
		client    proto.UserServiceClient
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				hostname:  tt.fields.hostname,
				port:      tt.fields.port,
				ErrorChan: tt.fields.ErrorChan,
				conn:      tt.fields.conn,
				client:    tt.fields.client,
			}
			c.ConnectWithRetry()
		})
	}
}

func TestClient_Disconnect(t *testing.T) {
	type fields struct {
		hostname  string
		port      uint16
		ErrorChan chan error
		conn      *grpc.ClientConn
		client    proto.UserServiceClient
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				hostname:  tt.fields.hostname,
				port:      tt.fields.port,
				ErrorChan: tt.fields.ErrorChan,
				conn:      tt.fields.conn,
				client:    tt.fields.client,
			}
			if err := c.Disconnect(); (err != nil) != tt.wantErr {
				t.Errorf("Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Reconnect(t *testing.T) {
	type fields struct {
		hostname  string
		port      uint16
		ErrorChan chan error
		conn      *grpc.ClientConn
		client    proto.UserServiceClient
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				hostname:  tt.fields.hostname,
				port:      tt.fields.port,
				ErrorChan: tt.fields.ErrorChan,
				conn:      tt.fields.conn,
				client:    tt.fields.client,
			}
			c.Reconnect()
		})
	}
}

func TestClient_Run(t *testing.T) {
	type fields struct {
		hostname  string
		port      uint16
		ErrorChan chan error
		conn      *grpc.ClientConn
		client    proto.UserServiceClient
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				hostname:  tt.fields.hostname,
				port:      tt.fields.port,
				ErrorChan: tt.fields.ErrorChan,
				conn:      tt.fields.conn,
				client:    tt.fields.client,
			}
			if err := c.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ValidateToken(t *testing.T) {
	type fields struct {
		hostname  string
		port      uint16
		ErrorChan chan error
		conn      *grpc.ClientConn
		client    proto.UserServiceClient
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		want1   int32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				hostname:  tt.fields.hostname,
				port:      tt.fields.port,
				ErrorChan: tt.fields.ErrorChan,
				conn:      tt.fields.conn,
				client:    tt.fields.client,
			}
			got, got1, err := c.ValidateToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_ping(t *testing.T) {
	type fields struct {
		hostname  string
		port      uint16
		ErrorChan chan error
		conn      *grpc.ClientConn
		client    proto.UserServiceClient
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				hostname:  tt.fields.hostname,
				port:      tt.fields.port,
				ErrorChan: tt.fields.ErrorChan,
				conn:      tt.fields.conn,
				client:    tt.fields.client,
			}
			if err := c.ping(); (err != nil) != tt.wantErr {
				t.Errorf("ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		hostname string
		port     uint16
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.hostname, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
