package sherlockcrawler

import (
	"context"
	"fmt"
	"time"

	sender "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyserfiletransfer"
	"github.com/pkg/errors"
)

/*
CHUNKSIZE represents the size of a chunk.
It is necessary to minimize the amount of bytes sent at once.
*/
const (
	CHUNKSIZE   int = 1024
	sendingtime     = 100
)

/*
SherlockStreamingServer will be the representation
of the StreamingServer pushing the files to the Analyser.
*/
type SherlockStreamingServer struct {
	Client sender.SenderService
	Queue  CrawlerQueue
}

/*
getChunkSize will return the chunksize for a chunk via fileupload.
*/
func getChunkSize() int {
	return CHUNKSIZE
}

/*
getQueue will return the current queue of the SherlockStreamingServer instance.
*/
func (c *SherlockStreamingServer) getQueue() *CrawlerQueue {
	return &c.Queue
}

/*
NewStreamingServer creates a new sender for the files collected by the crawler.
*/
func NewStreamingServer() (c SherlockStreamingServer) {
	return SherlockStreamingServer{
		Queue: NewCrawlerQueue(),
	}
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
sendFileToAnalyser will send a file + additional information to the analyser.
*/
func (c *SherlockStreamingServer) sendFileToAnalyser(ctx context.Context, ltask *CrawlerTaskRequest) error {

	stream, err := c.Client.Upload(ctx)
	if err != nil {
		return errors.New("failed to create upload stream for file")
	}

	if ltask.taskerror != nil {
		err = helpSendErrorCase(ctx, ltask, stream)
		if err != nil {
			return err
		}
	} else {
		err = helpSendInfos(ctx, ltask, stream)
		if err != nil {
			return err
		}
		err = helpSend(ctx, ltask, stream)
		if err != nil {
			return err
		}
	}

	err = stream.Close()
	if err != nil {
		return errors.New("error while closing stream")
	}

	c.getQueue().RemoveFromQueue(ltask.getTaskID())
	return nil
}

/*
helpSendErrorCase is a help method to reduce the size of the sendFileToAnalyser function. it is used in case there is an error.
*/
func helpSendErrorCase(ctx context.Context, ltask *CrawlerTaskRequest, stream sender.Sender_UploadService) (err error) {
	err = stream.SendMsg(&sender.ErrorCase{
		TaskId:       ltask.getTaskID(),
		Address:      ltask.getAddr(),
		TaskError:    ltask.getTaskError().Error(),
		ResponseTime: int64(ltask.getResponseTime()),
	})

	if err != nil {
		return errors.New("error while sending info message")
	}

	var status *sender.UploadStatus
	if status.Code != sender.UploadStatusCode_Ok {
		return errors.Errorf("upload error case failed - msg: %s", status.Message)
	}

	return err
}

/*
helpSendInfos is a help method to reduce the size of the sendFileToAnalyser function. it is responsible for sending the additional information. the analyser requires the information once.
*/
func helpSendInfos(ctx context.Context, ltask *CrawlerTaskRequest, stream sender.Sender_UploadService) (err error) {

	headerArr := []*sender.HeaderArray{}
	valueArr := []*sender.HeaderArrayValue{}
	for k, v := range ltask.getResponseHeader() {
		for _, i := range v {
			valueArr = append(valueArr, &sender.HeaderArrayValue{Value: i})
		}
		headerArr = append(headerArr, &sender.HeaderArray{Key: k, ValueArr: valueArr})
	}

	err = stream.SendMsg(&sender.Infos{
		TaskId:       ltask.getTaskID(),
		Address:      ltask.getAddr(),
		Header:       headerArr,
		StatusCode:   int32(ltask.getStatusCode()),
		ResponseTime: int64(ltask.getResponseTime()),
	})

	if err != nil {
		return errors.New("error while sending info message")
	}

	var status *sender.UploadStatus
	if status.Code != sender.UploadStatusCode_Ok {
		return errors.Errorf("upload infos failed - msg: %s", status.Message)
	}

	return err
}

/*
helpSend is a help method to reduce the size of the sendFileToAnalyser function. it is responsible for sending the http response with a taskid.
*/
func helpSend(ctx context.Context, ltask *CrawlerTaskRequest, stream sender.Sender_UploadService) (err error) {

	var lengthByteArray int = len(ltask.getResponseBodyInBytes())

	for i := 0; i < lengthByteArray; i += getChunkSize() {
		buf := ltask.getResponseBodyInBytes()[i:min(i+getChunkSize(), lengthByteArray)]

		err = stream.Send(&sender.Chunk{
			Content: buf,
			TaskId:  ltask.getTaskID(),
		})

		if err != nil {
			return errors.New("error while streaming")
		}
	}

	var status *sender.UploadStatus
	if status.Code != sender.UploadStatusCode_Ok {
		return errors.Errorf("upload failed - msg: %s", status.Message)
	}

	return err
}

/*
UploadFile cuts byte array in slices of chunksize and sends them to the analyzer.
*/
func (c *SherlockStreamingServer) Upload(ctx context.Context) error {
	for {
		for _, task := range *c.getQueue().getThisQueue() {
			if err := c.sendFileToAnalyser(ctx, task); err != nil {
				fmt.Printf("An error occurred while trying to submit a file. Error: %s ", err.Error())
			}
		}
		time.Tick(sendingtime * time.Millisecond)
	}
}
