// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: Alicja Kwiecinska (kwiecinskaa@google.com) github: alicjakwie

// Hermes's interface with Cloudprober

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.11.4
// source: github.com/googleinterns/step224-2020/alicja/config/proto/interface.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// TODO (#31): Complete ExitCode enum in interface.proto
type ExitCode int32

const (
	ExitCode_SUCCESS ExitCode = 0
)

// Enum value maps for ExitCode.
var (
	ExitCode_name = map[int32]string{
		0: "SUCCESS",
	}
	ExitCode_value = map[string]int32{
		"SUCCESS": 0,
	}
)

func (x ExitCode) Enum() *ExitCode {
	p := new(ExitCode)
	*p = x
	return p
}

func (x ExitCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ExitCode) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_enumTypes[0].Descriptor()
}

func (ExitCode) Type() protoreflect.EnumType {
	return &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_enumTypes[0]
}

func (x ExitCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ExitCode.Descriptor instead.
func (ExitCode) EnumDescriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{0}
}

type HermesFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`         // REQUIRED field
	Target   *TargetDefinition `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`     // REQUIRED field for PutFileRequest, GetFileRequest and DeleteFileRequest OPTIONAL field for GetFileResponse
	Contents string            `protobuf:"bytes,3,opt,name=contents,proto3" json:"contents,omitempty"` // REQUIRED field for PutFileRequest OPTIONAL field for GetFileResponse
}

func (x *HermesFile) Reset() {
	*x = HermesFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HermesFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HermesFile) ProtoMessage() {}

func (x *HermesFile) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HermesFile.ProtoReflect.Descriptor instead.
func (*HermesFile) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{0}
}

func (x *HermesFile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *HermesFile) GetTarget() *TargetDefinition {
	if x != nil {
		return x.Target
	}
	return nil
}

func (x *HermesFile) GetContents() string {
	if x != nil {
		return x.Contents
	}
	return ""
}

// HermesFile (name, contents and target) need to be specified so that Hermes knows what file to create, and which storage system to create it on
type PutFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *HermesFile `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *PutFileRequest) Reset() {
	*x = PutFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutFileRequest) ProtoMessage() {}

func (x *PutFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutFileRequest.ProtoReflect.Descriptor instead.
func (*PutFileRequest) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{1}
}

func (x *PutFileRequest) GetFile() *HermesFile {
	if x != nil {
		return x.File
	}
	return nil
}

// returns the exit code after the rpc GetFile
type PutFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExitCode ExitCode `protobuf:"varint,1,opt,name=exit_code,json=exitCode,proto3,enum=hermes.ExitCode" json:"exit_code,omitempty"` // REQUIRED field
}

func (x *PutFileResponse) Reset() {
	*x = PutFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutFileResponse) ProtoMessage() {}

func (x *PutFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutFileResponse.ProtoReflect.Descriptor instead.
func (*PutFileResponse) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{2}
}

func (x *PutFileResponse) GetExitCode() ExitCode {
	if x != nil {
		return x.ExitCode
	}
	return ExitCode_SUCCESS
}

// HermesFile (name and target, contents optional) need to be specified so that Hermes knows what file to retrieve and from which storage system to retrieve it from
type GetFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *HermesFile `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *GetFileRequest) Reset() {
	*x = GetFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileRequest) ProtoMessage() {}

func (x *GetFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileRequest.ProtoReflect.Descriptor instead.
func (*GetFileRequest) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{3}
}

func (x *GetFileRequest) GetFile() *HermesFile {
	if x != nil {
		return x.File
	}
	return nil
}

// returns the exit code after the rpc GetFile if GetFile suceeded we return the file (name and contents required) and if GetFile fails we do not
type GetFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExitCode ExitCode    `protobuf:"varint,1,opt,name=exit_code,json=exitCode,proto3,enum=hermes.ExitCode" json:"exit_code,omitempty"` // REQUIRED field
	File     *HermesFile `protobuf:"bytes,2,opt,name=file,proto3" json:"file,omitempty"`                                               // OPTIONAL field
}

func (x *GetFileResponse) Reset() {
	*x = GetFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileResponse) ProtoMessage() {}

func (x *GetFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileResponse.ProtoReflect.Descriptor instead.
func (*GetFileResponse) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{4}
}

func (x *GetFileResponse) GetExitCode() ExitCode {
	if x != nil {
		return x.ExitCode
	}
	return ExitCode_SUCCESS
}

func (x *GetFileResponse) GetFile() *HermesFile {
	if x != nil {
		return x.File
	}
	return nil
}

// HermesFile (name and target, contents optional) need to be specified so that Hermes knows what file to delete and from which storage system to delete it from
type DeleteFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *HermesFile `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"` // REQUIRED field
}

func (x *DeleteFileRequest) Reset() {
	*x = DeleteFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileRequest) ProtoMessage() {}

func (x *DeleteFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileRequest.ProtoReflect.Descriptor instead.
func (*DeleteFileRequest) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteFileRequest) GetFile() *HermesFile {
	if x != nil {
		return x.File
	}
	return nil
}

// returns the exit code after the rpc DeleteFile
type DeleteFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExitCode ExitCode `protobuf:"varint,1,opt,name=exit_code,json=exitCode,proto3,enum=hermes.ExitCode" json:"exit_code,omitempty"` // REQUIRED field
}

func (x *DeleteFileResponse) Reset() {
	*x = DeleteFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileResponse) ProtoMessage() {}

