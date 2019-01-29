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

// MonitorType_init_status ...
const (
	MonitorType_init_status MonitorType = 0
)

// MonitorType_name ...
var MonitorType_name = map[int32]string{
	0: "init_status",
}

// MonitorType_value ...
var MonitorType_value = map[string]int32{
	"init_status": 0,
}

// String ...
func (x MonitorType) String() string {
	return proto.EnumName(MonitorType_name, int32(x))
}

// EnumDescriptor ...
func (MonitorType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_44174b7b2a306b71, []int{0}
}

// MonitorInitRequest ...
type MonitorInitRequest struct {
	Bootstrap            string   `protobuf:"bytes,1,opt,name=bootstrap,proto3" json:"bootstrap,omitempty"`
	Secret               string   `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	Path                 string   `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	ClusterPath          string   `protobuf:"bytes,4,opt,name=cluster_path,json=clusterPath,proto3" json:"cluster_path,omitempty"`
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
	return fileDescriptor_44174b7b2a306b71, []int{0}
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

// GetPath ...
func (m *MonitorInitRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

// GetClusterPath ...
func (m *MonitorInitRequest) GetClusterPath() string {
	if m != nil {
		return m.ClusterPath
	}
	return ""
}

// MonitorProcRequest ...
type MonitorProcRequest struct {
	Type                 MonitorType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.MonitorType" json:"type,omitempty"`
	BootStrap            string      `protobuf:"bytes,2,opt,name=boot_strap,json=bootStrap,proto3" json:"boot_strap,omitempty"`
	Secret               string      `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
	IpfsPath             string      `protobuf:"bytes,4,opt,name=ipfs_path,json=ipfsPath,proto3" json:"ipfs_path,omitempty"`
	IpfsClusterPath      string      `protobuf:"bytes,5,opt,name=ipfs_cluster_path,json=ipfsClusterPath,proto3" json:"ipfs_cluster_path,omitempty"`
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
	return fileDescriptor_44174b7b2a306b71, []int{1}
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
	return MonitorType_init_status
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

// GetIpfsPath ...
func (m *MonitorProcRequest) GetIpfsPath() string {
	if m != nil {
		return m.IpfsPath
	}
	return ""
}

// GetIpfsClusterPath ...
func (m *MonitorProcRequest) GetIpfsClusterPath() string {
	if m != nil {
		return m.IpfsClusterPath
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
	return fileDescriptor_44174b7b2a306b71, []int{2}
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
	return fileDescriptor_44174b7b2a306b71, []int{3}
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
	proto.RegisterType((*MonitorInitRequest)(nil), "proto.MonitorInitRequest")
	proto.RegisterType((*MonitorProcRequest)(nil), "proto.MonitorProcRequest")
	proto.RegisterType((*MonitorCensorRequest)(nil), "proto.MonitorCensorRequest")
	proto.RegisterType((*MonitorReply)(nil), "proto.MonitorReply")
}

func init() { proto.RegisterFile("monitor.proto", fileDescriptor_44174b7b2a306b71) }

var fileDescriptor_44174b7b2a306b71 = []byte{
	// 384 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcb, 0x6e, 0xea, 0x30,
	0x10, 0x86, 0x49, 0xb8, 0x9c, 0xc3, 0xc0, 0x81, 0x53, 0x53, 0x55, 0xe9, 0x55, 0x34, 0x8b, 0xaa,
	0x62, 0x91, 0x05, 0x5d, 0x75, 0xd3, 0x4a, 0xc0, 0x86, 0x45, 0xa5, 0x28, 0x65, 0x8f, 0x42, 0x70,
	0x43, 0x24, 0x88, 0x53, 0x7b, 0xb2, 0x60, 0xd9, 0x57, 0xe8, 0xf3, 0xf4, 0xe1, 0x2a, 0x3b, 0x4e,
	0x71, 0x90, 0xba, 0x8a, 0x67, 0xfe, 0x99, 0xc9, 0xe7, 0xdf, 0x03, 0xff, 0x76, 0x2c, 0x4d, 0x90,
	0x71, 0x2f, 0xe3, 0x0c, 0x19, 0x69, 0xaa, 0x8f, 0xfb, 0x61, 0x01, 0x79, 0x29, 0x84, 0x79, 0x9a,
	0x60, 0x40, 0xdf, 0x73, 0x2a, 0x90, 0x5c, 0x41, 0x7b, 0xc5, 0x18, 0x0a, 0xe4, 0x61, 0xe6, 0x58,
	0x43, 0xeb, 0xbe, 0x1d, 0x1c, 0x12, 0xe4, 0x0c, 0x5a, 0x82, 0x46, 0x9c, 0xa2, 0x63, 0x2b, 0x49,
	0x47, 0x84, 0x40, 0x23, 0x0b, 0x71, 0xe3, 0xd4, 0x55, 0x56, 0x9d, 0xc9, 0x2d, 0x74, 0xa3, 0x6d,
	0x2e, 0x90, 0xf2, 0xa5, 0xd2, 0x1a, 0x4a, 0xeb, 0xe8, 0x9c, 0x1f, 0xe2, 0xc6, 0xfd, 0x3a, 0x30,
	0xf8, 0x9c, 0x45, 0x25, 0xc3, 0x1d, 0x34, 0x70, 0x9f, 0x51, 0xf5, 0xfb, 0xde, 0x98, 0x14, 0xdc,
	0x9e, 0x2e, 0x5c, 0xec, 0x33, 0x1a, 0x28, 0x9d, 0x5c, 0x03, 0x48, 0xb4, 0x65, 0x01, 0x6b, 0x1f,
	0x60, 0x5f, 0x8f, 0x60, 0xeb, 0x15, 0xd8, 0x4b, 0x68, 0x27, 0xd9, 0x9b, 0x30, 0xa9, 0xfe, 0xca,
	0x84, 0x44, 0x22, 0x23, 0x38, 0x51, 0x62, 0x05, 0xbd, 0xa9, 0x8a, 0xfa, 0x52, 0x98, 0x1a, 0xf8,
	0x4f, 0x70, 0xaa, 0xa1, 0xa6, 0x34, 0x15, 0x8c, 0x97, 0xfc, 0x3d, 0xb0, 0xe7, 0x33, 0x6d, 0x9e,
	0x3d, 0x9f, 0x49, 0x90, 0x35, 0xc5, 0x30, 0xd9, 0x96, 0xae, 0x15, 0x91, 0xbb, 0x80, 0xae, 0xee,
	0x0f, 0x68, 0xb6, 0xdd, 0x4b, 0x17, 0x23, 0xb6, 0x2e, 0xee, 0xdd, 0x0c, 0xd4, 0x99, 0x38, 0xf0,
	0x67, 0x47, 0x85, 0x08, 0x63, 0xaa, 0x9b, 0xcb, 0xd0, 0x98, 0x5a, 0x37, 0xa7, 0x8e, 0x6e, 0xa0,
	0x63, 0x58, 0x45, 0xfa, 0xd0, 0x49, 0xd2, 0x44, 0x9a, 0x14, 0x62, 0x2e, 0xfe, 0xd7, 0xc6, 0x9f,
	0x16, 0xf4, 0xf4, 0x2d, 0x74, 0x1d, 0x79, 0xfe, 0x69, 0x91, 0xab, 0x40, 0xce, 0xab, 0x8e, 0x1b,
	0xeb, 0x71, 0x31, 0xa8, 0x4a, 0x8a, 0xdb, 0xad, 0x19, 0x03, 0xe4, 0x3b, 0x1e, 0x0f, 0x30, 0xde,
	0xf6, 0x97, 0x01, 0x93, 0x47, 0x18, 0x46, 0x6c, 0xe7, 0xc5, 0x09, 0x6e, 0xf2, 0x95, 0x17, 0xb3,
	0x75, 0xc4, 0xd2, 0xd8, 0xab, 0x2c, 0xee, 0x64, 0x50, 0xa5, 0xf6, 0x65, 0xd2, 0xb7, 0x56, 0x2d,
	0xa5, 0x3e, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0x70, 0xf3, 0x1b, 0xcc, 0xe7, 0x02, 0x00, 0x00,
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

// ClusterMonitorServer is the server API for ClusterMonitor service.
type ClusterMonitorServer interface {
	MonitorInit(context.Context, *MonitorInitRequest) (*MonitorReply, error)
	MonitorProc(context.Context, *MonitorProcRequest) (*MonitorReply, error)
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "monitor.proto",
}
