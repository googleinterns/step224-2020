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
// This proto defines the service-level config for Hermes.
// This is also the external service API for Hermes.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.11.4
// source: github.com/googleinterns/step224-2020/config/proto/service.proto

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

// HermesProbeRequest is used for starting monitoring a new storage system using a Hermes probe.
type HermesProbeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProbeConfig *HermesProbeDef `protobuf:"bytes,1,opt,name=probe_config,json=probeConfig,proto3" json:"probe_config,omitempty"`
}

func (x *HermesProbeRequest) Reset() {
	*x = HermesProbeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HermesProbeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HermesProbeRequest) ProtoMessage() {}

func (x *HermesProbeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HermesProbeRequest.ProtoReflect.Descriptor instead.
func (*HermesProbeRequest) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescGZIP(), []int{0}
}

func (x *HermesProbeRequest) GetProbeConfig() *HermesProbeDef {
	if x != nil {
		return x.ProbeConfig
	}
	return nil
}

// TODO(#29) Add exit status to probe responses.
type HermesProbeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HermesProbeResponse) Reset() {
	*x = HermesProbeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HermesProbeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HermesProbeResponse) ProtoMessage() {}

func (x *HermesProbeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HermesProbeResponse.ProtoReflect.Descriptor instead.
func (*HermesProbeResponse) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescGZIP(), []int{1}
}

// StopMonitoringSystemRequest specifies a target so that Hermes knows which target to stop monitoring.
type StopMonitoringSystemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Targets *Targets `protobuf:"bytes,1,opt,name=targets,proto3" json:"targets,omitempty"`
}

func (x *StopMonitoringSystemRequest) Reset() {
	*x = StopMonitoringSystemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopMonitoringSystemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopMonitoringSystemRequest) ProtoMessage() {}

func (x *StopMonitoringSystemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopMonitoringSystemRequest.ProtoReflect.Descriptor instead.
func (*StopMonitoringSystemRequest) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescGZIP(), []int{2}
}

func (x *StopMonitoringSystemRequest) GetTargets() *Targets {
	if x != nil {
		return x.Targets
	}
	return nil
}

// TODO(#29) Add exit status for probe responses.
type StopMonitoringSystemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StopMonitoringSystemResponse) Reset() {
	*x = StopMonitoringSystemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopMonitoringSystemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopMonitoringSystemResponse) ProtoMessage() {}

func (x *StopMonitoringSystemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopMonitoringSystemResponse.ProtoReflect.Descriptor instead.
func (*StopMonitoringSystemResponse) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescGZIP(), []int{3}
}

type ListMonitoredSystemsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListMonitoredSystemsRequest) Reset() {
	*x = ListMonitoredSystemsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMonitoredSystemsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMonitoredSystemsRequest) ProtoMessage() {}

func (x *ListMonitoredSystemsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMonitoredSystemsRequest.ProtoReflect.Descriptor instead.
func (*ListMonitoredSystemsRequest) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescGZIP(), []int{4}
}

// ListMonitoredSystemsResponse message holds the list of targets returned by ListMonitoredStorageSystems()
type ListMonitoredSystemsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// TODO(#29) Add exit status to probe response.
	Targets *Targets `protobuf:"bytes,1,opt,name=targets,proto3" json:"targets,omitempty"`
}

func (x *ListMonitoredSystemsResponse) Reset() {
	*x = ListMonitoredSystemsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMonitoredSystemsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMonitoredSystemsResponse) ProtoMessage() {}

func (x *ListMonitoredSystemsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMonitoredSystemsResponse.ProtoReflect.Descriptor instead.
func (*ListMonitoredSystemsResponse) Descriptor() ([]byte, []int) {
	return file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescGZIP(), []int{5}
}

func (x *ListMonitoredSystemsResponse) GetTargets() *Targets {
	if x != nil {
		return x.Targets
	}
	return nil
}

var File_github_com_googleinterns_step224_2020_config_proto_service_proto protoreflect.FileDescriptor

var file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDesc = []byte{
	0x0a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x32,
	0x32, 0x34, 0x2d, 0x32, 0x30, 0x32, 0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x1a, 0x3e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x32, 0x32, 0x34, 0x2d, 0x32, 0x30, 0x32,
	0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70,
	0x72, 0x6f, 0x62, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x40, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x32, 0x32, 0x34, 0x2d, 0x32, 0x30, 0x32,
	0x30, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4f, 0x0a, 0x12,
	0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x39, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65,
	0x73, 0x2e, 0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x44, 0x65, 0x66,
	0x52, 0x0b, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x15, 0x0a,
	0x13, 0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x48, 0x0a, 0x1b, 0x53, 0x74, 0x6f, 0x70, 0x4d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x73, 0x52, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x22, 0x1e,
	0x0a, 0x1c, 0x53, 0x74, 0x6f, 0x70, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d,
	0x0a, 0x1b, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x53,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x49, 0x0a,
	0x1c, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x53, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a,
	0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x52,
	0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x32, 0xbb, 0x02, 0x0a, 0x06, 0x48, 0x65, 0x72,
	0x6d, 0x65, 0x73, 0x12, 0x59, 0x0a, 0x1c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x53, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x12, 0x1a, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x48, 0x65, 0x72,
	0x6d, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x48, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x50,
	0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6a,
	0x0a, 0x1b, 0x53, 0x74, 0x6f, 0x70, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x23, 0x2e,
	0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x4d, 0x6f, 0x6e, 0x69, 0x74,
	0x6f, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x24, 0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x70,
	0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6a, 0x0a, 0x1b, 0x4c, 0x69,
	0x73, 0x74, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x23, 0x2e, 0x68, 0x65, 0x72, 0x6d,
	0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x65, 0x64,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24,
	0x2e, 0x68, 0x65, 0x72, 0x6d, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x65, 0x64, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x73, 0x2f, 0x73, 0x74, 0x65, 0x70, 0x32, 0x32, 0x34, 0x2d, 0x32, 0x30, 0x32, 0x30, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescOnce sync.Once
	file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescData = file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDesc
)

func file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescGZIP() []byte {
	file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescOnce.Do(func() {
		file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescData)
	})
	return file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDescData
}

