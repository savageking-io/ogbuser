// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: user.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UserService_Ping_FullMethodName                        = "/user.UserService/Ping"
	UserService_AuthenticateUserCredentials_FullMethodName = "/user.UserService/AuthenticateUserCredentials"
	UserService_AuthenticatePlatform_FullMethodName        = "/user.UserService/AuthenticatePlatform"
	UserService_AuthenticateServer_FullMethodName          = "/user.UserService/AuthenticateServer"
	UserService_AuthenticateWebSocketToken_FullMethodName  = "/user.UserService/AuthenticateWebSocketToken"
	UserService_HasPermission_FullMethodName               = "/user.UserService/HasPermission"
	UserService_ValidateToken_FullMethodName               = "/user.UserService/ValidateToken"
	UserService_RenewToken_FullMethodName                  = "/user.UserService/RenewToken"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*PingMessage, error)
	AuthenticateUserCredentials(ctx context.Context, in *AuthUserCredentialsRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	AuthenticatePlatform(ctx context.Context, in *AuthPlatformRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	AuthenticateServer(ctx context.Context, in *AuthServerRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	AuthenticateWebSocketToken(ctx context.Context, in *AuthWebSocketTokenRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	HasPermission(ctx context.Context, in *HasPermissionRequest, opts ...grpc.CallOption) (*HasPermissionResponse, error)
	ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error)
	RenewToken(ctx context.Context, in *RenewTokenRequest, opts ...grpc.CallOption) (*RenewTokenResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*PingMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingMessage)
	err := c.cc.Invoke(ctx, UserService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AuthenticateUserCredentials(ctx context.Context, in *AuthUserCredentialsRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, UserService_AuthenticateUserCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AuthenticatePlatform(ctx context.Context, in *AuthPlatformRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, UserService_AuthenticatePlatform_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AuthenticateServer(ctx context.Context, in *AuthServerRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, UserService_AuthenticateServer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AuthenticateWebSocketToken(ctx context.Context, in *AuthWebSocketTokenRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, UserService_AuthenticateWebSocketToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) HasPermission(ctx context.Context, in *HasPermissionRequest, opts ...grpc.CallOption) (*HasPermissionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HasPermissionResponse)
	err := c.cc.Invoke(ctx, UserService_HasPermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidateTokenResponse)
	err := c.cc.Invoke(ctx, UserService_ValidateToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RenewToken(ctx context.Context, in *RenewTokenRequest, opts ...grpc.CallOption) (*RenewTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RenewTokenResponse)
	err := c.cc.Invoke(ctx, UserService_RenewToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility.
type UserServiceServer interface {
	Ping(context.Context, *PingMessage) (*PingMessage, error)
	AuthenticateUserCredentials(context.Context, *AuthUserCredentialsRequest) (*AuthResponse, error)
	AuthenticatePlatform(context.Context, *AuthPlatformRequest) (*AuthResponse, error)
	AuthenticateServer(context.Context, *AuthServerRequest) (*AuthResponse, error)
	AuthenticateWebSocketToken(context.Context, *AuthWebSocketTokenRequest) (*AuthResponse, error)
	HasPermission(context.Context, *HasPermissionRequest) (*HasPermissionResponse, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error)
	RenewToken(context.Context, *RenewTokenRequest) (*RenewTokenResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) Ping(context.Context, *PingMessage) (*PingMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedUserServiceServer) AuthenticateUserCredentials(context.Context, *AuthUserCredentialsRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateUserCredentials not implemented")
}
func (UnimplementedUserServiceServer) AuthenticatePlatform(context.Context, *AuthPlatformRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticatePlatform not implemented")
}
func (UnimplementedUserServiceServer) AuthenticateServer(context.Context, *AuthServerRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateServer not implemented")
}
func (UnimplementedUserServiceServer) AuthenticateWebSocketToken(context.Context, *AuthWebSocketTokenRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateWebSocketToken not implemented")
}
func (UnimplementedUserServiceServer) HasPermission(context.Context, *HasPermissionRequest) (*HasPermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasPermission not implemented")
}
func (UnimplementedUserServiceServer) ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (UnimplementedUserServiceServer) RenewToken(context.Context, *RenewTokenRequest) (*RenewTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenewToken not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}
func (UnimplementedUserServiceServer) testEmbeddedByValue()                     {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	// If the following call pancis, it indicates UnimplementedUserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Ping(ctx, req.(*PingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AuthenticateUserCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthUserCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AuthenticateUserCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AuthenticateUserCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AuthenticateUserCredentials(ctx, req.(*AuthUserCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AuthenticatePlatform_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthPlatformRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AuthenticatePlatform(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AuthenticatePlatform_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AuthenticatePlatform(ctx, req.(*AuthPlatformRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AuthenticateServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AuthenticateServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AuthenticateServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AuthenticateServer(ctx, req.(*AuthServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AuthenticateWebSocketToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthWebSocketTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AuthenticateWebSocketToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AuthenticateWebSocketToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AuthenticateWebSocketToken(ctx, req.(*AuthWebSocketTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_HasPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HasPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).HasPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_HasPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).HasPermission(ctx, req.(*HasPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_ValidateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ValidateToken(ctx, req.(*ValidateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RenewToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenewTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RenewToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_RenewToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RenewToken(ctx, req.(*RenewTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _UserService_Ping_Handler,
		},
		{
			MethodName: "AuthenticateUserCredentials",
			Handler:    _UserService_AuthenticateUserCredentials_Handler,
		},
		{
			MethodName: "AuthenticatePlatform",
			Handler:    _UserService_AuthenticatePlatform_Handler,
		},
		{
			MethodName: "AuthenticateServer",
			Handler:    _UserService_AuthenticateServer_Handler,
		},
		{
			MethodName: "AuthenticateWebSocketToken",
			Handler:    _UserService_AuthenticateWebSocketToken_Handler,
		},
		{
			MethodName: "HasPermission",
			Handler:    _UserService_HasPermission_Handler,
		},
		{
			MethodName: "ValidateToken",
			Handler:    _UserService_ValidateToken_Handler,
		},
		{
			MethodName: "RenewToken",
			Handler:    _UserService_RenewToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