func (x *DeleteFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileResponse.ProtoReflect.Descriptor instead.
func (*DeleteFileResponse) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteFileResponse) GetExitCode() ExitCode {
	if x != nil {
		return x.ExitCode
	}
	return ExitCode_SUCCESS
}

var File_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto protoreflect.FileDescriptor

var file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDesc = []byte{
	0x0a, 0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x32,
	0x32, 0x34, 0x2d, 0x32, 0x30, 0x32, 0x30, 0x2f, 0x61, 0x6c, 0x69, 0x63, 0x6a, 0x61, 0x2f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x68, 0x65, 0x72,
	0x6d, 0x65, 0x73, 0x1a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74,
	0x65, 0x70, 0x32, 0x32, 0x34, 0x2d, 0x32, 0x30, 0x32, 0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6e, 0x0a, 0x0a, 0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x46,
	0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73,
	0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x38, 0x0a, 0x0e, 0x50, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x48,
	0x65, 0x72, 0x6d, 0x65, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22,
	0x40, 0x0a, 0x0f, 0x50, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2d, 0x0a, 0x09, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x45,
	0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64,
	0x65, 0x22, 0x38, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x48, 0x65, 0x72, 0x6d, 0x65,
	0x73, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x68, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d,
	0x0a, 0x09, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x10, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x45, 0x78, 0x69, 0x74, 0x43,
	0x6f, 0x64, 0x65, 0x52, 0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x26, 0x0a,
	0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x68, 0x65,
	0x72, 0x6d, 0x65, 0x73, 0x2e, 0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x3b, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x66, 0x69,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65,
	0x73, 0x2e, 0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69,
	0x6c, 0x65, 0x22, 0x43, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x09, 0x65, 0x78, 0x69, 0x74,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x68, 0x65,
	0x72, 0x6d, 0x65, 0x73, 0x2e, 0x45, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x08, 0x65,
	0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x2a, 0x17, 0x0a, 0x08, 0x45, 0x78, 0x69, 0x74, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00,
	0x32, 0xcb, 0x01, 0x0a, 0x0c, 0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x62, 0x65,
	0x72, 0x12, 0x3a, 0x0a, 0x07, 0x50, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x68,
	0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x50, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x50, 0x75,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a,
	0x07, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65,
	0x73, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x17, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0a, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x19, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x34,
	0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x32,
	0x32, 0x34, 0x2d, 0x32, 0x30, 0x32, 0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescOnce sync.Once
	file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescData = file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDesc
)

func file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescGZIP() []byte {
	file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescOnce.Do(func() {
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescData)
	})
	return file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDescData
}

var file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_goTypes = []interface{}{
	(ExitCode)(0),              // 0: hermes.ExitCode
	(*HermesFile)(nil),         // 1: hermes.HermesFile
	(*PutFileRequest)(nil),     // 2: hermes.PutFileRequest
	(*PutFileResponse)(nil),    // 3: hermes.PutFileResponse
	(*GetFileRequest)(nil),     // 4: hermes.GetFileRequest
	(*GetFileResponse)(nil),    // 5: hermes.GetFileResponse
	(*DeleteFileRequest)(nil),  // 6: hermes.DeleteFileRequest
	(*DeleteFileResponse)(nil), // 7: hermes.DeleteFileResponse
	(*TargetDefinition)(nil),   // 8: hermes.TargetDefinition
}
var file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_depIdxs = []int32{
	8,  // 0: hermes.HermesFile.target:type_name -> hermes.TargetDefinition
	1,  // 1: hermes.PutFileRequest.file:type_name -> hermes.HermesFile
	0,  // 2: hermes.PutFileResponse.exit_code:type_name -> hermes.ExitCode
	1,  // 3: hermes.GetFileRequest.file:type_name -> hermes.HermesFile
	0,  // 4: hermes.GetFileResponse.exit_code:type_name -> hermes.ExitCode
	1,  // 5: hermes.GetFileResponse.file:type_name -> hermes.HermesFile
	1,  // 6: hermes.DeleteFileRequest.file:type_name -> hermes.HermesFile
	0,  // 7: hermes.DeleteFileResponse.exit_code:type_name -> hermes.ExitCode
	2,  // 8: hermes.HermesProber.PutFile:input_type -> hermes.PutFileRequest
	4,  // 9: hermes.HermesProber.GetFile:input_type -> hermes.GetFileRequest
	6,  // 10: hermes.HermesProber.DeleteFile:input_type -> hermes.DeleteFileRequest
	3,  // 11: hermes.HermesProber.PutFile:output_type -> hermes.PutFileResponse
	5,  // 12: hermes.HermesProber.GetFile:output_type -> hermes.GetFileResponse
	7,  // 13: hermes.HermesProber.DeleteFile:output_type -> hermes.DeleteFileResponse
	11, // [11:14] is the sub-list for method output_type
	8,  // [8:11] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_init() }
func file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_init() {
	if File_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto != nil {
		return
	}
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HermesFile); i {
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
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutFileRequest); i {
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
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutFileResponse); i {
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
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileRequest); i {
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
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileResponse); i {
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
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFileRequest); i {
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
		file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFileResponse); i {
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
			RawDescriptor: file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_goTypes,
		DependencyIndexes: file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_depIdxs,
		EnumInfos:         file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_enumTypes,
		MessageInfos:      file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_msgTypes,
	}.Build()
	File_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto = out.File
	file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_rawDesc = nil
	file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_goTypes = nil
	file_github_com_googleinterns_step224_2020_alicja_config_proto_interface_proto_depIdxs = nil
}