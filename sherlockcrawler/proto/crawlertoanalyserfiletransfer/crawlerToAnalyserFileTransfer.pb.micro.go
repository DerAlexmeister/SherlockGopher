// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: crawlerToAnalyserFileTransfer.proto

package crawlerproto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Sender service

type SenderService interface {
	Upload(ctx context.Context, opts ...client.CallOption) (Sender_UploadService, error)
	UploadInfos(ctx context.Context, in *Infos, opts ...client.CallOption) (*UploadStatus, error)
	UploadErrorCase(ctx context.Context, in *ErrorCase, opts ...client.CallOption) (*UploadStatus, error)
}

type senderService struct {
	c    client.Client
	name string
}

func NewSenderService(name string, c client.Client) SenderService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "crawlerproto"
	}
	return &senderService{
		c:    c,
		name: name,
	}
}

func (c *senderService) Upload(ctx context.Context, opts ...client.CallOption) (Sender_UploadService, error) {
	req := c.c.NewRequest(c.name, "Sender.Upload", &Chunk{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &senderServiceUpload{stream}, nil
}

type Sender_UploadService interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Chunk) error
}

type senderServiceUpload struct {
	stream client.Stream
}

func (x *senderServiceUpload) Close() error {
	return x.stream.Close()
}

func (x *senderServiceUpload) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *senderServiceUpload) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *senderServiceUpload) Send(m *Chunk) error {
	return x.stream.Send(m)
}

func (c *senderService) UploadInfos(ctx context.Context, in *Infos, opts ...client.CallOption) (*UploadStatus, error) {
	req := c.c.NewRequest(c.name, "Sender.UploadInfos", in)
	out := new(UploadStatus)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *senderService) UploadErrorCase(ctx context.Context, in *ErrorCase, opts ...client.CallOption) (*UploadStatus, error) {
	req := c.c.NewRequest(c.name, "Sender.UploadErrorCase", in)
	out := new(UploadStatus)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Sender service

type SenderHandler interface {
	Upload(context.Context, Sender_UploadStream) error
	UploadInfos(context.Context, *Infos, *UploadStatus) error
	UploadErrorCase(context.Context, *ErrorCase, *UploadStatus) error
}

func RegisterSenderHandler(s server.Server, hdlr SenderHandler, opts ...server.HandlerOption) error {
	type sender interface {
		Upload(ctx context.Context, stream server.Stream) error
		UploadInfos(ctx context.Context, in *Infos, out *UploadStatus) error
		UploadErrorCase(ctx context.Context, in *ErrorCase, out *UploadStatus) error
	}
	type Sender struct {
		sender
	}
	h := &senderHandler{hdlr}
	return s.Handle(s.NewHandler(&Sender{h}, opts...))
}

type senderHandler struct {
	SenderHandler
}

func (h *senderHandler) Upload(ctx context.Context, stream server.Stream) error {
	return h.SenderHandler.Upload(ctx, &senderUploadStream{stream})
}

type Sender_UploadStream interface {
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*Chunk, error)
}

type senderUploadStream struct {
	stream server.Stream
}

func (x *senderUploadStream) Close() error {
	return x.stream.Close()
}

func (x *senderUploadStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *senderUploadStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *senderUploadStream) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *senderHandler) UploadInfos(ctx context.Context, in *Infos, out *UploadStatus) error {
	return h.SenderHandler.UploadInfos(ctx, in, out)
}

func (h *senderHandler) UploadErrorCase(ctx context.Context, in *ErrorCase, out *UploadStatus) error {
	return h.SenderHandler.UploadErrorCase(ctx, in, out)
}
