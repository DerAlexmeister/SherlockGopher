package sherlockanalyser

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyserfiletransfer"
)

/*
ServerGRPC receives a byte array from the crawler
*/
type ServerGRPC struct {
	Queue      *AnalyserQueue
	Dependency *ServerDependency
}

/*
getCAddr getter for the address.
*/
func (s ServerGRPC) getQueue() *AnalyserQueue {
	return s.Queue
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
func NewServerGRPC(lqueue *AnalyserQueue) *ServerGRPC {
	return &ServerGRPC{
		Queue: lqueue,
	}
}

/*
DownloadFile gets chunks of a html response from the crawler, appends them and returns the result
*/
func (handler *ServerGRPC) Upload(ctx context.Context, stream proto.Sender_UploadStream) error {
	var task CrawlerData
	var arr []byte
	finished := false

	rec, err := stream.Recv()
	switch {
	case rec.TaskError == "" && rec.Address != "":
		headermap := http.Header{}
		for _, i := range rec.Header {
			for _, n := range i.ValueArr {
				headermap.Add(i.Key, n.Value)
			}
		}

		task.setCTaskID(rec.TaskId)
		task.setCAddr(rec.Address)
		task.setCResponseHeader(&headermap)
		task.setCStatusCode(int(rec.StatusCode))
		task.setCResponseTime(time.Duration(rec.ResponseTime))

		err = stream.Send(&proto.UploadStatus{
			Code: proto.UploadStatusCode_Ok,
		})
		if err != nil {
			return errors.New("failed to send status code")
		}
	case rec.TaskError != "":
		task.setCTaskID(rec.TaskId)
		task.setCAddr(rec.Address)
		task.setCTaskError(errors.New(rec.TaskError))
		task.setCResponseTime(time.Duration(rec.ResponseTime))

		err = stream.Send(&proto.UploadStatus{
			Code: proto.UploadStatusCode_Ok,
		})
		if err != nil {
			return errors.New("failed to send status code")
		}
		newTask := NewTask(task)
		handler.getQueue().AppendQueue(&newTask)
		return nil
	default:
		for !finished {
			arr = append(arr, rec.Content...)

			if err != nil {
				if err == io.EOF {
					finished = true
				} else {
					return errors.New("failed unexpectedely while reading chunks from stream")
				}
			}
			rec, err = stream.Recv()
		}
		task.setCResponseBody(arr)

		err = stream.Send(&proto.UploadStatus{
			Code: proto.UploadStatusCode_Ok,
		})
		if err != nil {
			return errors.New("failed to send status code")
		}
	}

	defer stream.Close()

	newTask := NewTask(task)
	handler.getQueue().AppendQueue(&newTask)

	return nil
}
