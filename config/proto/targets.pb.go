// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Evan Spendlove (@evanSpendlove)
//
//  The Targets proto defines the config for a Target.
//  Targets are used by Hermes and Cloudprober to identify storage systems to monitor.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.11.4
// source: github.com/googleinterns/step224-2020/config/proto/targets.proto

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

// TargetSystem is expected to be a storage system.
// It will stay constant until support for new storage systems is added.
type Target_TargetSystem int32

const (
	Target_GOOGLE_CLOUD_STORAGE Target_TargetSystem = 0
	Target_CEPH                 Target_TargetSystem = 1
)

// Enum value maps for Target_TargetSystem.
var (
	Target_TargetSystem_name = map[int32]string{
		0: "GOOGLE_CLOUD_STORAGE",
		1: "CEPH",
	}
	Target_TargetSystem_value = map[string]int32{
		"GOOGLE_CLOUD_STORAGE": 0,
		"CEPH":                 1,
	}
)

func (x Target_TargetSystem) Enum() *Target_TargetSystem {
	p := new(Target_TargetSystem)
	*p = x
	return p
}

func (x Target_TargetSystem) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Target_TargetSystem) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_googleinterns_step224_2020_config_proto_targets_proto_enumTypes[0].Descriptor()
}

func (Target_TargetSystem) Type() protoreflect.EnumType {
	return &file_github_com_googleinterns_step224_2020_config_proto_targets_proto_enumTypes[0]
}

func (x Target_TargetSystem) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Target_TargetSystem.Descriptor instead.
func (Target_TargetSystem) EnumDescriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescGZIP(), []int{0, 0}
}

// TODO(#30) Establish connection method for GCS and Ceph using Go libraries.
type Target_ConnectionType int32

const (
	Target_HTTP  Target_ConnectionType = 0
	Target_HTTPS Target_ConnectionType = 1
)

// Enum value maps for Target_ConnectionType.
var (
	Target_ConnectionType_name = map[int32]string{
		0: "HTTP",
		1: "HTTPS",
	}
	Target_ConnectionType_value = map[string]int32{
		"HTTP":  0,
		"HTTPS": 1,
	}
)

func (x Target_ConnectionType) Enum() *Target_ConnectionType {
	p := new(Target_ConnectionType)
	*p = x
	return p
}

func (x Target_ConnectionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Target_ConnectionType) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_googleinterns_step224_2020_config_proto_targets_proto_enumTypes[1].Descriptor()
}

func (Target_ConnectionType) Type() protoreflect.EnumType {
	return &file_github_com_googleinterns_step224_2020_config_proto_targets_proto_enumTypes[1]
}

func (x Target_ConnectionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Target_ConnectionType.Descriptor instead.
func (Target_ConnectionType) EnumDescriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescGZIP(), []int{0, 1}
}

// TargetDefinition contains all of the metadata necessary for Hermes to establish a connection to a storage system.
// Every probe request will require one or more targets.
type Target struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetSystem   Target_TargetSystem   `protobuf:"varint,1,opt,name=target_system,json=targetSystem,proto3,enum=hermes.Target_TargetSystem" json:"target_system,omitempty"`
	ConnectionType Target_ConnectionType `protobuf:"varint,2,opt,name=connection_type,json=connectionType,proto3,enum=hermes.Target_ConnectionType" json:"connection_type,omitempty"`
	// REQUIRED for Ceph S3 API
	// Port for connecting to Ceph S3 API
	Port int32 `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	// REQUIRED for Ceph S3 API
	// GCS uses service account credentials instead.
	ApiKey                 string `protobuf:"bytes,4,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	TotalSpaceAllocatedMib int64  `protobuf:"varint,5,opt,name=total_space_allocated_mib,json=totalSpaceAllocatedMib,proto3" json:"total_space_allocated_mib,omitempty"`
	TargetUrl              string `protobuf:"bytes,6,opt,name=target_url,json=targetUrl,proto3" json:"target_url,omitempty"` // URL for connecting to the the API of the target storage system.
}

func (x *Target) Reset() {
	*x = Target{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_config_proto_targets_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Target) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Target) ProtoMessage() {}

