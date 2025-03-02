// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api_stores.proto

package storespb

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
	StoresService_CreateStore_FullMethodName            = "/storespb.StoresService/CreateStore"
	StoresService_GetStore_FullMethodName               = "/storespb.StoresService/GetStore"
	StoresService_GetStores_FullMethodName              = "/storespb.StoresService/GetStores"
	StoresService_EnableParticipation_FullMethodName    = "/storespb.StoresService/EnableParticipation"
	StoresService_DisableParticipation_FullMethodName   = "/storespb.StoresService/DisableParticipation"
	StoresService_GetParticipatingStores_FullMethodName = "/storespb.StoresService/GetParticipatingStores"
	StoresService_AddProduct_FullMethodName             = "/storespb.StoresService/AddProduct"
	StoresService_RemoveProduct_FullMethodName          = "/storespb.StoresService/RemoveProduct"
	StoresService_GetProduct_FullMethodName             = "/storespb.StoresService/GetProduct"
	StoresService_GetCatalog_FullMethodName             = "/storespb.StoresService/GetCatalog"
)

// StoresServiceClient is the client API for StoresService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StoresServiceClient interface {
	CreateStore(ctx context.Context, in *CreateStoreRequest, opts ...grpc.CallOption) (*CreateStoreResponse, error)
	GetStore(ctx context.Context, in *GetStoreRequest, opts ...grpc.CallOption) (*GetStoreResponse, error)
	GetStores(ctx context.Context, in *GetStoresRequest, opts ...grpc.CallOption) (*GetStoresResponse, error)
	EnableParticipation(ctx context.Context, in *EnableParticipationRequest, opts ...grpc.CallOption) (*EnableParticipationResponse, error)
	DisableParticipation(ctx context.Context, in *DisableParticipationRequest, opts ...grpc.CallOption) (*DisableParticipationResponse, error)
	GetParticipatingStores(ctx context.Context, in *GetParticipatingStoresRequest, opts ...grpc.CallOption) (*GetParticipatingStoresResponse, error)
	AddProduct(ctx context.Context, in *AddProductRequest, opts ...grpc.CallOption) (*AddProductResponse, error)
	RemoveProduct(ctx context.Context, in *RemoveProductRequest, opts ...grpc.CallOption) (*RemoveProductResponse, error)
	GetProduct(ctx context.Context, in *GetProductRequest, opts ...grpc.CallOption) (*GetProductResponse, error)
	GetCatalog(ctx context.Context, in *GetCatalogRequest, opts ...grpc.CallOption) (*GetCatalogResponse, error)
}

type storesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStoresServiceClient(cc grpc.ClientConnInterface) StoresServiceClient {
	return &storesServiceClient{cc}
}

