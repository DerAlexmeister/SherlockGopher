package sherlockcrawler

import (
	"context"

	"github.com/micro/go-micro"
	sender "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto"
	"github.com/pkg/errors"
)

/*
chunkSize represents the size of a chunk. it is necessary to minimize the amount of bytes sent at once.
*/
const chunkSize int = 1024

/*
ClientGRPC will be the Client/Sender.
*/
type ClientGRPC struct {
	client    sender.SenderService
	chunkSize int
}

/*
NewClientGRPC creates a new sender.
*/
func NewClientGRPC(service micro.Service) (c ClientGRPC) {
	c.chunkSize = chunkSize
	c.client = sender.NewSenderService("client", service.Client())

	return c
}

/*
min is a help method. it receives 2 values and returns the smaller one.
*/
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

/*
UploadFile cuts byte array in slices of chunksize and sends them to the analyzer.
*/
func (c *ClientGRPC) UploadFile(ctx context.Context, ltask CrawlerTaskRequest) error {
	stream, err := c.client.Upload(ctx)

	if err != nil {
		return errors.New("failed to create upload stream for file")
	}

	var lengthByteArray int = len(ltask.getResponseBodyInBytes())

	for i := 0; i < lengthByteArray; i += chunkSize {
		buf := ltask.getResponseBodyInBytes()[i:min(i+chunkSize, lengthByteArray)]

		err = stream.Send(&sender.Chunk{
			Content: buf,
		})

		if err != nil {
			return errors.New("error while streaming")
		}
	}

	var status *sender.UploadStatus

	if status.Code != sender.UploadStatusCode_Ok {
		return errors.Errorf("upload failed - msg: %s", status.Message)
	}

	err = stream.Close()

	if err != nil {
		return errors.New("error while closing stream")
	}
	return nil
}
