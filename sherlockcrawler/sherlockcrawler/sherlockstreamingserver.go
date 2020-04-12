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
sendFileToAnalyser will send a file to the analyser.
*/
func (c *SherlockStreamingServer) sendFileToAnalyser(ctx context.Context, ltask *CrawlerTaskRequest, taskid uint64) error {

	stream, err := c.Client.Upload(ctx)
	if err != nil {
		return errors.New("failed to create upload stream for file")
	}

	if ltask.taskerror != nil {
		err = helpSendErrorCase(ctx, ltask, taskid, stream)
		if err != nil {
			return err
		}
	} else {
		err = helpSendInfos(ctx, ltask, taskid, stream)
		if err != nil {
			return err
		}
		err = helpSend(ctx, ltask, taskid, stream)
		if err != nil {
			return err
		}
	}

	err = stream.Close()
	if err != nil {
		return errors.New("error while closing stream")
	}

	//TODO remove from queue on success
	return nil
}

func helpSendErrorCase(ctx context.Context, ltask *CrawlerTaskRequest, taskid uint64, stream sender.Sender_UploadService) (err error) {
	//TODO send errorcase

	var status *sender.UploadStatus
	if status.Code != sender.UploadStatusCode_Ok {
		return errors.Errorf("upload failed - msg: %s", status.Message)
	}

	return err
}

func helpSendInfos(ctx context.Context, ltask *CrawlerTaskRequest, taskid uint64, stream sender.Sender_UploadService) (err error) {
	//TODO send infocase

	var status *sender.UploadStatus
	if status.Code != sender.UploadStatusCode_Ok {
		return errors.Errorf("upload failed - msg: %s", status.Message)
	}

	return err
}

func helpSend(ctx context.Context, ltask *CrawlerTaskRequest, taskid uint64, stream sender.Sender_UploadService) (err error) {

	//TODO send additional infos

	var lengthByteArray int = len(ltask.getResponseBodyInBytes())

	for i := 0; i < lengthByteArray; i += getChunkSize() {
		buf := ltask.getResponseBodyInBytes()[i:min(i+getChunkSize(), lengthByteArray)]

		err = stream.Send(&sender.Chunk{
			Content: buf,
			TaskId:  taskid,
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
Upload cuts byte array in slices of chunksize and sends them to the analyzer.
*/
func (c *SherlockStreamingServer) Upload(ctx context.Context) error {
	for {
		for id, task := range *c.getQueue().getThisQueue() {
			if err := c.sendFileToAnalyser(ctx, task, id); err != nil {
				fmt.Printf("An error occurred while trying to submit a file. Error: %s ", err.Error())
			}
		}
		time.Tick(sendingtime * time.Millisecond)
	}
}