func (c *storesServiceClient) CreateStore(ctx context.Context, in *CreateStoreRequest, opts ...grpc.CallOption) (*CreateStoreResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateStoreResponse)
	err := c.cc.Invoke(ctx, StoresService_CreateStore_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) GetStore(ctx context.Context, in *GetStoreRequest, opts ...grpc.CallOption) (*GetStoreResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetStoreResponse)
	err := c.cc.Invoke(ctx, StoresService_GetStore_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) GetStores(ctx context.Context, in *GetStoresRequest, opts ...grpc.CallOption) (*GetStoresResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetStoresResponse)
	err := c.cc.Invoke(ctx, StoresService_GetStores_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) EnableParticipation(ctx context.Context, in *EnableParticipationRequest, opts ...grpc.CallOption) (*EnableParticipationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EnableParticipationResponse)
	err := c.cc.Invoke(ctx, StoresService_EnableParticipation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) DisableParticipation(ctx context.Context, in *DisableParticipationRequest, opts ...grpc.CallOption) (*DisableParticipationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DisableParticipationResponse)
	err := c.cc.Invoke(ctx, StoresService_DisableParticipation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) GetParticipatingStores(ctx context.Context, in *GetParticipatingStoresRequest, opts ...grpc.CallOption) (*GetParticipatingStoresResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetParticipatingStoresResponse)
	err := c.cc.Invoke(ctx, StoresService_GetParticipatingStores_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) AddProduct(ctx context.Context, in *AddProductRequest, opts ...grpc.CallOption) (*AddProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddProductResponse)
	err := c.cc.Invoke(ctx, StoresService_AddProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) RemoveProduct(ctx context.Context, in *RemoveProductRequest, opts ...grpc.CallOption) (*RemoveProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveProductResponse)
	err := c.cc.Invoke(ctx, StoresService_RemoveProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) GetProduct(ctx context.Context, in *GetProductRequest, opts ...grpc.CallOption) (*GetProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProductResponse)
	err := c.cc.Invoke(ctx, StoresService_GetProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storesServiceClient) GetCatalog(ctx context.Context, in *GetCatalogRequest, opts ...grpc.CallOption) (*GetCatalogResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCatalogResponse)
	err := c.cc.Invoke(ctx, StoresService_GetCatalog_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StoresServiceServer is the server API for StoresService service.
// All implementations must embed UnimplementedStoresServiceServer
// for forward compatibility.
type StoresServiceServer interface {
	CreateStore(context.Context, *CreateStoreRequest) (*CreateStoreResponse, error)
	GetStore(context.Context, *GetStoreRequest) (*GetStoreResponse, error)
	GetStores(context.Context, *GetStoresRequest) (*GetStoresResponse, error)
	EnableParticipation(context.Context, *EnableParticipationRequest) (*EnableParticipationResponse, error)
	DisableParticipation(context.Context, *DisableParticipationRequest) (*DisableParticipationResponse, error)
	GetParticipatingStores(context.Context, *GetParticipatingStoresRequest) (*GetParticipatingStoresResponse, error)
	AddProduct(context.Context, *AddProductRequest) (*AddProductResponse, error)
	RemoveProduct(context.Context, *RemoveProductRequest) (*RemoveProductResponse, error)
	GetProduct(context.Context, *GetProductRequest) (*GetProductResponse, error)
	GetCatalog(context.Context, *GetCatalogRequest) (*GetCatalogResponse, error)
	mustEmbedUnimplementedStoresServiceServer()
}

// UnimplementedStoresServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStoresServiceServer struct{}

func (UnimplementedStoresServiceServer) CreateStore(context.Context, *CreateStoreRequest) (*CreateStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStore not implemented")
}
func (UnimplementedStoresServiceServer) GetStore(context.Context, *GetStoreRequest) (*GetStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStore not implemented")
}
func (UnimplementedStoresServiceServer) GetStores(context.Context, *GetStoresRequest) (*GetStoresResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStores not implemented")
}
func (UnimplementedStoresServiceServer) EnableParticipation(context.Context, *EnableParticipationRequest) (*EnableParticipationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableParticipation not implemented")
}
func (UnimplementedStoresServiceServer) DisableParticipation(context.Context, *DisableParticipationRequest) (*DisableParticipationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableParticipation not implemented")
}
func (UnimplementedStoresServiceServer) GetParticipatingStores(context.Context, *GetParticipatingStoresRequest) (*GetParticipatingStoresResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetParticipatingStores not implemented")
}
func (UnimplementedStoresServiceServer) AddProduct(context.Context, *AddProductRequest) (*AddProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProduct not implemented")
}
func (UnimplementedStoresServiceServer) RemoveProduct(context.Context, *RemoveProductRequest) (*RemoveProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveProduct not implemented")
}
func (UnimplementedStoresServiceServer) GetProduct(context.Context, *GetProductRequest) (*GetProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}
func (UnimplementedStoresServiceServer) GetCatalog(context.Context, *GetCatalogRequest) (*GetCatalogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCatalog not implemented")
}
func (UnimplementedStoresServiceServer) mustEmbedUnimplementedStoresServiceServer() {}
func (UnimplementedStoresServiceServer) testEmbeddedByValue()                       {}

// UnsafeStoresServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StoresServiceServer will
// result in compilation errors.
type UnsafeStoresServiceServer interface {
	mustEmbedUnimplementedStoresServiceServer()
}

func RegisterStoresServiceServer(s grpc.ServiceRegistrar, srv StoresServiceServer) {
	// If the following call pancis, it indicates UnimplementedStoresServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StoresService_ServiceDesc, srv)
}

func _StoresService_CreateStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).CreateStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_CreateStore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).CreateStore(ctx, req.(*CreateStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_GetStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).GetStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_GetStore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).GetStore(ctx, req.(*GetStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_GetStores_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStoresRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).GetStores(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_GetStores_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).GetStores(ctx, req.(*GetStoresRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_EnableParticipation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnableParticipationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).EnableParticipation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_EnableParticipation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).EnableParticipation(ctx, req.(*EnableParticipationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_DisableParticipation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableParticipationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).DisableParticipation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_DisableParticipation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).DisableParticipation(ctx, req.(*DisableParticipationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_GetParticipatingStores_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetParticipatingStoresRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).GetParticipatingStores(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_GetParticipatingStores_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).GetParticipatingStores(ctx, req.(*GetParticipatingStoresRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_AddProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).AddProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_AddProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).AddProduct(ctx, req.(*AddProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_RemoveProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).RemoveProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_RemoveProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).RemoveProduct(ctx, req.(*RemoveProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_GetProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).GetProduct(ctx, req.(*GetProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StoresService_GetCatalog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCatalogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoresServiceServer).GetCatalog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StoresService_GetCatalog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoresServiceServer).GetCatalog(ctx, req.(*GetCatalogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StoresService_ServiceDesc is the grpc.ServiceDesc for StoresService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StoresService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storespb.StoresService",
	HandlerType: (*StoresServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateStore",
			Handler:    _StoresService_CreateStore_Handler,
		},
		{
			MethodName: "GetStore",
			Handler:    _StoresService_GetStore_Handler,
		},
		{
			MethodName: "GetStores",
			Handler:    _StoresService_GetStores_Handler,
		},
		{
			MethodName: "EnableParticipation",
			Handler:    _StoresService_EnableParticipation_Handler,
		},
		{
			MethodName: "DisableParticipation",
			Handler:    _StoresService_DisableParticipation_Handler,
		},
		{
			MethodName: "GetParticipatingStores",
			Handler:    _StoresService_GetParticipatingStores_Handler,
		},
		{
			MethodName: "AddProduct",
			Handler:    _StoresService_AddProduct_Handler,
		},
		{
			MethodName: "RemoveProduct",
			Handler:    _StoresService_RemoveProduct_Handler,
		},
		{
			MethodName: "GetProduct",
			Handler:    _StoresService_GetProduct_Handler,
		},
		{
			MethodName: "GetCatalog",
			Handler:    _StoresService_GetCatalog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api_stores.proto",
}
