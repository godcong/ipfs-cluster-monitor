// Code generated by protoc-gen-go. DO NOT EDIT.
// source: monitor.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

// MonitorType ...
type MonitorType int32

// MonitorType_Init ...
const (
	MonitorType_Init  MonitorType = 0
	MonitorType_Reset MonitorType = 1
	MonitorType_Info  MonitorType = 2
)

// MonitorType_name ...
var MonitorType_name = map[int32]string{
	0: "Init",
	1: "Reset",
	2: "Info",
}

// MonitorType_value ...
var MonitorType_value = map[string]int32{
	"Init":  0,
	"Reset": 1,
	"Info":  2,
}

// String ...
func (x MonitorType) String() string {
	return proto.EnumName(MonitorType_name, int32(x))
}

// EnumDescriptor ...
func (MonitorType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{0}
}

// MonitorRequest ...
type MonitorRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *MonitorRequest) Reset() { *m = MonitorRequest{} }

// String ...
func (m *MonitorRequest) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*MonitorRequest) ProtoMessage() {}

// Descriptor ...
func (*MonitorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{0}
}

// XXX_Unmarshal ...
func (m *MonitorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorRequest.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *MonitorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorRequest.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *MonitorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorRequest.Merge(m, src)
}

// XXX_Size ...
func (m *MonitorRequest) XXX_Size() int {
	return xxx_messageInfo_MonitorRequest.Size(m)
}

