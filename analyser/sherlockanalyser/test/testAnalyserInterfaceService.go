package test

import (
	"context"

	"github.com/micro/go-micro/client"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto"
)

type TAnalyserInterfaceService struct {
}

func (test TAnalyserInterfaceService) ReceiveURL(ctx context.Context, in *crawlerproto.SubmitURLRequest, opts ...client.CallOption) (*crawlerproto.SubmitURLResponse, error) {
	return nil, nil
}

func (test TAnalyserInterfaceService) StatusOfTaskQueue(ctx context.Context, in *crawlerproto.TaskStatusRequest, opts ...client.CallOption) (*crawlerproto.TaskStatusResponse, error) {
	return nil, nil
}

func (test TAnalyserInterfaceService) SetState(ctx context.Context, in *crawlerproto.StateRequest, opts ...client.CallOption) (*crawlerproto.StateResponse, error) {
	return nil, nil
}

func (test TAnalyserInterfaceService) GetState(ctx context.Context, in *crawlerproto.StateGetRequest, opts ...client.CallOption) (*crawlerproto.StateGetResponse, error) {
	return nil, nil
}

/*
CreateTask.
*/
func (test TAnalyserInterfaceService) CreateTask(ctx context.Context, in *crawlerproto.CrawlTaskCreateRequest, opts ...client.CallOption) (*crawlerproto.CrawlTaskCreateResponse, error) {
	return nil, nil
}

/*
GetAnalyserInterfaceServiceInstance.
*/
func GetAnalyserInterfaceServiceInstance() crawlerproto.CrawlerService {
	return TAnalyserInterfaceService{}
}
