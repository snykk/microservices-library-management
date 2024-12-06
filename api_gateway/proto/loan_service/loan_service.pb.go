// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: loan_service.proto

package loan_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Loan message to represent loan data
type Loan struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId     string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	BookId     string `protobuf:"bytes,3,opt,name=book_id,json=bookId,proto3" json:"book_id,omitempty"`
	LoanDate   int64  `protobuf:"varint,4,opt,name=loan_date,json=loanDate,proto3" json:"loan_date,omitempty"`       // unix time
	ReturnDate int64  `protobuf:"varint,5,opt,name=return_date,json=returnDate,proto3" json:"return_date,omitempty"` // unix time
	Status     string `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`                            // Loan status (e.g., BORROWED, RETURNED, LOST)
	CreatedAt  int64  `protobuf:"varint,7,opt,name=createdAt,proto3" json:"createdAt,omitempty"`                     // unix time
	UpdatedAt  int64  `protobuf:"varint,8,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`                     // unit time
}

func (x *Loan) Reset() {
	*x = Loan{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Loan) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Loan) ProtoMessage() {}

func (x *Loan) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Loan.ProtoReflect.Descriptor instead.
func (*Loan) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{0}
}

func (x *Loan) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Loan) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Loan) GetBookId() string {
	if x != nil {
		return x.BookId
	}
	return ""
}

func (x *Loan) GetLoanDate() int64 {
	if x != nil {
		return x.LoanDate
	}
	return 0
}

func (x *Loan) GetReturnDate() int64 {
	if x != nil {
		return x.ReturnDate
	}
	return 0
}

func (x *Loan) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Loan) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Loan) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type CreateLoanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	BookId string `protobuf:"bytes,2,opt,name=book_id,json=bookId,proto3" json:"book_id,omitempty"`
	Email  string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"` // for sending loan notification to the email user
}

func (x *CreateLoanRequest) Reset() {
	*x = CreateLoanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLoanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLoanRequest) ProtoMessage() {}

