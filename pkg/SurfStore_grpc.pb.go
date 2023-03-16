// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: pkg/SurfStore.proto

package pkg

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BlockStoreClient is the client API for BlockStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockStoreClient interface {
	GetBlock(ctx context.Context, in *BlockHash, opts ...grpc.CallOption) (*Block, error)
	PutBlock(ctx context.Context, in *Block, opts ...grpc.CallOption) (*Success, error)
	GetBlockStoreAddr(ctx context.Context, in *BlockHash, opts ...grpc.CallOption) (*BlockStoreAddr, error)
	CleanBlockStore(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*Success, error)
}

type blockStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockStoreClient(cc grpc.ClientConnInterface) BlockStoreClient {
	return &blockStoreClient{cc}
}

func (c *blockStoreClient) GetBlock(ctx context.Context, in *BlockHash, opts ...grpc.CallOption) (*Block, error) {
	out := new(Block)
	err := c.cc.Invoke(ctx, "/pkg.BlockStore/GetBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockStoreClient) PutBlock(ctx context.Context, in *Block, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/pkg.BlockStore/PutBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockStoreClient) GetBlockStoreAddr(ctx context.Context, in *BlockHash, opts ...grpc.CallOption) (*BlockStoreAddr, error) {
	out := new(BlockStoreAddr)
	err := c.cc.Invoke(ctx, "/pkg.BlockStore/GetBlockStoreAddr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockStoreClient) CleanBlockStore(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/pkg.BlockStore/CleanBlockStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockStoreServer is the server API for BlockStore service.
// All implementations must embed UnimplementedBlockStoreServer
// for forward compatibility
type BlockStoreServer interface {
	GetBlock(context.Context, *BlockHash) (*Block, error)
	PutBlock(context.Context, *Block) (*Success, error)
	GetBlockStoreAddr(context.Context, *BlockHash) (*BlockStoreAddr, error)
	CleanBlockStore(context.Context, *BlockHashes) (*Success, error)
	mustEmbedUnimplementedBlockStoreServer()
}

// UnimplementedBlockStoreServer must be embedded to have forward compatible implementations.
type UnimplementedBlockStoreServer struct {
}

func (UnimplementedBlockStoreServer) GetBlock(context.Context, *BlockHash) (*Block, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlock not implemented")
}
func (UnimplementedBlockStoreServer) PutBlock(context.Context, *Block) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutBlock not implemented")
}
func (UnimplementedBlockStoreServer) GetBlockStoreAddr(context.Context, *BlockHash) (*BlockStoreAddr, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockStoreAddr not implemented")
}
func (UnimplementedBlockStoreServer) CleanBlockStore(context.Context, *BlockHashes) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanBlockStore not implemented")
}
func (UnimplementedBlockStoreServer) mustEmbedUnimplementedBlockStoreServer() {}

// UnsafeBlockStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockStoreServer will
// result in compilation errors.
type UnsafeBlockStoreServer interface {
	mustEmbedUnimplementedBlockStoreServer()
}

func RegisterBlockStoreServer(s grpc.ServiceRegistrar, srv BlockStoreServer) {
	s.RegisterService(&BlockStore_ServiceDesc, srv)
}

func _BlockStore_GetBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHash)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).GetBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.BlockStore/GetBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).GetBlock(ctx, req.(*BlockHash))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockStore_PutBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Block)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).PutBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.BlockStore/PutBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).PutBlock(ctx, req.(*Block))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockStore_GetBlockStoreAddr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHash)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).GetBlockStoreAddr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.BlockStore/GetBlockStoreAddr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).GetBlockStoreAddr(ctx, req.(*BlockHash))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockStore_CleanBlockStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHashes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).CleanBlockStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.BlockStore/CleanBlockStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).CleanBlockStore(ctx, req.(*BlockHashes))
	}
	return interceptor(ctx, in, info, handler)
}

// BlockStore_ServiceDesc is the grpc.ServiceDesc for BlockStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pkg.BlockStore",
	HandlerType: (*BlockStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBlock",
			Handler:    _BlockStore_GetBlock_Handler,
		},
		{
			MethodName: "PutBlock",
			Handler:    _BlockStore_PutBlock_Handler,
		},
		{
			MethodName: "GetBlockStoreAddr",
			Handler:    _BlockStore_GetBlockStoreAddr_Handler,
		},
		{
			MethodName: "CleanBlockStore",
			Handler:    _BlockStore_CleanBlockStore_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/SurfStore.proto",
}

// MetaStoreClient is the client API for MetaStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetaStoreClient interface {
	// raft
	AppendEntries(ctx context.Context, in *AppendEntryInput, opts ...grpc.CallOption) (*AppendEntryOutput, error)
	RequestVote(ctx context.Context, in *CandidateState, opts ...grpc.CallOption) (*VoteResult, error)
	GetStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Status, error)
	// metastore
	GetFileInfoMap(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FileInfoMap, error)
	UpdateFile(ctx context.Context, in *FileMetaData, opts ...grpc.CallOption) (*Version, error)
	GetDefualtBlockStoreAddrs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BlockStoreAddrs, error)
}

type metaStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewMetaStoreClient(cc grpc.ClientConnInterface) MetaStoreClient {
	return &metaStoreClient{cc}
}

func (c *metaStoreClient) AppendEntries(ctx context.Context, in *AppendEntryInput, opts ...grpc.CallOption) (*AppendEntryOutput, error) {
	out := new(AppendEntryOutput)
	err := c.cc.Invoke(ctx, "/pkg.MetaStore/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) RequestVote(ctx context.Context, in *CandidateState, opts ...grpc.CallOption) (*VoteResult, error) {
	out := new(VoteResult)
	err := c.cc.Invoke(ctx, "/pkg.MetaStore/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) GetStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/pkg.MetaStore/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) GetFileInfoMap(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FileInfoMap, error) {
	out := new(FileInfoMap)
	err := c.cc.Invoke(ctx, "/pkg.MetaStore/GetFileInfoMap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) UpdateFile(ctx context.Context, in *FileMetaData, opts ...grpc.CallOption) (*Version, error) {
	out := new(Version)
	err := c.cc.Invoke(ctx, "/pkg.MetaStore/UpdateFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) GetDefualtBlockStoreAddrs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BlockStoreAddrs, error) {
	out := new(BlockStoreAddrs)
	err := c.cc.Invoke(ctx, "/pkg.MetaStore/GetDefualtBlockStoreAddrs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetaStoreServer is the server API for MetaStore service.
// All implementations must embed UnimplementedMetaStoreServer
// for forward compatibility
type MetaStoreServer interface {
	// raft
	AppendEntries(context.Context, *AppendEntryInput) (*AppendEntryOutput, error)
	RequestVote(context.Context, *CandidateState) (*VoteResult, error)
	GetStatus(context.Context, *emptypb.Empty) (*Status, error)
	// metastore
	GetFileInfoMap(context.Context, *emptypb.Empty) (*FileInfoMap, error)
	UpdateFile(context.Context, *FileMetaData) (*Version, error)
	GetDefualtBlockStoreAddrs(context.Context, *emptypb.Empty) (*BlockStoreAddrs, error)
	mustEmbedUnimplementedMetaStoreServer()
}

// UnimplementedMetaStoreServer must be embedded to have forward compatible implementations.
type UnimplementedMetaStoreServer struct {
}

func (UnimplementedMetaStoreServer) AppendEntries(context.Context, *AppendEntryInput) (*AppendEntryOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedMetaStoreServer) RequestVote(context.Context, *CandidateState) (*VoteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedMetaStoreServer) GetStatus(context.Context, *emptypb.Empty) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedMetaStoreServer) GetFileInfoMap(context.Context, *emptypb.Empty) (*FileInfoMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileInfoMap not implemented")
}
func (UnimplementedMetaStoreServer) UpdateFile(context.Context, *FileMetaData) (*Version, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFile not implemented")
}
func (UnimplementedMetaStoreServer) GetDefualtBlockStoreAddrs(context.Context, *emptypb.Empty) (*BlockStoreAddrs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefualtBlockStoreAddrs not implemented")
}
func (UnimplementedMetaStoreServer) mustEmbedUnimplementedMetaStoreServer() {}

// UnsafeMetaStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetaStoreServer will
// result in compilation errors.
type UnsafeMetaStoreServer interface {
	mustEmbedUnimplementedMetaStoreServer()
}

func RegisterMetaStoreServer(s grpc.ServiceRegistrar, srv MetaStoreServer) {
	s.RegisterService(&MetaStore_ServiceDesc, srv)
}

func _MetaStore_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntryInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.MetaStore/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).AppendEntries(ctx, req.(*AppendEntryInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CandidateState)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.MetaStore/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).RequestVote(ctx, req.(*CandidateState))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.MetaStore/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).GetStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_GetFileInfoMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).GetFileInfoMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.MetaStore/GetFileInfoMap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).GetFileInfoMap(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_UpdateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileMetaData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).UpdateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.MetaStore/UpdateFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).UpdateFile(ctx, req.(*FileMetaData))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_GetDefualtBlockStoreAddrs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).GetDefualtBlockStoreAddrs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.MetaStore/GetDefualtBlockStoreAddrs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).GetDefualtBlockStoreAddrs(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MetaStore_ServiceDesc is the grpc.ServiceDesc for MetaStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetaStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pkg.MetaStore",
	HandlerType: (*MetaStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _MetaStore_AppendEntries_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _MetaStore_RequestVote_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _MetaStore_GetStatus_Handler,
		},
		{
			MethodName: "GetFileInfoMap",
			Handler:    _MetaStore_GetFileInfoMap_Handler,
		},
		{
			MethodName: "UpdateFile",
			Handler:    _MetaStore_UpdateFile_Handler,
		},
		{
			MethodName: "GetDefualtBlockStoreAddrs",
			Handler:    _MetaStore_GetDefualtBlockStoreAddrs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/SurfStore.proto",
}