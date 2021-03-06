// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protocol

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RepoClient is the client API for Repo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepoClient interface {
	ListRepos(ctx context.Context, in *ListRepoRequest, opts ...grpc.CallOption) (*ListRepoResponse, error)
	ListBranchesOfRepo(ctx context.Context, in *ListBranchesOfRepoRequest, opts ...grpc.CallOption) (*ListBranchesOfRepoResponse, error)
}

type repoClient struct {
	cc grpc.ClientConnInterface
}

func NewRepoClient(cc grpc.ClientConnInterface) RepoClient {
	return &repoClient{cc}
}

func (c *repoClient) ListRepos(ctx context.Context, in *ListRepoRequest, opts ...grpc.CallOption) (*ListRepoResponse, error) {
	out := new(ListRepoResponse)
	err := c.cc.Invoke(ctx, "/repo.Repo/ListRepos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repoClient) ListBranchesOfRepo(ctx context.Context, in *ListBranchesOfRepoRequest, opts ...grpc.CallOption) (*ListBranchesOfRepoResponse, error) {
	out := new(ListBranchesOfRepoResponse)
	err := c.cc.Invoke(ctx, "/repo.Repo/ListBranchesOfRepo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RepoServer is the server API for Repo service.
// All implementations must embed UnimplementedRepoServer
// for forward compatibility
type RepoServer interface {
	ListRepos(context.Context, *ListRepoRequest) (*ListRepoResponse, error)
	ListBranchesOfRepo(context.Context, *ListBranchesOfRepoRequest) (*ListBranchesOfRepoResponse, error)
	mustEmbedUnimplementedRepoServer()
}

// UnimplementedRepoServer must be embedded to have forward compatible implementations.
type UnimplementedRepoServer struct {
}

func (UnimplementedRepoServer) ListRepos(context.Context, *ListRepoRequest) (*ListRepoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRepos not implemented")
}
func (UnimplementedRepoServer) ListBranchesOfRepo(context.Context, *ListBranchesOfRepoRequest) (*ListBranchesOfRepoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBranchesOfRepo not implemented")
}
func (UnimplementedRepoServer) mustEmbedUnimplementedRepoServer() {}

// UnsafeRepoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepoServer will
// result in compilation errors.
type UnsafeRepoServer interface {
	mustEmbedUnimplementedRepoServer()
}

func RegisterRepoServer(s grpc.ServiceRegistrar, srv RepoServer) {
	s.RegisterService(&Repo_ServiceDesc, srv)
}

func _Repo_ListRepos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRepoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepoServer).ListRepos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repo.Repo/ListRepos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepoServer).ListRepos(ctx, req.(*ListRepoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Repo_ListBranchesOfRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBranchesOfRepoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepoServer).ListBranchesOfRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repo.Repo/ListBranchesOfRepo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepoServer).ListBranchesOfRepo(ctx, req.(*ListBranchesOfRepoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Repo_ServiceDesc is the grpc.ServiceDesc for Repo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Repo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "repo.Repo",
	HandlerType: (*RepoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRepos",
			Handler:    _Repo_ListRepos_Handler,
		},
		{
			MethodName: "ListBranchesOfRepo",
			Handler:    _Repo_ListBranchesOfRepo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "repo.proto",
}