func (x *CreateLoanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLoanRequest.ProtoReflect.Descriptor instead.
func (*CreateLoanRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateLoanRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *CreateLoanRequest) GetBookId() string {
	if x != nil {
		return x.BookId
	}
	return ""
}

func (x *CreateLoanRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type ReturnLoanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId     string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Email      string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"` // for sending loan notification to the email user
	ReturnDate int64  `protobuf:"varint,4,opt,name=return_date,json=returnDate,proto3" json:"return_date,omitempty"`
}

func (x *ReturnLoanRequest) Reset() {
	*x = ReturnLoanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReturnLoanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReturnLoanRequest) ProtoMessage() {}

func (x *ReturnLoanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReturnLoanRequest.ProtoReflect.Descriptor instead.
func (*ReturnLoanRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{2}
}

func (x *ReturnLoanRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ReturnLoanRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ReturnLoanRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ReturnLoanRequest) GetReturnDate() int64 {
	if x != nil {
		return x.ReturnDate
	}
	return 0
}

type GetLoanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetLoanRequest) Reset() {
	*x = GetLoanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLoanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLoanRequest) ProtoMessage() {}

func (x *GetLoanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLoanRequest.ProtoReflect.Descriptor instead.
func (*GetLoanRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetLoanRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpdateLoanStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status     string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`                            // New loan status (e.g., RETURNED, LOST)
	ReturnDate int64  `protobuf:"varint,3,opt,name=return_date,json=returnDate,proto3" json:"return_date,omitempty"` // unix time
}

func (x *UpdateLoanStatusRequest) Reset() {
	*x = UpdateLoanStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLoanStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLoanStatusRequest) ProtoMessage() {}

func (x *UpdateLoanStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLoanStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdateLoanStatusRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateLoanStatusRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateLoanStatusRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *UpdateLoanStatusRequest) GetReturnDate() int64 {
	if x != nil {
		return x.ReturnDate
	}
	return 0
}

type ListUserLoansRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // User ID
}

func (x *ListUserLoansRequest) Reset() {
	*x = ListUserLoansRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUserLoansRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserLoansRequest) ProtoMessage() {}

func (x *ListUserLoansRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserLoansRequest.ProtoReflect.Descriptor instead.
func (*ListUserLoansRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{5}
}

func (x *ListUserLoansRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ListLoansRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListLoansRequest) Reset() {
	*x = ListLoansRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLoansRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLoansRequest) ProtoMessage() {}

func (x *ListLoansRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLoansRequest.ProtoReflect.Descriptor instead.
func (*ListLoansRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{6}
}

type LoanResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Loan *Loan `protobuf:"bytes,1,opt,name=loan,proto3" json:"loan,omitempty"`
}

func (x *LoanResponse) Reset() {
	*x = LoanResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoanResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoanResponse) ProtoMessage() {}

func (x *LoanResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoanResponse.ProtoReflect.Descriptor instead.
func (*LoanResponse) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{7}
}

func (x *LoanResponse) GetLoan() *Loan {
	if x != nil {
		return x.Loan
	}
	return nil
}

type ListLoansResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Loans []*Loan `protobuf:"bytes,1,rep,name=loans,proto3" json:"loans,omitempty"` // List of loans
}

func (x *ListLoansResponse) Reset() {
	*x = ListLoansResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLoansResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLoansResponse) ProtoMessage() {}

func (x *ListLoansResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLoansResponse.ProtoReflect.Descriptor instead.
func (*ListLoansResponse) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{8}
}

func (x *ListLoansResponse) GetLoans() []*Loan {
	if x != nil {
		return x.Loans
	}
	return nil
}

type GetUserLoansByStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"` // Loan status to filter by (e.g., BORROWED, RETURNED, LOST)
}

func (x *GetUserLoansByStatusRequest) Reset() {
	*x = GetUserLoansByStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserLoansByStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserLoansByStatusRequest) ProtoMessage() {}

func (x *GetUserLoansByStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserLoansByStatusRequest.ProtoReflect.Descriptor instead.
func (*GetUserLoansByStatusRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{9}
}

func (x *GetUserLoansByStatusRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetUserLoansByStatusRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type GetLoansByStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"` // Loan status to filter by (e.g., BORROWED, RETURNED, LOST)
}

func (x *GetLoansByStatusRequest) Reset() {
	*x = GetLoansByStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loan_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLoansByStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLoansByStatusRequest) ProtoMessage() {}

func (x *GetLoansByStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loan_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLoansByStatusRequest.ProtoReflect.Descriptor instead.
func (*GetLoansByStatusRequest) Descriptor() ([]byte, []int) {
	return file_loan_service_proto_rawDescGZIP(), []int{10}
}

func (x *GetLoansByStatusRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_loan_service_proto protoreflect.FileDescriptor

var file_loan_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x22, 0xda, 0x01, 0x0a, 0x04, 0x4c, 0x6f, 0x61, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x6c, 0x6f, 0x61, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22,
	0x5b, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a,
	0x07, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x62, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x73, 0x0a, 0x11,
	0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x44, 0x61, 0x74,
	0x65, 0x22, 0x20, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x62, 0x0a, 0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x61,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e,
	0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x65, 0x74,
	0x75, 0x72, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x22, 0x2f, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x12, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74,
	0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x36, 0x0a, 0x0c,
	0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x04,
	0x6c, 0x6f, 0x61, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6c, 0x6f, 0x61,
	0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x04,
	0x6c, 0x6f, 0x61, 0x6e, 0x22, 0x3d, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x61, 0x6e,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x05, 0x6c, 0x6f, 0x61,
	0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x05, 0x6c, 0x6f,
	0x61, 0x6e, 0x73, 0x22, 0x4e, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4c, 0x6f,
	0x61, 0x6e, 0x73, 0x42, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x31, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x42,
	0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0xa3, 0x05, 0x0a, 0x0b, 0x4c, 0x6f, 0x61, 0x6e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4c, 0x6f, 0x61, 0x6e, 0x12, 0x1f, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x49, 0x0a, 0x0a, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x4c, 0x6f, 0x61, 0x6e, 0x12,
	0x1f, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52,
	0x65, 0x74, 0x75, 0x72, 0x6e, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1a, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x12, 0x1c, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x55, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x61, 0x6e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x61, 0x6e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6c,
	0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x61, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x12, 0x22, 0x2e, 0x6c, 0x6f, 0x61, 0x6e,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c,
	0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x12, 0x1e, 0x2e, 0x6c, 0x6f,
	0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c,
	0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6c, 0x6f,
	0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c,
	0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x14,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x42, 0x79, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x29, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x6e, 0x73,
	0x42, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1f, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x5a, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x42, 0x79, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x2e, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x61, 0x6e, 0x73, 0x42, 0x79, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6c, 0x6f,
	0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c,
	0x6f, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d,
	0x2f, 0x6c, 0x6f, 0x61, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_loan_service_proto_rawDescOnce sync.Once
	file_loan_service_proto_rawDescData = file_loan_service_proto_rawDesc
)

func file_loan_service_proto_rawDescGZIP() []byte {
	file_loan_service_proto_rawDescOnce.Do(func() {
		file_loan_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_loan_service_proto_rawDescData)
	})
	return file_loan_service_proto_rawDescData
}

var file_loan_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_loan_service_proto_goTypes = []interface{}{
	(*Loan)(nil),                        // 0: loan_service.Loan
	(*CreateLoanRequest)(nil),           // 1: loan_service.CreateLoanRequest
	(*ReturnLoanRequest)(nil),           // 2: loan_service.ReturnLoanRequest
	(*GetLoanRequest)(nil),              // 3: loan_service.GetLoanRequest
	(*UpdateLoanStatusRequest)(nil),     // 4: loan_service.UpdateLoanStatusRequest
	(*ListUserLoansRequest)(nil),        // 5: loan_service.ListUserLoansRequest
	(*ListLoansRequest)(nil),            // 6: loan_service.ListLoansRequest
	(*LoanResponse)(nil),                // 7: loan_service.LoanResponse
	(*ListLoansResponse)(nil),           // 8: loan_service.ListLoansResponse
	(*GetUserLoansByStatusRequest)(nil), // 9: loan_service.GetUserLoansByStatusRequest
	(*GetLoansByStatusRequest)(nil),     // 10: loan_service.GetLoansByStatusRequest
}
var file_loan_service_proto_depIdxs = []int32{
	0,  // 0: loan_service.LoanResponse.loan:type_name -> loan_service.Loan
	0,  // 1: loan_service.ListLoansResponse.loans:type_name -> loan_service.Loan
	1,  // 2: loan_service.LoanService.CreateLoan:input_type -> loan_service.CreateLoanRequest
	2,  // 3: loan_service.LoanService.ReturnLoan:input_type -> loan_service.ReturnLoanRequest
	3,  // 4: loan_service.LoanService.GetLoan:input_type -> loan_service.GetLoanRequest
	4,  // 5: loan_service.LoanService.UpdateLoanStatus:input_type -> loan_service.UpdateLoanStatusRequest
	5,  // 6: loan_service.LoanService.ListUserLoans:input_type -> loan_service.ListUserLoansRequest
	6,  // 7: loan_service.LoanService.ListLoans:input_type -> loan_service.ListLoansRequest
	9,  // 8: loan_service.LoanService.GetUserLoansByStatus:input_type -> loan_service.GetUserLoansByStatusRequest
	10, // 9: loan_service.LoanService.GetLoansByStatus:input_type -> loan_service.GetLoansByStatusRequest
	7,  // 10: loan_service.LoanService.CreateLoan:output_type -> loan_service.LoanResponse
	7,  // 11: loan_service.LoanService.ReturnLoan:output_type -> loan_service.LoanResponse
	7,  // 12: loan_service.LoanService.GetLoan:output_type -> loan_service.LoanResponse
	7,  // 13: loan_service.LoanService.UpdateLoanStatus:output_type -> loan_service.LoanResponse
	8,  // 14: loan_service.LoanService.ListUserLoans:output_type -> loan_service.ListLoansResponse
	8,  // 15: loan_service.LoanService.ListLoans:output_type -> loan_service.ListLoansResponse
	8,  // 16: loan_service.LoanService.GetUserLoansByStatus:output_type -> loan_service.ListLoansResponse
	8,  // 17: loan_service.LoanService.GetLoansByStatus:output_type -> loan_service.ListLoansResponse
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_loan_service_proto_init() }
func file_loan_service_proto_init() {
	if File_loan_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_loan_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Loan); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLoanRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReturnLoanRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLoanRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateLoanStatusRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUserLoansRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLoansRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoanResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLoansResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserLoansByStatusRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_loan_service_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLoansByStatusRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_loan_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_loan_service_proto_goTypes,
		DependencyIndexes: file_loan_service_proto_depIdxs,
		MessageInfos:      file_loan_service_proto_msgTypes,
	}.Build()
	File_loan_service_proto = out.File
	file_loan_service_proto_rawDesc = nil
	file_loan_service_proto_goTypes = nil
	file_loan_service_proto_depIdxs = nil
}
