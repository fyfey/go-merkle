// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// MerkleClient is the client API for Merkle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MerkleClient interface {
	GetPart(ctx context.Context, in *PartRequest, opts ...grpc.CallOption) (*Part, error)
	GetMetadata(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Metadata, error)
}

type merkleClient struct {
	cc grpc.ClientConnInterface
}

func NewMerkleClient(cc grpc.ClientConnInterface) MerkleClient {
	return &merkleClient{cc}
}

func (c *merkleClient) GetPart(ctx context.Context, in *PartRequest, opts ...grpc.CallOption) (*Part, error) {
	out := new(Part)
	err := c.cc.Invoke(ctx, "/merkle.Merkle/GetPart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *merkleClient) GetMetadata(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Metadata, error) {
	out := new(Metadata)
	err := c.cc.Invoke(ctx, "/merkle.Merkle/GetMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MerkleServer is the server API for Merkle service.
// All implementations must embed UnimplementedMerkleServer
// for forward compatibility
type MerkleServer interface {
	GetPart(context.Context, *PartRequest) (*Part, error)
	GetMetadata(context.Context, *Empty) (*Metadata, error)
	mustEmbedUnimplementedMerkleServer()
}

// UnimplementedMerkleServer must be embedded to have forward compatible implementations.
type UnimplementedMerkleServer struct {
}

func (UnimplementedMerkleServer) GetPart(context.Context, *PartRequest) (*Part, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPart not implemented")
}
func (UnimplementedMerkleServer) GetMetadata(context.Context, *Empty) (*Metadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetadata not implemented")
}
func (UnimplementedMerkleServer) mustEmbedUnimplementedMerkleServer() {}

// UnsafeMerkleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MerkleServer will
// result in compilation errors.
type UnsafeMerkleServer interface {
	mustEmbedUnimplementedMerkleServer()
}

func RegisterMerkleServer(s grpc.ServiceRegistrar, srv MerkleServer) {
	s.RegisterService(&Merkle_ServiceDesc, srv)
}

func _Merkle_GetPart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MerkleServer).GetPart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/merkle.Merkle/GetPart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MerkleServer).GetPart(ctx, req.(*PartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Merkle_GetMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MerkleServer).GetMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/merkle.Merkle/GetMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MerkleServer).GetMetadata(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Merkle_ServiceDesc is the grpc.ServiceDesc for Merkle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Merkle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "merkle.Merkle",
	HandlerType: (*MerkleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPart",
			Handler:    _Merkle_GetPart_Handler,
		},
		{
			MethodName: "GetMetadata",
			Handler:    _Merkle_GetMetadata_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "merkle.proto",
}