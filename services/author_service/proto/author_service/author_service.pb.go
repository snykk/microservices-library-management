// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: author_service.proto

package author_service

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type Author struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Biography string `protobuf:"bytes,3,opt,name=biography,proto3" json:"biography,omitempty"`
	Version   int32  `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty"`
	CreatedAt int64  `protobuf:"varint,5,opt,name=createdAt,proto3" json:"createdAt,omitempty"` // unix time
	UpdatedAt int64  `protobuf:"varint,6,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"` // unit time
}

func (x *Author) Reset() {
	*x = Author{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Author) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Author) ProtoMessage() {}

func (x *Author) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Author.ProtoReflect.Descriptor instead.
func (*Author) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{0}
}

func (x *Author) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Author) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Author) GetBiography() string {
	if x != nil {
		return x.Biography
	}
	return ""
}

func (x *Author) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Author) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Author) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type CreateAuthorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Biography string `protobuf:"bytes,2,opt,name=biography,proto3" json:"biography,omitempty"`
}

func (x *CreateAuthorRequest) Reset() {
	*x = CreateAuthorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAuthorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAuthorRequest) ProtoMessage() {}

func (x *CreateAuthorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAuthorRequest.ProtoReflect.Descriptor instead.
func (*CreateAuthorRequest) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateAuthorRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateAuthorRequest) GetBiography() string {
	if x != nil {
		return x.Biography
	}
	return ""
}

type CreateAuthorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Author *Author `protobuf:"bytes,1,opt,name=author,proto3" json:"author,omitempty"`
}

func (x *CreateAuthorResponse) Reset() {
	*x = CreateAuthorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAuthorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAuthorResponse) ProtoMessage() {}

func (x *CreateAuthorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAuthorResponse.ProtoReflect.Descriptor instead.
func (*CreateAuthorResponse) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{2}
}

func (x *CreateAuthorResponse) GetAuthor() *Author {
	if x != nil {
		return x.Author
	}
	return nil
}

type GetAuthorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // ID must not be empty
}

func (x *GetAuthorRequest) Reset() {
	*x = GetAuthorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthorRequest) ProtoMessage() {}

func (x *GetAuthorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthorRequest.ProtoReflect.Descriptor instead.
func (*GetAuthorRequest) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetAuthorRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetAuthorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Author *Author `protobuf:"bytes,1,opt,name=author,proto3" json:"author,omitempty"`
}

func (x *GetAuthorResponse) Reset() {
	*x = GetAuthorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAuthorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAuthorResponse) ProtoMessage() {}

func (x *GetAuthorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAuthorResponse.ProtoReflect.Descriptor instead.
func (*GetAuthorResponse) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetAuthorResponse) GetAuthor() *Author {
	if x != nil {
		return x.Author
	}
	return nil
}

type ListAuthorsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page     int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`         // Page must be >= 1
	PageSize int32 `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"` // Page size must be >= 1
}

func (x *ListAuthorsRequest) Reset() {
	*x = ListAuthorsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAuthorsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAuthorsRequest) ProtoMessage() {}

func (x *ListAuthorsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAuthorsRequest.ProtoReflect.Descriptor instead.
func (*ListAuthorsRequest) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{5}
}

func (x *ListAuthorsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListAuthorsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ListAuthorsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authors    []*Author `protobuf:"bytes,1,rep,name=authors,proto3" json:"authors,omitempty"`
	TotalItems int32     `protobuf:"varint,2,opt,name=totalItems,proto3" json:"totalItems,omitempty"` // Total number of items
	TotalPages int32     `protobuf:"varint,3,opt,name=totalPages,proto3" json:"totalPages,omitempty"` // Total number of pages
}

func (x *ListAuthorsResponse) Reset() {
	*x = ListAuthorsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAuthorsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAuthorsResponse) ProtoMessage() {}

func (x *ListAuthorsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAuthorsResponse.ProtoReflect.Descriptor instead.
func (*ListAuthorsResponse) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{6}
}

func (x *ListAuthorsResponse) GetAuthors() []*Author {
	if x != nil {
		return x.Authors
	}
	return nil
}

func (x *ListAuthorsResponse) GetTotalItems() int32 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

func (x *ListAuthorsResponse) GetTotalPages() int32 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

type UpdateAuthorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // ID must not be empty
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Biography string `protobuf:"bytes,3,opt,name=biography,proto3" json:"biography,omitempty"`
	Version   int32  `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *UpdateAuthorRequest) Reset() {
	*x = UpdateAuthorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAuthorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAuthorRequest) ProtoMessage() {}

func (x *UpdateAuthorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAuthorRequest.ProtoReflect.Descriptor instead.
func (*UpdateAuthorRequest) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateAuthorRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateAuthorRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateAuthorRequest) GetBiography() string {
	if x != nil {
		return x.Biography
	}
	return ""
}

func (x *UpdateAuthorRequest) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

type UpdateAuthorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Author *Author `protobuf:"bytes,1,opt,name=author,proto3" json:"author,omitempty"`
}

func (x *UpdateAuthorResponse) Reset() {
	*x = UpdateAuthorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAuthorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAuthorResponse) ProtoMessage() {}

func (x *UpdateAuthorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAuthorResponse.ProtoReflect.Descriptor instead.
func (*UpdateAuthorResponse) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateAuthorResponse) GetAuthor() *Author {
	if x != nil {
		return x.Author
	}
	return nil
}

type DeleteAuthorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // ID must not be empty
	Version int32  `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *DeleteAuthorRequest) Reset() {
	*x = DeleteAuthorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAuthorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAuthorRequest) ProtoMessage() {}

