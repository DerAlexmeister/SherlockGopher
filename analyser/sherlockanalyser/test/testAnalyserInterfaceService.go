package test

import (
	"context"
	"github.com/micro/go-micro/client"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
)

type testAnalyserInterfaceService struct {
}

func (test testAnalyserInterfaceService) CreateTask(ctx context.Context, in *crawlerproto.CrawlTaskCreateRequest, opts ...client.CallOption) (*crawlerproto.CrawlTaskCreateResponse, error) {
	return nil, nil
}

func GetAnalyserInterfaceServiceInstance() crawlerproto.AnalyserInterfaceService {
	return testAnalyserInterfaceService{}
}