func (x *Target) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_config_proto_targets_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Target.ProtoReflect.Descriptor instead.
func (*Target) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescGZIP(), []int{0}
}

func (x *Target) GetTargetSystem() Target_TargetSystem {
	if x != nil {
		return x.TargetSystem
	}
	return Target_GOOGLE_CLOUD_STORAGE
}

func (x *Target) GetConnectionType() Target_ConnectionType {
	if x != nil {
		return x.ConnectionType
	}
	return Target_HTTP
}

func (x *Target) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Target) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *Target) GetTotalSpaceAllocatedMib() int64 {
	if x != nil {
		return x.TotalSpaceAllocatedMib
	}
	return 0
}

func (x *Target) GetTargetUrl() string {
	if x != nil {
		return x.TargetUrl
	}
	return ""
}

var File_github_com_googleinterns_step224_2020_config_proto_targets_proto protoreflect.FileDescriptor

var file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDesc = []byte{
	0x0a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x32,
	0x32, 0x34, 0x2d, 0x32, 0x30, 0x32, 0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x22, 0xf4, 0x02, 0x0a, 0x06, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x40, 0x0a, 0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x68,
	0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2e, 0x54, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x46, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1d, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x12, 0x39, 0x0a, 0x19,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x61, 0x6c, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x6d, 0x69, 0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x16, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x70, 0x61, 0x63, 0x65, 0x41, 0x6c, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x65, 0x64, 0x4d, 0x69, 0x62, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x55, 0x72, 0x6c, 0x22, 0x32, 0x0a, 0x0c, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x18, 0x0a, 0x14, 0x47, 0x4f, 0x4f, 0x47, 0x4c, 0x45,
	0x5f, 0x43, 0x4c, 0x4f, 0x55, 0x44, 0x5f, 0x53, 0x54, 0x4f, 0x52, 0x41, 0x47, 0x45, 0x10, 0x00,
	0x12, 0x08, 0x0a, 0x04, 0x43, 0x45, 0x50, 0x48, 0x10, 0x01, 0x22, 0x25, 0x0a, 0x0e, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04,
	0x48, 0x54, 0x54, 0x50, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x48, 0x54, 0x54, 0x50, 0x53, 0x10,
	0x01, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74,
	0x65, 0x70, 0x32, 0x32, 0x34, 0x2d, 0x32, 0x30, 0x32, 0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescOnce sync.Once
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescData = file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDesc
)

func file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescGZIP() []byte {
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescOnce.Do(func() {
		file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescData)
	})
	return file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDescData
}

var file_github_com_googleinterns_step224_2020_config_proto_targets_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_github_com_googleinterns_step224_2020_config_proto_targets_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_googleinterns_step224_2020_config_proto_targets_proto_goTypes = []interface{}{
	(Target_TargetSystem)(0),   // 0: hermes.Target.TargetSystem
	(Target_ConnectionType)(0), // 1: hermes.Target.ConnectionType
	(*Target)(nil),             // 2: hermes.Target
}
var file_github_com_googleinterns_step224_2020_config_proto_targets_proto_depIdxs = []int32{
	0, // 0: hermes.Target.target_system:type_name -> hermes.Target.TargetSystem
	1, // 1: hermes.Target.connection_type:type_name -> hermes.Target.ConnectionType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_github_com_googleinterns_step224_2020_config_proto_targets_proto_init() }
func file_github_com_googleinterns_step224_2020_config_proto_targets_proto_init() {
	if File_github_com_googleinterns_step224_2020_config_proto_targets_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_googleinterns_step224_2020_config_proto_targets_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Target); i {
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
			RawDescriptor: file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_googleinterns_step224_2020_config_proto_targets_proto_goTypes,
		DependencyIndexes: file_github_com_googleinterns_step224_2020_config_proto_targets_proto_depIdxs,
		EnumInfos:         file_github_com_googleinterns_step224_2020_config_proto_targets_proto_enumTypes,
		MessageInfos:      file_github_com_googleinterns_step224_2020_config_proto_targets_proto_msgTypes,
	}.Build()
	File_github_com_googleinterns_step224_2020_config_proto_targets_proto = out.File
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_rawDesc = nil
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_goTypes = nil
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_depIdxs = nil
}