func (x *DeleteAuthorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAuthorRequest.ProtoReflect.Descriptor instead.
func (*DeleteAuthorRequest) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteAuthorRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DeleteAuthorRequest) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

type DeleteAuthorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *DeleteAuthorResponse) Reset() {
	*x = DeleteAuthorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_author_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAuthorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAuthorResponse) ProtoMessage() {}

func (x *DeleteAuthorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_author_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAuthorResponse.ProtoReflect.Descriptor instead.
func (*DeleteAuthorResponse) Descriptor() ([]byte, []int) {
	return file_author_service_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteAuthorResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_author_service_proto protoreflect.FileDescriptor

var file_author_service_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xa0, 0x01, 0x0a, 0x06, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x22, 0x59, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x03,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61,
	0x70, 0x68, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x10, 0x01, 0x52, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x22, 0x46, 0x0a,
	0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x06, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x22, 0x2b, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x43, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52,
	0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x22, 0x56, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x1a, 0x02, 0x28, 0x01, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x08, 0x70, 0x61,
	0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x1a, 0x02, 0x28, 0x01, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22,
	0x87, 0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x07, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x22, 0x95, 0x01, 0x0a, 0x13, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10,
	0x03, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x09, 0x62, 0x69, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x79, 0x12, 0x21,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x1a, 0x02, 0x28, 0x01, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x22, 0x46, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x22, 0x51, 0x0a, 0x13, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x1a,
	0x02, 0x28, 0x01, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x30, 0x0a, 0x14,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xca,
	0x03, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x59, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x12, 0x23, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x20, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a,
	0x0b, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x12, 0x22, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x23, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x23, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x59, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x12, 0x23, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x2f,
	0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_author_service_proto_rawDescOnce sync.Once
	file_author_service_proto_rawDescData = file_author_service_proto_rawDesc
)

func file_author_service_proto_rawDescGZIP() []byte {
	file_author_service_proto_rawDescOnce.Do(func() {
		file_author_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_author_service_proto_rawDescData)
	})
	return file_author_service_proto_rawDescData
}

var file_author_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_author_service_proto_goTypes = []interface{}{
	(*Author)(nil),               // 0: author_service.Author
	(*CreateAuthorRequest)(nil),  // 1: author_service.CreateAuthorRequest
	(*CreateAuthorResponse)(nil), // 2: author_service.CreateAuthorResponse
	(*GetAuthorRequest)(nil),     // 3: author_service.GetAuthorRequest
	(*GetAuthorResponse)(nil),    // 4: author_service.GetAuthorResponse
	(*ListAuthorsRequest)(nil),   // 5: author_service.ListAuthorsRequest
	(*ListAuthorsResponse)(nil),  // 6: author_service.ListAuthorsResponse
	(*UpdateAuthorRequest)(nil),  // 7: author_service.UpdateAuthorRequest
	(*UpdateAuthorResponse)(nil), // 8: author_service.UpdateAuthorResponse
	(*DeleteAuthorRequest)(nil),  // 9: author_service.DeleteAuthorRequest
	(*DeleteAuthorResponse)(nil), // 10: author_service.DeleteAuthorResponse
}
var file_author_service_proto_depIdxs = []int32{
	0,  // 0: author_service.CreateAuthorResponse.author:type_name -> author_service.Author
	0,  // 1: author_service.GetAuthorResponse.author:type_name -> author_service.Author
	0,  // 2: author_service.ListAuthorsResponse.authors:type_name -> author_service.Author
	0,  // 3: author_service.UpdateAuthorResponse.author:type_name -> author_service.Author
	1,  // 4: author_service.AuthorService.CreateAuthor:input_type -> author_service.CreateAuthorRequest
	3,  // 5: author_service.AuthorService.GetAuthor:input_type -> author_service.GetAuthorRequest
	5,  // 6: author_service.AuthorService.ListAuthors:input_type -> author_service.ListAuthorsRequest
	7,  // 7: author_service.AuthorService.UpdateAuthor:input_type -> author_service.UpdateAuthorRequest
	9,  // 8: author_service.AuthorService.DeleteAuthor:input_type -> author_service.DeleteAuthorRequest
	2,  // 9: author_service.AuthorService.CreateAuthor:output_type -> author_service.CreateAuthorResponse
	4,  // 10: author_service.AuthorService.GetAuthor:output_type -> author_service.GetAuthorResponse
	6,  // 11: author_service.AuthorService.ListAuthors:output_type -> author_service.ListAuthorsResponse
	8,  // 12: author_service.AuthorService.UpdateAuthor:output_type -> author_service.UpdateAuthorResponse
	10, // 13: author_service.AuthorService.DeleteAuthor:output_type -> author_service.DeleteAuthorResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_author_service_proto_init() }
func file_author_service_proto_init() {
	if File_author_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_author_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Author); i {
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
		file_author_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAuthorRequest); i {
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
		file_author_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAuthorResponse); i {
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
		file_author_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthorRequest); i {
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
		file_author_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAuthorResponse); i {
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
		file_author_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAuthorsRequest); i {
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
		file_author_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAuthorsResponse); i {
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
		file_author_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAuthorRequest); i {
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
		file_author_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAuthorResponse); i {
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
		file_author_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAuthorRequest); i {
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
		file_author_service_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAuthorResponse); i {
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
			RawDescriptor: file_author_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_author_service_proto_goTypes,
		DependencyIndexes: file_author_service_proto_depIdxs,
		MessageInfos:      file_author_service_proto_msgTypes,
	}.Build()
	File_author_service_proto = out.File
	file_author_service_proto_rawDesc = nil
	file_author_service_proto_goTypes = nil
	file_author_service_proto_depIdxs = nil
}
