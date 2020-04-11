package streamreceiver

import (
	"context"
	"errors"
	"io"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto/filestreamproto"
)

/*
ServerGRPC receives a byte array from the crawler
*/
type ServerGRPC struct{}

/*
NewServerGRPC returns a new ServerGRPC
*/
func NewServerGRPC() *ServerGRPC {
	return &ServerGRPC{}
}

/*
DownloadFile gets chunks of a html response from the crawler, appends them and returns the result
*/
func (handler *ServerGRPC) DownloadFile(ctx context.Context, stream proto.Receiver_UploadStream) (arr []byte, err error) {
	finished := false

	for !finished {
		chunk, err := stream.Recv()
		arr = append(arr, chunk.Content...)

		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				return nil, errors.New("failed unexpectedely while reading chunks from stream")
			}
		}
	}

	err = stream.SendMsg(&proto.UploadStatus{
		Message: "Upload received with success",
		Code:    proto.UploadStatusCode_Ok,
	})
	if err != nil {
		return nil, errors.New("failed to send status code")
	}

	err = stream.Close()

	if err != nil {
		return nil, errors.New("error while closing stream")
	}

	return arr, nil
}
