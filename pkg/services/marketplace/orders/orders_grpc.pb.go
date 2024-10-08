// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.2
// source: orders.proto

package orders

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
	OrderService_GetSellerOrders_FullMethodName = "/orders.OrderService/GetSellerOrders"
	OrderService_GetUserOrders_FullMethodName   = "/orders.OrderService/GetUserOrders"
	OrderService_CreateOrder_FullMethodName     = "/orders.OrderService/CreateOrder"
	OrderService_GetOrder_FullMethodName        = "/orders.OrderService/GetOrder"
	OrderService_AcceptOrder_FullMethodName     = "/orders.OrderService/AcceptOrder"
	OrderService_CompleteOrder_FullMethodName   = "/orders.OrderService/CompleteOrder"
	OrderService_CheckUserBlock_FullMethodName  = "/orders.OrderService/CheckUserBlock"
	OrderService_BlockUser_FullMethodName       = "/orders.OrderService/BlockUser"
	OrderService_UnblockUser_FullMethodName     = "/orders.OrderService/UnblockUser"
)

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderServiceClient interface {
	GetSellerOrders(ctx context.Context, in *Seller, opts ...grpc.CallOption) (*AllOrders, error)
	GetUserOrders(ctx context.Context, in *Buyer, opts ...grpc.CallOption) (*AllOrders, error)
	CreateOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*OrderID, error)
	GetOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*Order, error)
	AcceptOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*Response, error)
	CompleteOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*Response, error)
	CheckUserBlock(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserBlock, error)
	BlockUser(ctx context.Context, in *UserBlock, opts ...grpc.CallOption) (*Response, error)
	UnblockUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*Response, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) GetSellerOrders(ctx context.Context, in *Seller, opts ...grpc.CallOption) (*AllOrders, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllOrders)
	err := c.cc.Invoke(ctx, OrderService_GetSellerOrders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetUserOrders(ctx context.Context, in *Buyer, opts ...grpc.CallOption) (*AllOrders, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AllOrders)
	err := c.cc.Invoke(ctx, OrderService_GetUserOrders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CreateOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*OrderID, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderID)
	err := c.cc.Invoke(ctx, OrderService_CreateOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*Order, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Order)
	err := c.cc.Invoke(ctx, OrderService_GetOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) AcceptOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, OrderService_AcceptOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CompleteOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, OrderService_CompleteOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) CheckUserBlock(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserBlock, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserBlock)
	err := c.cc.Invoke(ctx, OrderService_CheckUserBlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) BlockUser(ctx context.Context, in *UserBlock, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, OrderService_BlockUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) UnblockUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, OrderService_UnblockUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServiceServer is the server API for OrderService service.
// All implementations must embed UnimplementedOrderServiceServer
// for forward compatibility.
type OrderServiceServer interface {
	GetSellerOrders(context.Context, *Seller) (*AllOrders, error)
	GetUserOrders(context.Context, *Buyer) (*AllOrders, error)
	CreateOrder(context.Context, *Order) (*OrderID, error)
	GetOrder(context.Context, *OrderID) (*Order, error)
	AcceptOrder(context.Context, *OrderID) (*Response, error)
	CompleteOrder(context.Context, *OrderID) (*Response, error)
	CheckUserBlock(context.Context, *User) (*UserBlock, error)
	BlockUser(context.Context, *UserBlock) (*Response, error)
	UnblockUser(context.Context, *User) (*Response, error)
	mustEmbedUnimplementedOrderServiceServer()
}

// UnimplementedOrderServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOrderServiceServer struct{}

func (UnimplementedOrderServiceServer) GetSellerOrders(context.Context, *Seller) (*AllOrders, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSellerOrders not implemented")
}
func (UnimplementedOrderServiceServer) GetUserOrders(context.Context, *Buyer) (*AllOrders, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserOrders not implemented")
}
func (UnimplementedOrderServiceServer) CreateOrder(context.Context, *Order) (*OrderID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedOrderServiceServer) GetOrder(context.Context, *OrderID) (*Order, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
}
func (UnimplementedOrderServiceServer) AcceptOrder(context.Context, *OrderID) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptOrder not implemented")
}
func (UnimplementedOrderServiceServer) CompleteOrder(context.Context, *OrderID) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteOrder not implemented")
}
func (UnimplementedOrderServiceServer) CheckUserBlock(context.Context, *User) (*UserBlock, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUserBlock not implemented")
}
func (UnimplementedOrderServiceServer) BlockUser(context.Context, *UserBlock) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlockUser not implemented")
}
func (UnimplementedOrderServiceServer) UnblockUser(context.Context, *User) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnblockUser not implemented")
}
func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}
func (UnimplementedOrderServiceServer) testEmbeddedByValue()                      {}

// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServiceServer will
// result in compilation errors.
type UnsafeOrderServiceServer interface {
	mustEmbedUnimplementedOrderServiceServer()
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	// If the following call pancis, it indicates UnimplementedOrderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_GetSellerOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Seller)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetSellerOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_GetSellerOrders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetSellerOrders(ctx, req.(*Seller))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_GetUserOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Buyer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetUserOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_GetUserOrders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetUserOrders(ctx, req.(*Buyer))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Order)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_CreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CreateOrder(ctx, req.(*Order))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_GetOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_GetOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetOrder(ctx, req.(*OrderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_AcceptOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).AcceptOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_AcceptOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).AcceptOrder(ctx, req.(*OrderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CompleteOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CompleteOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_CompleteOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CompleteOrder(ctx, req.(*OrderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_CheckUserBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CheckUserBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_CheckUserBlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CheckUserBlock(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_BlockUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserBlock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).BlockUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_BlockUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).BlockUser(ctx, req.(*UserBlock))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_UnblockUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).UnblockUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderService_UnblockUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).UnblockUser(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "orders.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSellerOrders",
			Handler:    _OrderService_GetSellerOrders_Handler,
		},
		{
			MethodName: "GetUserOrders",
			Handler:    _OrderService_GetUserOrders_Handler,
		},
		{
			MethodName: "CreateOrder",
			Handler:    _OrderService_CreateOrder_Handler,
		},
		{
			MethodName: "GetOrder",
			Handler:    _OrderService_GetOrder_Handler,
		},
		{
			MethodName: "AcceptOrder",
			Handler:    _OrderService_AcceptOrder_Handler,
		},
		{
			MethodName: "CompleteOrder",
			Handler:    _OrderService_CompleteOrder_Handler,
		},
		{
			MethodName: "CheckUserBlock",
			Handler:    _OrderService_CheckUserBlock_Handler,
		},
		{
			MethodName: "BlockUser",
			Handler:    _OrderService_BlockUser_Handler,
		},
		{
			MethodName: "UnblockUser",
			Handler:    _OrderService_UnblockUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "orders.proto",
}
