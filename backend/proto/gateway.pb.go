// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gateway.proto

package proto

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Category int32

const (
	Category_CMD    Category = 0
	Category_CONFIG Category = 1
)

var Category_name = map[int32]string{
	0: "CMD",
	1: "CONFIG",
}

var Category_value = map[string]int32{
	"CMD":    0,
	"CONFIG": 1,
}

func (x Category) String() string {
	return proto.EnumName(Category_name, int32(x))
}

func (Category) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{0}
}

type Device int32

const (
	Device_LAMP   Device = 0
	Device_HEATER Device = 1
)

var Device_name = map[int32]string{
	0: "LAMP",
	1: "HEATER",
}

var Device_value = map[string]int32{
	"LAMP":   0,
	"HEATER": 1,
}

func (x Device) String() string {
	return proto.EnumName(Device_name, int32(x))
}

func (Device) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{1}
}

type Operation int32

const (
	Operation_OFF Operation = 0
	Operation_ON  Operation = 1
)

var Operation_name = map[int32]string{
	0: "OFF",
	1: "ON",
}

var Operation_value = map[string]int32{
	"OFF": 0,
	"ON":  1,
}

func (x Operation) String() string {
	return proto.EnumName(Operation_name, int32(x))
}

func (Operation) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{2}
}

//message是固定的。UserInfo是类名，可以随意指定，符合规范即可
type UplinkFrame struct {
	FrameType            int32    `protobuf:"varint,1,opt,name=FrameType,proto3" json:"FrameType,omitempty"`
	DevName              string   `protobuf:"bytes,2,opt,name=DevName,proto3" json:"DevName,omitempty"`
	DevId                string   `protobuf:"bytes,3,opt,name=DevId,proto3" json:"DevId,omitempty"`
	FrameNum             uint32   `protobuf:"varint,4,opt,name=FrameNum,proto3" json:"FrameNum,omitempty"`
	Port                 uint32   `protobuf:"varint,5,opt,name=Port,proto3" json:"Port,omitempty"`
	PhyPayload           []byte   `protobuf:"bytes,6,opt,name=PhyPayload,proto3" json:"PhyPayload,omitempty"`
	UplinkId             uint32   `protobuf:"varint,7,opt,name=UplinkId,proto3" json:"UplinkId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UplinkFrame) Reset()         { *m = UplinkFrame{} }
func (m *UplinkFrame) String() string { return proto.CompactTextString(m) }
func (*UplinkFrame) ProtoMessage()    {}
func (*UplinkFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{0}
}

func (m *UplinkFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UplinkFrame.Unmarshal(m, b)
}
func (m *UplinkFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UplinkFrame.Marshal(b, m, deterministic)
}
func (m *UplinkFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UplinkFrame.Merge(m, src)
}
func (m *UplinkFrame) XXX_Size() int {
	return xxx_messageInfo_UplinkFrame.Size(m)
}
func (m *UplinkFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_UplinkFrame.DiscardUnknown(m)
}

var xxx_messageInfo_UplinkFrame proto.InternalMessageInfo

func (m *UplinkFrame) GetFrameType() int32 {
	if m != nil {
		return m.FrameType
	}
	return 0
}

func (m *UplinkFrame) GetDevName() string {
	if m != nil {
		return m.DevName
	}
	return ""
}

func (m *UplinkFrame) GetDevId() string {
	if m != nil {
		return m.DevId
	}
	return ""
}

func (m *UplinkFrame) GetFrameNum() uint32 {
	if m != nil {
		return m.FrameNum
	}
	return 0
}

func (m *UplinkFrame) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *UplinkFrame) GetPhyPayload() []byte {
	if m != nil {
		return m.PhyPayload
	}
	return nil
}

func (m *UplinkFrame) GetUplinkId() uint32 {
	if m != nil {
		return m.UplinkId
	}
	return 0
}

type DownlinkFrame struct {
	FrameType            int32    `protobuf:"varint,1,opt,name=FrameType,proto3" json:"FrameType,omitempty"`
	DevName              string   `protobuf:"bytes,2,opt,name=DevName,proto3" json:"DevName,omitempty"`
	DevId                string   `protobuf:"bytes,3,opt,name=DevId,proto3" json:"DevId,omitempty"`
	FrameNum             uint32   `protobuf:"varint,4,opt,name=FrameNum,proto3" json:"FrameNum,omitempty"`
	Port                 uint32   `protobuf:"varint,5,opt,name=Port,proto3" json:"Port,omitempty"`
	PhyPayload           []byte   `protobuf:"bytes,6,opt,name=PhyPayload,proto3" json:"PhyPayload,omitempty"`
	DownlinkId           uint32   `protobuf:"varint,7,opt,name=DownlinkId,proto3" json:"DownlinkId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DownlinkFrame) Reset()         { *m = DownlinkFrame{} }
func (m *DownlinkFrame) String() string { return proto.CompactTextString(m) }
func (*DownlinkFrame) ProtoMessage()    {}
func (*DownlinkFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{1}
}

func (m *DownlinkFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DownlinkFrame.Unmarshal(m, b)
}
func (m *DownlinkFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DownlinkFrame.Marshal(b, m, deterministic)
}
func (m *DownlinkFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownlinkFrame.Merge(m, src)
}
func (m *DownlinkFrame) XXX_Size() int {
	return xxx_messageInfo_DownlinkFrame.Size(m)
}
func (m *DownlinkFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_DownlinkFrame.DiscardUnknown(m)
}

var xxx_messageInfo_DownlinkFrame proto.InternalMessageInfo

func (m *DownlinkFrame) GetFrameType() int32 {
	if m != nil {
		return m.FrameType
	}
	return 0
}

func (m *DownlinkFrame) GetDevName() string {
	if m != nil {
		return m.DevName
	}
	return ""
}

func (m *DownlinkFrame) GetDevId() string {
	if m != nil {
		return m.DevId
	}
	return ""
}

func (m *DownlinkFrame) GetFrameNum() uint32 {
	if m != nil {
		return m.FrameNum
	}
	return 0
}

func (m *DownlinkFrame) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *DownlinkFrame) GetPhyPayload() []byte {
	if m != nil {
		return m.PhyPayload
	}
	return nil
}

