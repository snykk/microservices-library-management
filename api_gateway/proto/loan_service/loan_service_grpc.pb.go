// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: loan_service.proto

package loan_service

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

// LoanServiceClient is the client API for LoanService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoanServiceClient interface {
	CreateLoan(ctx context.Context, in *CreateLoanRequest, opts ...grpc.CallOption) (*LoanResponse, error)
	GetLoan(ctx context.Context, in *GetLoanRequest, opts ...grpc.CallOption) (*LoanResponse, error)
	UpdateLoanStatus(ctx context.Context, in *UpdateLoanStatusRequest, opts ...grpc.CallOption) (*LoanResponse, error)
	ListUserLoans(ctx context.Context, in *ListUserLoansRequest, opts ...grpc.CallOption) (*ListLoansResponse, error)
	ListLoans(ctx context.Context, in *ListLoansRequest, opts ...grpc.CallOption) (*ListLoansResponse, error)
}

type loanServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLoanServiceClient(cc grpc.ClientConnInterface) LoanServiceClient {
	return &loanServiceClient{cc}
}

func (c *loanServiceClient) CreateLoan(ctx context.Context, in *CreateLoanRequest, opts ...grpc.CallOption) (*LoanResponse, error) {
	out := new(LoanResponse)
	err := c.cc.Invoke(ctx, "/loan_service.LoanService/CreateLoan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loanServiceClient) GetLoan(ctx context.Context, in *GetLoanRequest, opts ...grpc.CallOption) (*LoanResponse, error) {
	out := new(LoanResponse)
	err := c.cc.Invoke(ctx, "/loan_service.LoanService/GetLoan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loanServiceClient) UpdateLoanStatus(ctx context.Context, in *UpdateLoanStatusRequest, opts ...grpc.CallOption) (*LoanResponse, error) {
	out := new(LoanResponse)
	err := c.cc.Invoke(ctx, "/loan_service.LoanService/UpdateLoanStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loanServiceClient) ListUserLoans(ctx context.Context, in *ListUserLoansRequest, opts ...grpc.CallOption) (*ListLoansResponse, error) {
	out := new(ListLoansResponse)
	err := c.cc.Invoke(ctx, "/loan_service.LoanService/ListUserLoans", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loanServiceClient) ListLoans(ctx context.Context, in *ListLoansRequest, opts ...grpc.CallOption) (*ListLoansResponse, error) {
	out := new(ListLoansResponse)
	err := c.cc.Invoke(ctx, "/loan_service.LoanService/ListLoans", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoanServiceServer is the server API for LoanService service.
// All implementations must embed UnimplementedLoanServiceServer
// for forward compatibility
type LoanServiceServer interface {
	CreateLoan(context.Context, *CreateLoanRequest) (*LoanResponse, error)
	GetLoan(context.Context, *GetLoanRequest) (*LoanResponse, error)
	UpdateLoanStatus(context.Context, *UpdateLoanStatusRequest) (*LoanResponse, error)
	ListUserLoans(context.Context, *ListUserLoansRequest) (*ListLoansResponse, error)
	ListLoans(context.Context, *ListLoansRequest) (*ListLoansResponse, error)
	mustEmbedUnimplementedLoanServiceServer()
}

// UnimplementedLoanServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLoanServiceServer struct {
}

func (UnimplementedLoanServiceServer) CreateLoan(context.Context, *CreateLoanRequest) (*LoanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLoan not implemented")
}
func (UnimplementedLoanServiceServer) GetLoan(context.Context, *GetLoanRequest) (*LoanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLoan not implemented")
}
func (UnimplementedLoanServiceServer) UpdateLoanStatus(context.Context, *UpdateLoanStatusRequest) (*LoanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLoanStatus not implemented")
}
func (UnimplementedLoanServiceServer) ListUserLoans(context.Context, *ListUserLoansRequest) (*ListLoansResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserLoans not implemented")
}
func (UnimplementedLoanServiceServer) ListLoans(context.Context, *ListLoansRequest) (*ListLoansResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLoans not implemented")
}
func (UnimplementedLoanServiceServer) mustEmbedUnimplementedLoanServiceServer() {}

// UnsafeLoanServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoanServiceServer will
// result in compilation errors.
type UnsafeLoanServiceServer interface {
	mustEmbedUnimplementedLoanServiceServer()
}

func RegisterLoanServiceServer(s grpc.ServiceRegistrar, srv LoanServiceServer) {
	s.RegisterService(&LoanService_ServiceDesc, srv)
}

func _LoanService_CreateLoan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLoanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoanServiceServer).CreateLoan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loan_service.LoanService/CreateLoan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoanServiceServer).CreateLoan(ctx, req.(*CreateLoanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoanService_GetLoan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLoanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoanServiceServer).GetLoan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loan_service.LoanService/GetLoan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoanServiceServer).GetLoan(ctx, req.(*GetLoanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoanService_UpdateLoanStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLoanStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoanServiceServer).UpdateLoanStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loan_service.LoanService/UpdateLoanStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoanServiceServer).UpdateLoanStatus(ctx, req.(*UpdateLoanStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoanService_ListUserLoans_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserLoansRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoanServiceServer).ListUserLoans(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loan_service.LoanService/ListUserLoans",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoanServiceServer).ListUserLoans(ctx, req.(*ListUserLoansRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoanService_ListLoans_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLoansRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoanServiceServer).ListLoans(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loan_service.LoanService/ListLoans",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoanServiceServer).ListLoans(ctx, req.(*ListLoansRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LoanService_ServiceDesc is the grpc.ServiceDesc for LoanService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoanService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loan_service.LoanService",
	HandlerType: (*LoanServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLoan",
			Handler:    _LoanService_CreateLoan_Handler,
		},
		{
			MethodName: "GetLoan",
			Handler:    _LoanService_GetLoan_Handler,
		},
		{
			MethodName: "UpdateLoanStatus",
			Handler:    _LoanService_UpdateLoanStatus_Handler,
		},
		{
			MethodName: "ListUserLoans",
			Handler:    _LoanService_ListUserLoans_Handler,
		},
		{
			MethodName: "ListLoans",
			Handler:    _LoanService_ListLoans_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "loan_service.proto",
}
