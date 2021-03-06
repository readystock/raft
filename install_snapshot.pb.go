// Code generated by protoc-gen-go. DO NOT EDIT.
// source: install_snapshot.proto

package raft

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SnapshotVersion int32

const (
	SnapshotVersion_SnapshotVersionMin SnapshotVersion = 0
	SnapshotVersion_SnapshotVersionMax SnapshotVersion = 1
)

var SnapshotVersion_name = map[int32]string{
	0: "SnapshotVersionMin",
	1: "SnapshotVersionMax",
}

var SnapshotVersion_value = map[string]int32{
	"SnapshotVersionMin": 0,
	"SnapshotVersionMax": 1,
}

func (x SnapshotVersion) String() string {
	return proto.EnumName(SnapshotVersion_name, int32(x))
}

func (SnapshotVersion) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3878e8ee4342d62f, []int{0}
}

type InstallSnapshotRequestWrapper struct {
	Request *InstallSnapshotRequest `protobuf:"bytes,1,opt,name=Request,proto3" json:"Request,omitempty"`
	// With gRPC we want to send the snapshot data along with the message.
	Snapshot             []byte   `protobuf:"bytes,2,opt,name=Snapshot,proto3" json:"Snapshot,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallSnapshotRequestWrapper) Reset()         { *m = InstallSnapshotRequestWrapper{} }
func (m *InstallSnapshotRequestWrapper) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotRequestWrapper) ProtoMessage()    {}
func (*InstallSnapshotRequestWrapper) Descriptor() ([]byte, []int) {
	return fileDescriptor_3878e8ee4342d62f, []int{0}
}

func (m *InstallSnapshotRequestWrapper) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallSnapshotRequestWrapper.Unmarshal(m, b)
}
func (m *InstallSnapshotRequestWrapper) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallSnapshotRequestWrapper.Marshal(b, m, deterministic)
}
func (m *InstallSnapshotRequestWrapper) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallSnapshotRequestWrapper.Merge(m, src)
}
func (m *InstallSnapshotRequestWrapper) XXX_Size() int {
	return xxx_messageInfo_InstallSnapshotRequestWrapper.Size(m)
}
func (m *InstallSnapshotRequestWrapper) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallSnapshotRequestWrapper.DiscardUnknown(m)
}

var xxx_messageInfo_InstallSnapshotRequestWrapper proto.InternalMessageInfo

func (m *InstallSnapshotRequestWrapper) GetRequest() *InstallSnapshotRequest {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *InstallSnapshotRequestWrapper) GetSnapshot() []byte {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

type InstallSnapshotRequest struct {
	// Required field on all requests.
	Header          *RPCHeader      `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"`
	SnapshotVersion SnapshotVersion `protobuf:"varint,2,opt,name=SnapshotVersion,proto3,enum=raft.SnapshotVersion" json:"SnapshotVersion,omitempty"`
	Term            uint64          `protobuf:"varint,3,opt,name=Term,proto3" json:"Term,omitempty"`
	Leader          []byte          `protobuf:"bytes,4,opt,name=Leader,proto3" json:"Leader,omitempty"`
	// These are the last index/term included in the snapshot
	LastLogIndex uint64 `protobuf:"varint,5,opt,name=LastLogIndex,proto3" json:"LastLogIndex,omitempty"`
	LastLogTerm  uint64 `protobuf:"varint,6,opt,name=LastLogTerm,proto3" json:"LastLogTerm,omitempty"`
	// Peer Set in the snapshot. This is deprecated in favor of Configuration
	// but remains here in case we receive an InstallSnapshot from a leader
	// that's running old code.
	Peers []byte `protobuf:"bytes,7,opt,name=Peers,proto3" json:"Peers,omitempty"`
	// Cluster membership.
	Configuration []byte `protobuf:"bytes,8,opt,name=Configuration,proto3" json:"Configuration,omitempty"`
	// Log index where 'Configuration' entry was originally written.
	ConfigurationIndex uint64 `protobuf:"varint,9,opt,name=ConfigurationIndex,proto3" json:"ConfigurationIndex,omitempty"`
	// Size of the snapshot
	Size                 int64    `protobuf:"varint,10,opt,name=Size,proto3" json:"Size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InstallSnapshotRequest) Reset()         { *m = InstallSnapshotRequest{} }
func (m *InstallSnapshotRequest) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotRequest) ProtoMessage()    {}
func (*InstallSnapshotRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3878e8ee4342d62f, []int{1}
}

func (m *InstallSnapshotRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallSnapshotRequest.Unmarshal(m, b)
}
func (m *InstallSnapshotRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallSnapshotRequest.Marshal(b, m, deterministic)
}
func (m *InstallSnapshotRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallSnapshotRequest.Merge(m, src)
}
func (m *InstallSnapshotRequest) XXX_Size() int {
	return xxx_messageInfo_InstallSnapshotRequest.Size(m)
}
func (m *InstallSnapshotRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallSnapshotRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InstallSnapshotRequest proto.InternalMessageInfo

func (m *InstallSnapshotRequest) GetHeader() *RPCHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *InstallSnapshotRequest) GetSnapshotVersion() SnapshotVersion {
	if m != nil {
		return m.SnapshotVersion
	}
	return SnapshotVersion_SnapshotVersionMin
}

func (m *InstallSnapshotRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *InstallSnapshotRequest) GetLeader() []byte {
	if m != nil {
		return m.Leader
	}
	return nil
}

func (m *InstallSnapshotRequest) GetLastLogIndex() uint64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

func (m *InstallSnapshotRequest) GetLastLogTerm() uint64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

func (m *InstallSnapshotRequest) GetPeers() []byte {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *InstallSnapshotRequest) GetConfiguration() []byte {
	if m != nil {
		return m.Configuration
	}
	return nil
}

func (m *InstallSnapshotRequest) GetConfigurationIndex() uint64 {
	if m != nil {
		return m.ConfigurationIndex
	}
	return 0
}

func (m *InstallSnapshotRequest) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type InstallSnapshotResponse struct {
	// Required field on all requests.
	Header               *RPCHeader `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"`
	Term                 uint64     `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	Success              bool       `protobuf:"varint,3,opt,name=Success,proto3" json:"Success,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *InstallSnapshotResponse) Reset()         { *m = InstallSnapshotResponse{} }
func (m *InstallSnapshotResponse) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotResponse) ProtoMessage()    {}
func (*InstallSnapshotResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3878e8ee4342d62f, []int{2}
}