func (m *DownlinkFrame) GetDownlinkId() uint32 {
	if m != nil {
		return m.DownlinkId
	}
	return 0
}

type Payload struct {
	Kind                 uint32   `protobuf:"varint,1,opt,name=kind,proto3" json:"kind,omitempty"`
	Key                  uint32   `protobuf:"varint,2,opt,name=key,proto3" json:"key,omitempty"`
	Val                  []byte   `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Payload) Reset()         { *m = Payload{} }
func (m *Payload) String() string { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()    {}
func (*Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1a937782ebbded5, []int{2}
}

func (m *Payload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Payload.Unmarshal(m, b)
}
func (m *Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Payload.Marshal(b, m, deterministic)
}
func (m *Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload.Merge(m, src)
}
func (m *Payload) XXX_Size() int {
	return xxx_messageInfo_Payload.Size(m)
}
func (m *Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_Payload proto.InternalMessageInfo

func (m *Payload) GetKind() uint32 {
	if m != nil {
		return m.Kind
	}
	return 0
}

func (m *Payload) GetKey() uint32 {
	if m != nil {
		return m.Key
	}
	return 0
}

func (m *Payload) GetVal() []byte {
	if m != nil {
		return m.Val
	}
	return nil
}

func init() {
	proto.RegisterEnum("proto.Category", Category_name, Category_value)
	proto.RegisterEnum("proto.Device", Device_name, Device_value)
	proto.RegisterEnum("proto.Operation", Operation_name, Operation_value)
	proto.RegisterType((*UplinkFrame)(nil), "proto.UplinkFrame")
	proto.RegisterType((*DownlinkFrame)(nil), "proto.DownlinkFrame")
	proto.RegisterType((*Payload)(nil), "proto.Payload")
}

func init() { proto.RegisterFile("gateway.proto", fileDescriptor_f1a937782ebbded5) }

var fileDescriptor_f1a937782ebbded5 = []byte{
	// 317 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x92, 0xdf, 0x6a, 0xfa, 0x30,
	0x1c, 0xc5, 0x8d, 0xf6, 0x8f, 0x7e, 0x7f, 0x16, 0xc2, 0x97, 0xdf, 0x45, 0x18, 0xe2, 0x8a, 0x57,
	0xc5, 0x8b, 0xdd, 0xec, 0x09, 0xc4, 0xda, 0x4d, 0x98, 0x6d, 0x09, 0xee, 0x01, 0xb2, 0x35, 0xb8,
	0xa2, 0x36, 0xa5, 0x74, 0x95, 0x3c, 0xe1, 0xae, 0xf6, 0x4e, 0xa3, 0x29, 0x75, 0xbe, 0xc2, 0xae,
	0x7a, 0xce, 0xa7, 0xc9, 0xe1, 0x1c, 0x08, 0x78, 0x07, 0x51, 0xcb, 0x8b, 0xd0, 0x0f, 0x65, 0xa5,
	0x6a, 0x85, 0xb6, 0xf9, 0x2c, 0xbe, 0x08, 0xfc, 0x7b, 0x2d, 0x4f, 0x79, 0x71, 0x8c, 0x2a, 0x71,
	0x96, 0x38, 0x83, 0x89, 0x11, 0x7b, 0x5d, 0x4a, 0x46, 0x7c, 0x12, 0xd8, 0xfc, 0x17, 0x20, 0x03,
	0x37, 0x94, 0x4d, 0x2c, 0xce, 0x92, 0x0d, 0x7d, 0x12, 0x4c, 0x78, 0x6f, 0xf1, 0x3f, 0xd8, 0xa1,
	0x6c, 0xb6, 0x19, 0x1b, 0x19, 0xde, 0x19, 0xbc, 0x83, 0xb1, 0xb9, 0x1c, 0x7f, 0x9e, 0x99, 0xe5,
	0x93, 0xc0, 0xe3, 0x57, 0x8f, 0x08, 0x56, 0xaa, 0xaa, 0x9a, 0xd9, 0x86, 0x1b, 0x8d, 0x73, 0x80,
	0xf4, 0x43, 0xa7, 0x42, 0x9f, 0x94, 0xc8, 0x98, 0xe3, 0x93, 0x60, 0xca, 0x6f, 0x48, 0x9b, 0xd7,
	0x95, 0xdd, 0x66, 0xcc, 0xed, 0xf2, 0x7a, 0xbf, 0xf8, 0x26, 0xe0, 0x85, 0xea, 0x52, 0xfc, 0x95,
	0x2d, 0x73, 0x80, 0xbe, 0xee, 0x75, 0xcd, 0x0d, 0x59, 0xac, 0xc0, 0xed, 0x8f, 0x22, 0x58, 0xc7,
	0xbc, 0xc8, 0xcc, 0x06, 0x8f, 0x1b, 0x8d, 0x14, 0x46, 0x47, 0xa9, 0x4d, 0x75, 0x8f, 0xb7, 0xb2,
	0x25, 0x8d, 0x38, 0x99, 0xd2, 0x53, 0xde, 0xca, 0xe5, 0x3d, 0x8c, 0xd7, 0xa2, 0x96, 0x07, 0x55,
	0x69, 0x74, 0x61, 0xb4, 0xde, 0x85, 0x74, 0x80, 0x00, 0xce, 0x3a, 0x89, 0xa3, 0xed, 0x13, 0x25,
	0xcb, 0x39, 0x38, 0xa1, 0x6c, 0xf2, 0x77, 0x89, 0x63, 0xb0, 0x5e, 0x56, 0xbb, 0xb4, 0xfb, 0xff,
	0xbc, 0x59, 0xed, 0x37, 0x9c, 0x92, 0xe5, 0x0c, 0x26, 0x49, 0x29, 0x2b, 0x51, 0xe7, 0xaa, 0x68,
	0x13, 0x92, 0x28, 0xa2, 0x03, 0x74, 0x60, 0x98, 0xc4, 0x94, 0xbc, 0x39, 0xe6, 0x09, 0x3d, 0xfe,
	0x04, 0x00, 0x00, 0xff, 0xff, 0x21, 0x96, 0x5b, 0x7d, 0x5a, 0x02, 0x00, 0x00,
}
