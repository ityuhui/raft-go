// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package raft_rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RaftServiceClient is the client API for RaftService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RaftServiceClient interface {
	// define the interface and data type
	TellMyHeartBeatToFollower(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatReply, error)
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error)
}

type raftServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftServiceClient(cc grpc.ClientConnInterface) RaftServiceClient {
	return &raftServiceClient{cc}
}

func (c *raftServiceClient) TellMyHeartBeatToFollower(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatReply, error) {
	out := new(HeartBeatReply)
	err := c.cc.Invoke(ctx, "/raft_rpc.RaftService/TellMyHeartBeatToFollower", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftServiceClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error) {
	out := new(VoteReply)
	err := c.cc.Invoke(ctx, "/raft_rpc.RaftService/Vote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftServiceServer is the server API for RaftService service.
// All implementations must embed UnimplementedRaftServiceServer
// for forward compatibility
type RaftServiceServer interface {
	// define the interface and data type
	TellMyHeartBeatToFollower(context.Context, *HeartBeatRequest) (*HeartBeatReply, error)
	Vote(context.Context, *VoteRequest) (*VoteReply, error)
	mustEmbedUnimplementedRaftServiceServer()
}

// UnimplementedRaftServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRaftServiceServer struct {
}

func (*UnimplementedRaftServiceServer) TellMyHeartBeatToFollower(context.Context, *HeartBeatRequest) (*HeartBeatReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TellMyHeartBeatToFollower not implemented")
}
func (*UnimplementedRaftServiceServer) Vote(context.Context, *VoteRequest) (*VoteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
func (*UnimplementedRaftServiceServer) mustEmbedUnimplementedRaftServiceServer() {}

func RegisterRaftServiceServer(s *grpc.Server, srv RaftServiceServer) {
	s.RegisterService(&_RaftService_serviceDesc, srv)
}

func _RaftService_TellMyHeartBeatToFollower_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).TellMyHeartBeatToFollower(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft_rpc.RaftService/TellMyHeartBeatToFollower",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).TellMyHeartBeatToFollower(ctx, req.(*HeartBeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftService_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServiceServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raft_rpc.RaftService/Vote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServiceServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RaftService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "raft_rpc.RaftService",
	HandlerType: (*RaftServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TellMyHeartBeatToFollower",
			Handler:    _RaftService_TellMyHeartBeatToFollower_Handler,
		},
		{
			MethodName: "Vote",
			Handler:    _RaftService_Vote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/raft.proto",
}