var file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_github_com_googleinterns_step224_2020_config_proto_service_proto_goTypes = []interface{}{
	(*HermesProbeRequest)(nil),           // 0: hermes.HermesProbeRequest
	(*HermesProbeResponse)(nil),          // 1: hermes.HermesProbeResponse
	(*StopMonitoringSystemRequest)(nil),  // 2: hermes.StopMonitoringSystemRequest
	(*StopMonitoringSystemResponse)(nil), // 3: hermes.StopMonitoringSystemResponse
	(*ListMonitoredSystemsRequest)(nil),  // 4: hermes.ListMonitoredSystemsRequest
	(*ListMonitoredSystemsResponse)(nil), // 5: hermes.ListMonitoredSystemsResponse
	(*HermesProbeDef)(nil),               // 6: hermes.HermesProbeDef
	(*Targets)(nil),                      // 7: hermes.Targets
}
var file_github_com_googleinterns_step224_2020_config_proto_service_proto_depIdxs = []int32{
	6, // 0: hermes.HermesProbeRequest.probe_config:type_name -> hermes.HermesProbeDef
	7, // 1: hermes.StopMonitoringSystemRequest.targets:type_name -> hermes.Targets
	7, // 2: hermes.ListMonitoredSystemsResponse.targets:type_name -> hermes.Targets
	0, // 3: hermes.Hermes.StartMonitoringStorageSystem:input_type -> hermes.HermesProbeRequest
	2, // 4: hermes.Hermes.StopMonitoringStorageSystem:input_type -> hermes.StopMonitoringSystemRequest
	4, // 5: hermes.Hermes.ListMonitoredStorageSystems:input_type -> hermes.ListMonitoredSystemsRequest
	1, // 6: hermes.Hermes.StartMonitoringStorageSystem:output_type -> hermes.HermesProbeResponse
	3, // 7: hermes.Hermes.StopMonitoringStorageSystem:output_type -> hermes.StopMonitoringSystemResponse
	5, // 8: hermes.Hermes.ListMonitoredStorageSystems:output_type -> hermes.ListMonitoredSystemsResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_github_com_googleinterns_step224_2020_config_proto_service_proto_init() }
func file_github_com_googleinterns_step224_2020_config_proto_service_proto_init() {
	if File_github_com_googleinterns_step224_2020_config_proto_service_proto != nil {
		return
	}
	file_github_com_googleinterns_step224_2020_config_proto_probe_proto_init()
	file_github_com_googleinterns_step224_2020_config_proto_targets_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HermesProbeRequest); i {
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
		file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HermesProbeResponse); i {
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
		file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopMonitoringSystemRequest); i {
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
		file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopMonitoringSystemResponse); i {
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
		file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMonitoredSystemsRequest); i {
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
		file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMonitoredSystemsResponse); i {
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
			RawDescriptor: file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_googleinterns_step224_2020_config_proto_service_proto_goTypes,
		DependencyIndexes: file_github_com_googleinterns_step224_2020_config_proto_service_proto_depIdxs,
		MessageInfos:      file_github_com_googleinterns_step224_2020_config_proto_service_proto_msgTypes,
	}.Build()
	File_github_com_googleinterns_step224_2020_config_proto_service_proto = out.File
	file_github_com_googleinterns_step224_2020_config_proto_service_proto_rawDesc = nil
	file_github_com_googleinterns_step224_2020_config_proto_service_proto_goTypes = nil
	file_github_com_googleinterns_step224_2020_config_proto_service_proto_depIdxs = nil
}
