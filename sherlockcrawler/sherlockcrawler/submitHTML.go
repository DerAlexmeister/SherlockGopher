package sherlockcrawler

import (
	"context"
	"io"
	"os"

	"github.com/micro/go-micro"
	sender "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/sherlockcrawler/proto"
	"github.com/pkg/errors"
)

const chunkSize int = 1024
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

func (c *ClientGRPC) UploadFile(ctx context.Context) error {
	writing := true
	file := readF()

	stream, err := c.client.Upload(ctx)

	if err != nil {
		return errors.New("failed to create upload stream for file")
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
			return errors.New("error while copying from file to buf")
		}

		err = stream.Send(&sender.Chunk{
			Content: buf[:n],
		})

		if err != nil {
			return errors.New("error while streaming")
		}
	}

	//receive fehlt

	err = stream.Close()

	if err != nil {
		return errors.New("error while closing stream"))
	}
	return nil
}