func (m *InstallSnapshotResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstallSnapshotResponse.Unmarshal(m, b)
}
func (m *InstallSnapshotResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstallSnapshotResponse.Marshal(b, m, deterministic)
}
func (m *InstallSnapshotResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstallSnapshotResponse.Merge(m, src)
}
func (m *InstallSnapshotResponse) XXX_Size() int {
	return xxx_messageInfo_InstallSnapshotResponse.Size(m)
}
func (m *InstallSnapshotResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InstallSnapshotResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InstallSnapshotResponse proto.InternalMessageInfo

func (m *InstallSnapshotResponse) GetHeader() *RPCHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *InstallSnapshotResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *InstallSnapshotResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterEnum("raft.SnapshotVersion", SnapshotVersion_name, SnapshotVersion_value)
	proto.RegisterType((*InstallSnapshotRequestWrapper)(nil), "raft.InstallSnapshotRequestWrapper")
	proto.RegisterType((*InstallSnapshotRequest)(nil), "raft.InstallSnapshotRequest")
	proto.RegisterType((*InstallSnapshotResponse)(nil), "raft.InstallSnapshotResponse")
}

func init() { proto.RegisterFile("install_snapshot.proto", fileDescriptor_3878e8ee4342d62f) }

var fileDescriptor_3878e8ee4342d62f = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x5f, 0x4b, 0xfb, 0x30,
	0x14, 0xfd, 0x75, 0xeb, 0xba, 0xfd, 0xee, 0xa6, 0x1b, 0x17, 0xad, 0x65, 0x28, 0x94, 0x22, 0x58,
	0x7c, 0xe8, 0xc3, 0x04, 0x5f, 0x45, 0xf6, 0xe2, 0x60, 0xc2, 0xc8, 0x44, 0x1f, 0x47, 0xdc, 0xb2,
	0xad, 0x30, 0x9b, 0x98, 0x64, 0x30, 0xfc, 0xdc, 0x7e, 0x00, 0x59, 0xd2, 0xca, 0xfe, 0xf4, 0xc1,
	0xb7, 0x9c, 0x73, 0x72, 0xee, 0xb9, 0x39, 0x2d, 0xf8, 0x69, 0xa6, 0x34, 0x5d, 0xad, 0x26, 0x2a,
	0xa3, 0x42, 0x2d, 0xb9, 0x4e, 0x84, 0xe4, 0x9a, 0xa3, 0x2b, 0xe9, 0x5c, 0x77, 0x3b, 0x52, 0x4c,
	0x27, 0x4b, 0x46, 0x67, 0x4c, 0x5a, 0x3e, 0x52, 0x70, 0x35, 0xb0, 0x8e, 0x71, 0x6e, 0x20, 0xec,
	0x73, 0xcd, 0x94, 0x7e, 0x93, 0x54, 0x08, 0x26, 0xf1, 0x1e, 0xea, 0x39, 0x13, 0x38, 0xa1, 0x13,
	0x37, 0x7b, 0x97, 0xc9, 0x76, 0x54, 0x52, 0xee, 0x22, 0xc5, 0x65, 0xec, 0x42, 0xa3, 0xd0, 0x82,
	0x4a, 0xe8, 0xc4, 0x2d, 0xf2, 0x8b, 0xa3, 0xef, 0x0a, 0xf8, 0xe5, 0x7e, 0xbc, 0x01, 0xef, 0xc9,
	0xec, 0x97, 0xa7, 0xb5, 0x6d, 0x1a, 0x19, 0xf5, 0x2d, 0x4d, 0x72, 0x19, 0x1f, 0xa0, 0x5d, 0x78,
	0x5f, 0x99, 0x54, 0x29, 0xcf, 0x4c, 0xcc, 0x69, 0xef, 0xdc, 0x3a, 0x0e, 0x44, 0x72, 0x78, 0x1b,
	0x11, 0xdc, 0x17, 0x26, 0x3f, 0x82, 0x6a, 0xe8, 0xc4, 0x2e, 0x31, 0x67, 0xf4, 0xc1, 0x1b, 0xda,
	0x74, 0xd7, 0xac, 0x9c, 0x23, 0x8c, 0xa0, 0x35, 0xa4, 0x4a, 0x0f, 0xf9, 0x62, 0x90, 0xcd, 0xd8,
	0x26, 0xa8, 0x19, 0xcf, 0x1e, 0x87, 0x21, 0x34, 0x73, 0x6c, 0xc6, 0x7a, 0xe6, 0xca, 0x2e, 0x85,
	0x67, 0x50, 0x1b, 0x31, 0x26, 0x55, 0x50, 0x37, 0xc3, 0x2d, 0xc0, 0x6b, 0x38, 0xe9, 0xf3, 0x6c,
	0x9e, 0x2e, 0xd6, 0x92, 0xea, 0xed, 0x33, 0x1a, 0x46, 0xdd, 0x27, 0x31, 0x01, 0xdc, 0x23, 0xec,
	0x1e, 0xff, 0x4d, 0x48, 0x89, 0xb2, 0x7d, 0xdd, 0x38, 0xfd, 0x62, 0x01, 0x84, 0x4e, 0x5c, 0x25,
	0xe6, 0x1c, 0x09, 0xb8, 0x38, 0x6a, 0x5d, 0x09, 0x9e, 0x29, 0xf6, 0xf7, 0xda, 0x8b, 0xd6, 0x2a,
	0x3b, 0xad, 0x05, 0x50, 0x1f, 0xaf, 0xa7, 0x53, 0xa6, 0x94, 0x29, 0xb3, 0x41, 0x0a, 0x78, 0xfb,
	0x78, 0xf4, 0x91, 0xd0, 0x07, 0x3c, 0xa0, 0x9e, 0xd3, 0xac, 0xf3, 0xaf, 0x8c, 0xa7, 0x9b, 0x8e,
	0xf3, 0xee, 0x99, 0xff, 0xf4, 0xee, 0x27, 0x00, 0x00, 0xff, 0xff, 0x72, 0x6f, 0xa8, 0x67, 0xd9,
	0x02, 0x00, 0x00,
}
