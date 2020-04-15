package streamreceiver

import (
	"context"
	"errors"
	"io"
	"reflect"

	analyser "github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockanalyser"
	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyserfiletransfer"
	crawl "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/sherlockcrawler"
)

/*
ServerGRPC receives a byte array from the crawler
*/
type ServerGRPC struct {
	Queue      *analyser.AnalyserQueue
	Dependency *ServerDependency
}

/*
ServerDependency will be all dependencys for the ServerGRPC.
*/
type ServerDependency struct {
	Crawler func()
}

/*
NewServerGRPC returns a new ServerGRPC
*/
func NewServerGRPC(lqueue *analyser.AnalyserQueue) *ServerGRPC {
	return &ServerGRPC{
		Queue: lqueue,
	}
}

/*
DownloadFile gets chunks of a html response from the crawler, appends them and returns the result
*/
func (handler *ServerGRPC) DownloadFile(ctx context.Context, stream proto.Sender_UploadStream) error {
	var task crawl.CrawlerTaskRequest

	rec, err := stream.Recv()
	switch reflect.TypeOf(rec) {
	case *proto.Infos:
		task.setTaskID(rec.TaskId)
		task.setAddr(rec.Address)
		task.setResponseHeader(rec.Header)
		task.setStatusCode(rec.StatusCode)
		task.setResponseTime(rec.ResponseTime)

		err = stream.SendMsg(&proto.UploadStatus{
			Code: proto.UploadStatusCode_Ok,
		})
		if err != nil {
			return errors.New("failed to send status code")
		}
	case *proto.ErrorCase:
		task.setTaskID(rec.TaskId)
		task.setAddr(rec.Address)
		task.setTaskError(errors.New(rec.TaskError))
		task.setResponseTime(rec.ResponseTime)

		err = stream.SendMsg(&proto.UploadStatus{
			Code: proto.UploadStatusCode_Ok,
		})
		if err != nil {
			return errors.New("failed to send status code")
		}
		return nil
	}

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

	err = stream.SendMsg(&proto.UploadStatus{
		Code: proto.UploadStatusCode_Ok,
	})
	if err != nil {
		return errors.New("failed to send status code")
	}

	defer stream.Close()

	//TODO werte weitergeben: arr und task
	return nil
}
