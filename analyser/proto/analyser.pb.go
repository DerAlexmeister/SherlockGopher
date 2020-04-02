// Code generated by protoc-gen-go. DO NOT EDIT.
// source: analyser.proto

package messaging

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

type UploadStatusCode int32

const (
	UploadStatusCode_Unknown UploadStatusCode = 0
	UploadStatusCode_Ok      UploadStatusCode = 1
	UploadStatusCode_Failed  UploadStatusCode = 2
)

var UploadStatusCode_name = map[int32]string{
	0: "Unknown",
	1: "Ok",
	2: "Failed",
}

var UploadStatusCode_value = map[string]int32{
	"Unknown": 0,
	"Ok":      1,
	"Failed":  2,
}

func (x UploadStatusCode) String() string {
	return proto.EnumName(UploadStatusCode_name, int32(x))
}

func (UploadStatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_18bf88fceba00cff, []int{0}
}

type Chunk struct {
	Content              []byte   `protobuf:"bytes,1,opt,name=Content,proto3" json:"Content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Chunk) Reset()         { *m = Chunk{} }
func (m *Chunk) String() string { return proto.CompactTextString(m) }
func (*Chunk) ProtoMessage()    {}
func (*Chunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_18bf88fceba00cff, []int{0}
}

func (m *Chunk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Chunk.Unmarshal(m, b)
}
func (m *Chunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Chunk.Marshal(b, m, deterministic)
}
func (m *Chunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Chunk.Merge(m, src)
}
func (m *Chunk) XXX_Size() int {
	return xxx_messageInfo_Chunk.Size(m)
}
func (m *Chunk) XXX_DiscardUnknown() {
	xxx_messageInfo_Chunk.DiscardUnknown(m)
}

var xxx_messageInfo_Chunk proto.InternalMessageInfo

func (m *Chunk) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type UploadStatus struct {
	Message              string           `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	Code                 UploadStatusCode `protobuf:"varint,2,opt,name=Code,proto3,enum=messaging.UploadStatusCode" json:"Code,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *UploadStatus) Reset()         { *m = UploadStatus{} }
func (m *UploadStatus) String() string { return proto.CompactTextString(m) }
func (*UploadStatus) ProtoMessage()    {}
func (*UploadStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_18bf88fceba00cff, []int{1}
}

func (m *UploadStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadStatus.Unmarshal(m, b)
}
func (m *UploadStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadStatus.Marshal(b, m, deterministic)
}
func (m *UploadStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadStatus.Merge(m, src)
}
func (m *UploadStatus) XXX_Size() int {
	return xxx_messageInfo_UploadStatus.Size(m)
}
func (m *UploadStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadStatus.DiscardUnknown(m)
}

var xxx_messageInfo_UploadStatus proto.InternalMessageInfo

func (m *UploadStatus) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *UploadStatus) GetCode() UploadStatusCode {
	if m != nil {
		return m.Code
	}
	return UploadStatusCode_Unknown
}

func init() {
	proto.RegisterEnum("messaging.UploadStatusCode", UploadStatusCode_name, UploadStatusCode_value)
	proto.RegisterType((*Chunk)(nil), "messaging.Chunk")
	proto.RegisterType((*UploadStatus)(nil), "messaging.UploadStatus")
}

func init() {
	proto.RegisterFile("analyser.proto", fileDescriptor_18bf88fceba00cff)
}

var fileDescriptor_18bf88fceba00cff = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0xcc, 0x4b, 0xcc,
	0xa9, 0x2c, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xcc, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0xcf, 0xcc, 0x4b, 0x57, 0x52, 0xe4, 0x62, 0x75, 0xce, 0x28, 0xcd, 0xcb, 0x16, 0x92,
	0xe0, 0x62, 0x77, 0xce, 0xcf, 0x2b, 0x49, 0xcd, 0x2b, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x09,
	0x82, 0x71, 0x95, 0x22, 0xb9, 0x78, 0x42, 0x0b, 0x72, 0xf2, 0x13, 0x53, 0x82, 0x4b, 0x12, 0x4b,
	0x4a, 0x8b, 0x41, 0x2a, 0x7d, 0xc1, 0xfa, 0x53, 0xc1, 0x2a, 0x39, 0x83, 0x60, 0x5c, 0x21, 0x7d,
	0x2e, 0x16, 0xe7, 0xfc, 0x94, 0x54, 0x09, 0x26, 0xa0, 0x30, 0x9f, 0x91, 0xb4, 0x1e, 0xdc, 0x1a,
	0x3d, 0x64, 0x03, 0x40, 0x4a, 0x82, 0xc0, 0x0a, 0xb5, 0x8c, 0xb9, 0x04, 0xd0, 0x65, 0x84, 0xb8,
	0xb9, 0xd8, 0x43, 0xf3, 0xb2, 0xf3, 0xf2, 0xcb, 0xf3, 0x04, 0x18, 0x84, 0xd8, 0xb8, 0x98, 0xfc,
	0xb3, 0x05, 0x18, 0x85, 0xb8, 0xb8, 0xd8, 0xdc, 0x12, 0x33, 0x73, 0x52, 0x53, 0x04, 0x98, 0x8c,
	0x9c, 0xb9, 0x38, 0x1c, 0xa1, 0xfe, 0x11, 0x32, 0xe7, 0x62, 0x83, 0x18, 0x20, 0x24, 0x80, 0x64,
	0x1b, 0xd8, 0x47, 0x52, 0xe2, 0x38, 0xec, 0x57, 0x62, 0xd0, 0x60, 0x4c, 0x62, 0x03, 0x87, 0x84,
	0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xce, 0x8f, 0x6b, 0x0c, 0x1b, 0x01, 0x00, 0x00,
}
