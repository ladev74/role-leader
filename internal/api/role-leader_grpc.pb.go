// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: role-leader.proto

package api

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
	RoleLeader_CreateFeedback_FullMethodName = "/api.RoleLeader/CreateFeedback"
	RoleLeader_GetCall_FullMethodName        = "/api.RoleLeader/GetCall"
	RoleLeader_GetLeaderCalls_FullMethodName = "/api.RoleLeader/GetLeaderCalls"
)

// RoleLeaderClient is the client API for RoleLeader service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoleLeaderClient interface {
	CreateFeedback(ctx context.Context, in *CreateFeedbackRequest, opts ...grpc.CallOption) (*CreateFeedbackResponse, error)
	GetCall(ctx context.Context, in *GetCallRequest, opts ...grpc.CallOption) (*GetCallResponse, error)
	GetLeaderCalls(ctx context.Context, in *GetLeaderCallsRequest, opts ...grpc.CallOption) (*GetLeaderCallsResponse, error)
}

type roleLeaderClient struct {
	cc grpc.ClientConnInterface
}

func NewRoleLeaderClient(cc grpc.ClientConnInterface) RoleLeaderClient {
	return &roleLeaderClient{cc}
}

func (c *roleLeaderClient) CreateFeedback(ctx context.Context, in *CreateFeedbackRequest, opts ...grpc.CallOption) (*CreateFeedbackResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFeedbackResponse)
	err := c.cc.Invoke(ctx, RoleLeader_CreateFeedback_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleLeaderClient) GetCall(ctx context.Context, in *GetCallRequest, opts ...grpc.CallOption) (*GetCallResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCallResponse)
	err := c.cc.Invoke(ctx, RoleLeader_GetCall_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleLeaderClient) GetLeaderCalls(ctx context.Context, in *GetLeaderCallsRequest, opts ...grpc.CallOption) (*GetLeaderCallsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetLeaderCallsResponse)
	err := c.cc.Invoke(ctx, RoleLeader_GetLeaderCalls_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoleLeaderServer is the server API for RoleLeader service.
// All implementations must embed UnimplementedRoleLeaderServer
// for forward compatibility.
type RoleLeaderServer interface {
	CreateFeedback(context.Context, *CreateFeedbackRequest) (*CreateFeedbackResponse, error)
	GetCall(context.Context, *GetCallRequest) (*GetCallResponse, error)
	GetLeaderCalls(context.Context, *GetLeaderCallsRequest) (*GetLeaderCallsResponse, error)
	mustEmbedUnimplementedRoleLeaderServer()
}

// UnimplementedRoleLeaderServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRoleLeaderServer struct{}

func (UnimplementedRoleLeaderServer) CreateFeedback(context.Context, *CreateFeedbackRequest) (*CreateFeedbackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFeedback not implemented")
}
func (UnimplementedRoleLeaderServer) GetCall(context.Context, *GetCallRequest) (*GetCallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCall not implemented")
}
func (UnimplementedRoleLeaderServer) GetLeaderCalls(context.Context, *GetLeaderCallsRequest) (*GetLeaderCallsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeaderCalls not implemented")
}
func (UnimplementedRoleLeaderServer) mustEmbedUnimplementedRoleLeaderServer() {}
func (UnimplementedRoleLeaderServer) testEmbeddedByValue()                    {}

// UnsafeRoleLeaderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoleLeaderServer will
// result in compilation errors.
type UnsafeRoleLeaderServer interface {
	mustEmbedUnimplementedRoleLeaderServer()
}

func RegisterRoleLeaderServer(s grpc.ServiceRegistrar, srv RoleLeaderServer) {
	// If the following call pancis, it indicates UnimplementedRoleLeaderServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RoleLeader_ServiceDesc, srv)
}

func _RoleLeader_CreateFeedback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFeedbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleLeaderServer).CreateFeedback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleLeader_CreateFeedback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleLeaderServer).CreateFeedback(ctx, req.(*CreateFeedbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleLeader_GetCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleLeaderServer).GetCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleLeader_GetCall_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleLeaderServer).GetCall(ctx, req.(*GetCallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleLeader_GetLeaderCalls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLeaderCallsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleLeaderServer).GetLeaderCalls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleLeader_GetLeaderCalls_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleLeaderServer).GetLeaderCalls(ctx, req.(*GetLeaderCallsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoleLeader_ServiceDesc is the grpc.ServiceDesc for RoleLeader service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoleLeader_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.RoleLeader",
	HandlerType: (*RoleLeaderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFeedback",
			Handler:    _RoleLeader_CreateFeedback_Handler,
		},
		{
			MethodName: "GetCall",
			Handler:    _RoleLeader_GetCall_Handler,
		},
		{
			MethodName: "GetLeaderCalls",
			Handler:    _RoleLeader_GetLeaderCalls_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "role-leader.proto",
}
