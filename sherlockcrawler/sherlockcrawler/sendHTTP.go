package sherlockcrawler

import (
	"context"
	"io"
	"os"

	"github.com/micro/go-micro"
	sender "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/sherlockcrawler/proto"
	"github.com/pkg/errors"
)

const chunkSize int = 10
const filepath string = "testdatei.txt"

type ClientGRPC struct {
	client    sender.SenderService
	chunkSize int
}

func readF() *os.File {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	return file
}

func NewClientGRPC(service micro.Service) (c ClientGRPC) {
	c.chunkSize = chunkSize
	c.client = sender.NewSenderService("client", service.Client())
	return c
}

func (c *ClientGRPC) UploadFile(ctx context.Context) (err error) {
	writing := true
	file := readF()

	// Open a stream-based connection with the gRPC server
	stream, err := c.client.Upload(ctx)

	if err != nil {
		err = errors.Wrapf(err, "failed to create upload stream for file %s", file)
		return
	}

	buf := make([]byte, chunkSize)
	var n int
	for writing {
		n, err = file.Read(buf)

		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			err = errors.Wrapf(err, "errored while copying from file to buf")
			return
		}

		err = stream.Send(&sender.Chunk{
			Content: buf[:n],
		})

		if err != nil {
			err = errors.Wrapf(err, "failed to send chunk via stream")
			return
		}
	}

	//receive

	err = stream.Close()

	if err != nil {
		err = errors.Wrapf(err, "failed to receive upstream status response")
		return
	}
	return
}
