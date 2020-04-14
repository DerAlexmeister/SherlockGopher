package streamreceiver

import (
	"context"
	"errors"
	"io"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyserfiletransfer"
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
func (handler *ServerGRPC) DownloadFile(ctx context.Context, stream proto.Sender_UploadStream) (err error) {
	//errorcase oder infos

	//File erhalten
	var arr []byte
	finished := false

	for !finished {
		chunk, err := stream.Recv()
		arr = append(arr, chunk.Content...)

		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				return errors.New("failed unexpectedely while reading chunks from stream")
			}
		}
	}

	err = stream.Close()

	if err != nil {
		return errors.New("error while closing stream")
	}

	return nil
}

func (handler *ServerGRPC) UploadInfos(ctx context.Context, in *proto.Infos, out *proto.UploadStatus) error {

}

func (handler *ServerGRPC) UploadErrorCase(ctx context.Context, in *proto.ErrorCase, out *proto.UploadStatus) error {

}

//TODO werte weitergeben
