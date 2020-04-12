// Code generated by protoc-gen-go. DO NOT EDIT.
// source: crawlerToAnalyser.proto

package analyserinterface

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

type URL_STATUS int32

const (
	URL_STATUS_ok      URL_STATUS = 0
	URL_STATUS_failure URL_STATUS = 1
)

var URL_STATUS_name = map[int32]string{
	0: "ok",
	1: "failure",
}

var URL_STATUS_value = map[string]int32{
	"ok":      0,
	"failure": 1,
}

func (x URL_STATUS) String() string {
	return proto.EnumName(URL_STATUS_name, int32(x))
}

func (URL_STATUS) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_17c4b44578c50216, []int{0}
}

type CrawlTaskCreateRequest struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CrawlTaskCreateRequest) Reset()         { *m = CrawlTaskCreateRequest{} }
func (m *CrawlTaskCreateRequest) String() string { return proto.CompactTextString(m) }
func (*CrawlTaskCreateRequest) ProtoMessage()    {}
func (*CrawlTaskCreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_17c4b44578c50216, []int{0}
}

func (m *CrawlTaskCreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrawlTaskCreateRequest.Unmarshal(m, b)
}
func (m *CrawlTaskCreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrawlTaskCreateRequest.Marshal(b, m, deterministic)
}
func (m *CrawlTaskCreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrawlTaskCreateRequest.Merge(m, src)
}
func (m *CrawlTaskCreateRequest) XXX_Size() int {
	return xxx_messageInfo_CrawlTaskCreateRequest.Size(m)
}
func (m *CrawlTaskCreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CrawlTaskCreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CrawlTaskCreateRequest proto.InternalMessageInfo

func (m *CrawlTaskCreateRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type CrawlTaskCreateResponse struct {
	Statuscode           URL_STATUS `protobuf:"varint,1,opt,name=statuscode,proto3,enum=analyserinterface.URL_STATUS" json:"statuscode,omitempty"`
	Taskid               uint64     `protobuf:"varint,2,opt,name=taskid,proto3" json:"taskid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *CrawlTaskCreateResponse) Reset()         { *m = CrawlTaskCreateResponse{} }
func (m *CrawlTaskCreateResponse) String() string { return proto.CompactTextString(m) }
func (*CrawlTaskCreateResponse) ProtoMessage()    {}
func (*CrawlTaskCreateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_17c4b44578c50216, []int{1}
}

func (m *CrawlTaskCreateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrawlTaskCreateResponse.Unmarshal(m, b)
}
func (m *CrawlTaskCreateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrawlTaskCreateResponse.Marshal(b, m, deterministic)
}
func (m *CrawlTaskCreateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrawlTaskCreateResponse.Merge(m, src)
}
func (m *CrawlTaskCreateResponse) XXX_Size() int {
	return xxx_messageInfo_CrawlTaskCreateResponse.Size(m)
}
func (m *CrawlTaskCreateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CrawlTaskCreateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CrawlTaskCreateResponse proto.InternalMessageInfo

func (m *CrawlTaskCreateResponse) GetStatuscode() URL_STATUS {
	if m != nil {
		return m.Statuscode
	}
	return URL_STATUS_ok
}

func (m *CrawlTaskCreateResponse) GetTaskid() uint64 {
	if m != nil {
		return m.Taskid
	}
	return 0
}

func init() {
	proto.RegisterEnum("analyserinterface.URL_STATUS", URL_STATUS_name, URL_STATUS_value)
	proto.RegisterType((*CrawlTaskCreateRequest)(nil), "analyserinterface.CrawlTaskCreateRequest")
	proto.RegisterType((*CrawlTaskCreateResponse)(nil), "analyserinterface.CrawlTaskCreateResponse")
}

func init() {
	proto.RegisterFile("crawlerToAnalyser.proto", fileDescriptor_17c4b44578c50216)
}

var fileDescriptor_17c4b44578c50216 = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x90, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x9b, 0x2a, 0x11, 0x47, 0x90, 0x74, 0x0e, 0x6d, 0x11, 0x04, 0xcd, 0x49, 0x73, 0xc8,
	0xa1, 0x9e, 0x3d, 0x94, 0x9e, 0x04, 0x4f, 0xdb, 0xed, 0x59, 0xc6, 0x74, 0x0a, 0xa1, 0x4b, 0x36,
	0xce, 0xee, 0x22, 0xfa, 0xeb, 0x6d, 0xda, 0x04, 0x85, 0xe4, 0xe0, 0x6d, 0x1f, 0xbc, 0x6f, 0xdf,
	0x7b, 0x03, 0xb3, 0x42, 0xe8, 0xd3, 0xb0, 0x68, 0xbb, 0xac, 0xc8, 0x7c, 0x39, 0x96, 0xbc, 0x16,
	0xeb, 0x2d, 0x4e, 0xa8, 0xd5, 0x65, 0xe5, 0x59, 0x76, 0x54, 0x70, 0x9a, 0xc1, 0x74, 0xd5, 0xb8,
	0x35, 0xb9, 0xfd, 0x4a, 0x98, 0x3c, 0x2b, 0xfe, 0x08, 0xec, 0x3c, 0x26, 0x70, 0x16, 0xc4, 0xcc,
	0xa3, 0xbb, 0xe8, 0xe1, 0x52, 0x35, 0xcf, 0xb4, 0x86, 0x59, 0xcf, 0xeb, 0x6a, 0x5b, 0x39, 0xc6,
	0x67, 0x00, 0xe7, 0xc9, 0x07, 0x57, 0xd8, 0x2d, 0x1f, 0x99, 0xeb, 0xc5, 0x6d, 0xde, 0x8b, 0xcb,
	0x37, 0xea, 0xf5, 0x6d, 0xad, 0x97, 0x7a, 0xb3, 0x56, 0x7f, 0x00, 0x9c, 0x42, 0xec, 0x0f, 0x9f,
	0x96, 0xdb, 0xf9, 0xf8, 0x80, 0x9e, 0xab, 0x56, 0x65, 0xf7, 0x00, 0xbf, 0x04, 0xc6, 0x30, 0xb6,
	0xfb, 0x64, 0x84, 0x57, 0x70, 0xb1, 0xa3, 0xd2, 0x04, 0xe1, 0x24, 0x5a, 0x7c, 0xc3, 0xa4, 0x5b,
	0xf9, 0xd2, 0xc5, 0x20, 0x03, 0x9c, 0x0a, 0x36, 0x55, 0xf1, 0x71, 0xa0, 0xc8, 0xf0, 0xe8, 0x9b,
	0xec, 0x3f, 0xd6, 0xd3, 0xe6, 0x74, 0xf4, 0x1e, 0x1f, 0xcf, 0xfa, 0xf4, 0x13, 0x00, 0x00, 0xff,
	0xff, 0xd4, 0x4f, 0xa8, 0x46, 0x71, 0x01, 0x00, 0x00,
}