// XXX_DiscardUnknown ...
func (m *MonitorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorRequest proto.InternalMessageInfo

// MonitorInitRequest ...
type MonitorInitRequest struct {
	Bootstrap            string   `protobuf:"bytes,1,opt,name=bootstrap,proto3" json:"bootstrap,omitempty"`
	Secret               string   `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	Workspace            string   `protobuf:"bytes,3,opt,name=workspace,proto3" json:"workspace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *MonitorInitRequest) Reset() { *m = MonitorInitRequest{} }

// String ...
func (m *MonitorInitRequest) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*MonitorInitRequest) ProtoMessage() {}

// Descriptor ...
func (*MonitorInitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{1}
}

// XXX_Unmarshal ...
func (m *MonitorInitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorInitRequest.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *MonitorInitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorInitRequest.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *MonitorInitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorInitRequest.Merge(m, src)
}

// XXX_Size ...
func (m *MonitorInitRequest) XXX_Size() int {
	return xxx_messageInfo_MonitorInitRequest.Size(m)
}

// XXX_DiscardUnknown ...
func (m *MonitorInitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorInitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorInitRequest proto.InternalMessageInfo

// GetBootstrap ...
func (m *MonitorInitRequest) GetBootstrap() string {
	if m != nil {
		return m.Bootstrap
	}
	return ""
}

// GetSecret ...
func (m *MonitorInitRequest) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

// GetWorkspace ...
func (m *MonitorInitRequest) GetWorkspace() string {
	if m != nil {
		return m.Workspace
	}
	return ""
}

// MonitorBootstrapReply ...
type MonitorBootstrapReply struct {
	Bootstraps           []string `protobuf:"bytes,1,rep,name=bootstraps,proto3" json:"bootstraps,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *MonitorBootstrapReply) Reset() { *m = MonitorBootstrapReply{} }

// String ...
func (m *MonitorBootstrapReply) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*MonitorBootstrapReply) ProtoMessage() {}

// Descriptor ...
func (*MonitorBootstrapReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{2}
}

// XXX_Unmarshal ...
func (m *MonitorBootstrapReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorBootstrapReply.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *MonitorBootstrapReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorBootstrapReply.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *MonitorBootstrapReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorBootstrapReply.Merge(m, src)
}

// XXX_Size ...
func (m *MonitorBootstrapReply) XXX_Size() int {
	return xxx_messageInfo_MonitorBootstrapReply.Size(m)
}

// XXX_DiscardUnknown ...
func (m *MonitorBootstrapReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorBootstrapReply.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorBootstrapReply proto.InternalMessageInfo

// GetBootstraps ...
func (m *MonitorBootstrapReply) GetBootstraps() []string {
	if m != nil {
		return m.Bootstraps
	}
	return nil
}

// MonitorPinReply ...
type MonitorPinReply struct {
	Pins                 []string `protobuf:"bytes,1,rep,name=pins,proto3" json:"pins,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *MonitorPinReply) Reset() { *m = MonitorPinReply{} }

// String ...
func (m *MonitorPinReply) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*MonitorPinReply) ProtoMessage() {}

// Descriptor ...
func (*MonitorPinReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{3}
}

// XXX_Unmarshal ...
func (m *MonitorPinReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorPinReply.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *MonitorPinReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorPinReply.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *MonitorPinReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorPinReply.Merge(m, src)
}

// XXX_Size ...
func (m *MonitorPinReply) XXX_Size() int {
	return xxx_messageInfo_MonitorPinReply.Size(m)
}

// XXX_DiscardUnknown ...
func (m *MonitorPinReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorPinReply.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorPinReply proto.InternalMessageInfo

// GetPins ...
func (m *MonitorPinReply) GetPins() []string {
	if m != nil {
		return m.Pins
	}
	return nil
}

// MonitorProcRequest ...
type MonitorProcRequest struct {
	Type                 MonitorType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.MonitorType" json:"type,omitempty"`
	BootStrap            string      `protobuf:"bytes,2,opt,name=boot_strap,json=bootStrap,proto3" json:"boot_strap,omitempty"`
	Secret               string      `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
	Workspace            string      `protobuf:"bytes,4,opt,name=workspace,proto3" json:"workspace,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

// Reset ...
func (m *MonitorProcRequest) Reset() { *m = MonitorProcRequest{} }

// String ...
func (m *MonitorProcRequest) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*MonitorProcRequest) ProtoMessage() {}

// Descriptor ...
func (*MonitorProcRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{4}
}

// XXX_Unmarshal ...
func (m *MonitorProcRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorProcRequest.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *MonitorProcRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorProcRequest.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *MonitorProcRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorProcRequest.Merge(m, src)
}

// XXX_Size ...
func (m *MonitorProcRequest) XXX_Size() int {
	return xxx_messageInfo_MonitorProcRequest.Size(m)
}

// XXX_DiscardUnknown ...
func (m *MonitorProcRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorProcRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorProcRequest proto.InternalMessageInfo

// GetType ...
func (m *MonitorProcRequest) GetType() MonitorType {
	if m != nil {
		return m.Type
	}
	return MonitorType_Init
}

// GetBootStrap ...
func (m *MonitorProcRequest) GetBootStrap() string {
	if m != nil {
		return m.BootStrap
	}
	return ""
}

// GetSecret ...
func (m *MonitorProcRequest) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

// GetWorkspace ...
func (m *MonitorProcRequest) GetWorkspace() string {
	if m != nil {
		return m.Workspace
	}
	return ""
}

// MonitorCensorRequest ...
type MonitorCensorRequest struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Detail               string   `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *MonitorCensorRequest) Reset() { *m = MonitorCensorRequest{} }

// String ...
func (m *MonitorCensorRequest) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*MonitorCensorRequest) ProtoMessage() {}

// Descriptor ...
func (*MonitorCensorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{5}
}

// XXX_Unmarshal ...
func (m *MonitorCensorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorCensorRequest.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *MonitorCensorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorCensorRequest.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *MonitorCensorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorCensorRequest.Merge(m, src)
}

// XXX_Size ...
func (m *MonitorCensorRequest) XXX_Size() int {
	return xxx_messageInfo_MonitorCensorRequest.Size(m)
}

// XXX_DiscardUnknown ...
func (m *MonitorCensorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorCensorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorCensorRequest proto.InternalMessageInfo

// GetID ...
func (m *MonitorCensorRequest) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

// GetDetail ...
func (m *MonitorCensorRequest) GetDetail() string {
	if m != nil {
		return m.Detail
	}
	return ""
}

// MonitorReply ...
type MonitorReply struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Detail               string   `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *MonitorReply) Reset() { *m = MonitorReply{} }

// String ...
func (m *MonitorReply) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*MonitorReply) ProtoMessage() {}

// Descriptor ...
func (*MonitorReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{6}
}

// XXX_Unmarshal ...
func (m *MonitorReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorReply.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *MonitorReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorReply.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *MonitorReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorReply.Merge(m, src)
}

// XXX_Size ...
func (m *MonitorReply) XXX_Size() int {
	return xxx_messageInfo_MonitorReply.Size(m)
}

// XXX_DiscardUnknown ...
func (m *MonitorReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorReply.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorReply proto.InternalMessageInfo

// GetCode ...
func (m *MonitorReply) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

// GetMessage ...
func (m *MonitorReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

// GetDetail ...
func (m *MonitorReply) GetDetail() string {
	if m != nil {
		return m.Detail
	}
	return ""
}

func init() {
	proto.RegisterEnum("proto.MonitorType", MonitorType_name, MonitorType_value)
	proto.RegisterType((*MonitorRequest)(nil), "proto.MonitorRequest")
	proto.RegisterType((*MonitorInitRequest)(nil), "proto.MonitorInitRequest")
	proto.RegisterType((*MonitorBootstrapReply)(nil), "proto.MonitorBootstrapReply")
	proto.RegisterType((*MonitorPinReply)(nil), "proto.MonitorPinReply")
	proto.RegisterType((*MonitorProcRequest)(nil), "proto.MonitorProcRequest")
	proto.RegisterType((*MonitorCensorRequest)(nil), "proto.MonitorCensorRequest")
	proto.RegisterType((*MonitorReply)(nil), "proto.MonitorReply")
}

func init() { proto.RegisterFile("monitor.proto", fileDescriptor_44174b7b2a306b71) }

var fileDescriptor_44174b7b2a306b71 = []byte{
	// 429 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0xcd, 0x57, 0x17, 0x32, 0x40, 0x88, 0x66, 0xd9, 0x55, 0x40, 0x0b, 0xac, 0x2c, 0x81, 0x56,
	0x08, 0xe5, 0xb0, 0x1c, 0xe0, 0x02, 0x48, 0xdd, 0x5e, 0x72, 0x40, 0x8a, 0x42, 0xef, 0x28, 0x4d,
	0x4d, 0x89, 0x48, 0xe3, 0x60, 0xbb, 0x45, 0xfd, 0x19, 0xdc, 0xf9, 0xb1, 0xab, 0xd8, 0x4e, 0x9a,
	0x44, 0xea, 0x29, 0xf6, 0x1b, 0xbf, 0x67, 0xbf, 0x37, 0x13, 0x78, 0xb2, 0x65, 0x75, 0x29, 0x19,
	0x8f, 0x1b, 0xce, 0x24, 0xc3, 0x99, 0xfa, 0x90, 0x10, 0x82, 0x6f, 0x1a, 0xcf, 0xe8, 0x9f, 0x1d,
	0x15, 0x92, 0xfc, 0x02, 0x34, 0x48, 0x52, 0x97, 0xd2, 0xa0, 0x78, 0x05, 0xfe, 0x8a, 0x31, 0x29,
	0x24, 0xcf, 0x9b, 0xc8, 0xbe, 0xb6, 0x6f, 0xfc, 0xec, 0x08, 0xe0, 0x25, 0x9c, 0x09, 0x5a, 0x70,
	0x2a, 0x23, 0x47, 0x95, 0xcc, 0xae, 0x65, 0xfd, 0x65, 0xfc, 0xb7, 0x68, 0xf2, 0x82, 0x46, 0xae,
	0x66, 0xf5, 0x00, 0xf9, 0x08, 0x17, 0xe6, 0xa6, 0x79, 0xa7, 0x94, 0xd1, 0xa6, 0x3a, 0xe0, 0x2b,
	0x80, 0x5e, 0x5b, 0x44, 0xf6, 0xb5, 0x7b, 0xe3, 0x67, 0x03, 0x84, 0xbc, 0x81, 0xa7, 0x86, 0x98,
	0x96, 0xb5, 0xa6, 0x20, 0x78, 0x4d, 0x59, 0x77, 0x87, 0xd5, 0x9a, 0xfc, 0xb3, 0x7b, 0x2b, 0x29,
	0x67, 0x45, 0x67, 0xe5, 0x2d, 0x78, 0xf2, 0xd0, 0x50, 0xe5, 0x22, 0xb8, 0x45, 0x9d, 0x47, 0x6c,
	0x0e, 0x2e, 0x0f, 0x0d, 0xcd, 0x54, 0x1d, 0x5f, 0xea, 0x57, 0xfc, 0xd0, 0x9e, 0x9d, 0xa3, 0xe7,
	0xef, 0x13, 0xcf, 0xee, 0x69, 0xcf, 0xde, 0xd4, 0xf3, 0x17, 0x78, 0x66, 0x6e, 0xba, 0xa3, 0xb5,
	0xe8, 0x53, 0xc7, 0x00, 0x9c, 0x64, 0x61, 0x82, 0x75, 0x92, 0x45, 0xab, 0xbe, 0xa6, 0x32, 0x2f,
	0xab, 0x2e, 0x51, 0xbd, 0x23, 0x4b, 0x78, 0xdc, 0xf7, 0xcb, 0xf8, 0x2e, 0xd8, 0x5a, 0x9b, 0x99,
	0x65, 0x6a, 0x8d, 0x11, 0x3c, 0xd8, 0x52, 0x21, 0xf2, 0x0d, 0x35, 0xe4, 0x6e, 0x3b, 0x50, 0x75,
	0x87, 0xaa, 0xef, 0xde, 0xc3, 0xa3, 0x81, 0x7f, 0x7c, 0x08, 0x5e, 0xdb, 0xfb, 0xd0, 0x42, 0x1f,
	0x66, 0x19, 0x15, 0x54, 0x86, 0xb6, 0x06, 0x7f, 0xb2, 0xd0, 0xb9, 0xfd, 0xef, 0x40, 0x70, 0x57,
	0xed, 0x84, 0xa4, 0xdc, 0xb0, 0xf0, 0x6b, 0x2f, 0xd0, 0x12, 0xf1, 0xf9, 0x38, 0xd4, 0xc1, 0x20,
	0xbd, 0x38, 0x1f, 0x97, 0x94, 0x0b, 0x62, 0x0d, 0x04, 0xda, 0x56, 0x4d, 0x05, 0x06, 0xed, 0x3b,
	0x25, 0x90, 0x40, 0x38, 0x1d, 0x26, 0xbc, 0x98, 0x1e, 0xd5, 0x0a, 0x57, 0x63, 0x78, 0x3c, 0x7c,
	0xc4, 0xc2, 0xcf, 0x00, 0xc7, 0xf1, 0x3a, 0x25, 0x72, 0x39, 0x79, 0xa1, 0x19, 0x44, 0x62, 0xcd,
	0x3f, 0xc1, 0xeb, 0x82, 0x6d, 0xe3, 0x7d, 0x95, 0xef, 0x79, 0xbc, 0x61, 0xeb, 0x82, 0xd5, 0x9b,
	0x78, 0xf4, 0xf3, 0xcd, 0xcf, 0xc7, 0xf1, 0xa5, 0x2d, 0x98, 0xda, 0xab, 0x33, 0x55, 0xfd, 0x70,
	0x1f, 0x00, 0x00, 0xff, 0xff, 0x41, 0x5e, 0xca, 0x65, 0xab, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ClusterMonitorClient is the client API for ClusterMonitor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterMonitorClient interface {
	MonitorInit(ctx context.Context, in *MonitorInitRequest, opts ...grpc.CallOption) (*MonitorReply, error)
	MonitorProc(ctx context.Context, in *MonitorProcRequest, opts ...grpc.CallOption) (*MonitorReply, error)
	MonitorBootstrap(ctx context.Context, in *MonitorRequest, opts ...grpc.CallOption) (*MonitorBootstrapReply, error)
	MonitorPin(ctx context.Context, in *MonitorRequest, opts ...grpc.CallOption) (*MonitorPinReply, error)
}

type clusterMonitorClient struct {
	cc *grpc.ClientConn
}

// NewClusterMonitorClient ...
func NewClusterMonitorClient(cc *grpc.ClientConn) ClusterMonitorClient {
	return &clusterMonitorClient{cc}
}

// MonitorInit ...
func (c *clusterMonitorClient) MonitorInit(ctx context.Context, in *MonitorInitRequest, opts ...grpc.CallOption) (*MonitorReply, error) {
	out := new(MonitorReply)
	err := c.cc.Invoke(ctx, "/proto.ClusterMonitor/MonitorInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitorProc ...
func (c *clusterMonitorClient) MonitorProc(ctx context.Context, in *MonitorProcRequest, opts ...grpc.CallOption) (*MonitorReply, error) {
	out := new(MonitorReply)
	err := c.cc.Invoke(ctx, "/proto.ClusterMonitor/MonitorProc", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitorBootstrap ...
func (c *clusterMonitorClient) MonitorBootstrap(ctx context.Context, in *MonitorRequest, opts ...grpc.CallOption) (*MonitorBootstrapReply, error) {
	out := new(MonitorBootstrapReply)
	err := c.cc.Invoke(ctx, "/proto.ClusterMonitor/MonitorBootstrap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitorPin ...
func (c *clusterMonitorClient) MonitorPin(ctx context.Context, in *MonitorRequest, opts ...grpc.CallOption) (*MonitorPinReply, error) {
	out := new(MonitorPinReply)
	err := c.cc.Invoke(ctx, "/proto.ClusterMonitor/MonitorPin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterMonitorServer is the server API for ClusterMonitor service.
type ClusterMonitorServer interface {
	MonitorInit(context.Context, *MonitorInitRequest) (*MonitorReply, error)
	MonitorProc(context.Context, *MonitorProcRequest) (*MonitorReply, error)
	MonitorBootstrap(context.Context, *MonitorRequest) (*MonitorBootstrapReply, error)
	MonitorPin(context.Context, *MonitorRequest) (*MonitorPinReply, error)
}

// RegisterClusterMonitorServer ...
func RegisterClusterMonitorServer(s *grpc.Server, srv ClusterMonitorServer) {
	s.RegisterService(&_ClusterMonitor_serviceDesc, srv)
}

func _ClusterMonitor_MonitorInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonitorInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterMonitorServer).MonitorInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClusterMonitor/MonitorInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterMonitorServer).MonitorInit(ctx, req.(*MonitorInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterMonitor_MonitorProc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonitorProcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterMonitorServer).MonitorProc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClusterMonitor/MonitorProc",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterMonitorServer).MonitorProc(ctx, req.(*MonitorProcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterMonitor_MonitorBootstrap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonitorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterMonitorServer).MonitorBootstrap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClusterMonitor/MonitorBootstrap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterMonitorServer).MonitorBootstrap(ctx, req.(*MonitorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterMonitor_MonitorPin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonitorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterMonitorServer).MonitorPin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClusterMonitor/MonitorPin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterMonitorServer).MonitorPin(ctx, req.(*MonitorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClusterMonitor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ClusterMonitor",
	HandlerType: (*ClusterMonitorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MonitorInit",
			Handler:    _ClusterMonitor_MonitorInit_Handler,
		},
		{
			MethodName: "MonitorProc",
			Handler:    _ClusterMonitor_MonitorProc_Handler,
		},
		{
			MethodName: "MonitorBootstrap",
			Handler:    _ClusterMonitor_MonitorBootstrap_Handler,
		},
		{
			MethodName: "MonitorPin",
			Handler:    _ClusterMonitor_MonitorPin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "monitor.proto",
}
