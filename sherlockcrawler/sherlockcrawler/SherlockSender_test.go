package sherlockcrawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"

	aproto "github.com/DerAlexx/SherlockGopher/analyser/proto"
	mockanalyser "github.com/DerAlexx/SherlockGopher/sherlockcrawler/sherlockcrawler/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

const (
	errorMessage = "This error is expected!"
)

type SenderTestError struct{}

func (err *SenderTestError) Error() string {
	return errorMessage
}

func TestSendWebsiteData(t *testing.T) {
	sut := NewSherlockCrawlerService()

	ctrl := gomock.NewController(t)

	mockAnalyser := mockanalyser.NewMockAnalyserService(ctrl)
	mockWebsiteData := mockanalyser.NewMockAnalyser_WebsiteDataService(ctrl)
	mockAnalyser.EXPECT().WebsiteData(gomock.Any()).Return(mockWebsiteData, nil)
	deps := SherlockDependencies{Analyser: func() aproto.AnalyserService { return mockAnalyser }}
	sut.InjectDependency(&deps)
	mockWebsiteData.EXPECT().Send(gomock.Any()).Return(nil).MinTimes(1)
	mockWebsiteData.EXPECT().Recv().Return(&aproto.CrawlerAck{
		Id: 0,
	}, nil).MinTimes(1)
	mockWebsiteData.EXPECT().Close()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	header := http.Header{}
	header.Add("1", "example.org")
	have := sut.SendWebsiteData(mockAnalyser, &CrawlerTaskRequest{
		taskID:       0,
		addr:         "https://www.hm.edu",
		taskState:    0,
		taskError:    nil,
		taskErrorTry: 0,
		response: &http.Response{
			Status:           "",
			StatusCode:       0,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           header,
			Body:             ioutil.NopCloser(strings.NewReader("body")),
			ContentLength:    0,
			TransferEncoding: nil,
			Close:            false,
			Uncompressed:     false,
			Trailer:          nil,
			Request:          nil,
			TLS:              nil,
		},
		responseHeader:    &header,
		responseBody:      "This is a test body",
		responseBodyBytes: nil,
		statusCode:        0,
		responseTime:      0,
	}, wg)
	fmt.Print(have)
	if have != nil {
		assert.Equal(t, have, nil)
	}
}

func TestSendWebsiteDataError(t *testing.T) {
	sut := NewSherlockCrawlerService()

	ctrl := gomock.NewController(t)

	mockAnalyser := mockanalyser.NewMockAnalyserService(ctrl)
	mockWebsiteData := mockanalyser.NewMockAnalyser_WebsiteDataService(ctrl)
	mockAnalyser.EXPECT().WebsiteData(gomock.Any()).Return(mockWebsiteData, nil)
	deps := SherlockDependencies{Analyser: func() aproto.AnalyserService { return mockAnalyser }}
	sut.InjectDependency(&deps)
	mockWebsiteData.EXPECT().Send(gomock.Any()).Return(nil).MinTimes(1)
	mockWebsiteData.EXPECT().Recv().Return(&aproto.CrawlerAck{
		Id: 0,
	}, nil).MinTimes(1)
	mockWebsiteData.EXPECT().Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	header := http.Header{}
	header.Add("2", "example.org")
	have := sut.SendWebsiteData(mockAnalyser, &CrawlerTaskRequest{
		taskID:       0,
		addr:         "https://www.hm.edu",
		taskState:    0,
		taskError:    &SenderTestError{},
		taskErrorTry: 0,
		response: &http.Response{
			Status:           "",
			StatusCode:       0,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           header,
			Body:             ioutil.NopCloser(strings.NewReader("body")),
			ContentLength:    0,
			TransferEncoding: nil,
			Close:            false,
			Uncompressed:     false,
			Trailer:          nil,
			Request:          nil,
			TLS:              nil,
		},
		responseHeader:    &header,
		responseBody:      "",
		responseBodyBytes: nil,
		statusCode:        0,
		responseTime:      0,
	}, wg)
	fmt.Print(have)
	if have != nil {
		assert.Equal(t, have, nil)
	}
}

func TestSendDataWebsiteDataError(t *testing.T) {
	sut := NewSherlockCrawlerService()

	ctrl := gomock.NewController(t)

	mockAnalyser := mockanalyser.NewMockAnalyserService(ctrl)

	mockAnalyser.EXPECT().WebsiteData(gomock.Any()).Return(nil, &SenderTestError{})
	deps := SherlockDependencies{Analyser: func() aproto.AnalyserService { return mockAnalyser }}
	sut.InjectDependency(&deps)
	err := sut.SendData(1, nil, nil, nil)
	if err.Error() != errorMessage {
		t.Fail()
	}
}

func TestSendDataSendError(t *testing.T) {
	sut := NewSherlockCrawlerService()

	ctrl := gomock.NewController(t)

	mockAnalyser := mockanalyser.NewMockAnalyserService(ctrl)
	mockWebsiteData := mockanalyser.NewMockAnalyser_WebsiteDataService(ctrl)
	mockAnalyser.EXPECT().WebsiteData(gomock.Any()).Return(mockWebsiteData, nil)
	deps := SherlockDependencies{Analyser: func() aproto.AnalyserService { return mockAnalyser }}
	sut.InjectDependency(&deps)
	mockWebsiteData.EXPECT().Send(gomock.Any()).Return(&SenderTestError{}).MinTimes(1)
	err := sut.SendData(1, [][]byte{{2}}, nil, &CrawlerTaskRequest{addr: "example.com"})

	if err.Error() != errorMessage {
		t.Fail()
	}
}

func TestSendReceiveError(t *testing.T) {
	sut := NewSherlockCrawlerService()

	ctrl := gomock.NewController(t)

	mockAnalyser := mockanalyser.NewMockAnalyserService(ctrl)
	mockWebsiteData := mockanalyser.NewMockAnalyser_WebsiteDataService(ctrl)
	mockAnalyser.EXPECT().WebsiteData(gomock.Any()).Return(mockWebsiteData, nil)
	deps := SherlockDependencies{Analyser: func() aproto.AnalyserService { return mockAnalyser }}
	sut.InjectDependency(&deps)
	mockWebsiteData.EXPECT().Send(gomock.Any()).Return(nil).MinTimes(1)
	mockWebsiteData.EXPECT().Recv().Return(nil, &SenderTestError{}).MinTimes(1)
	mockWebsiteData.EXPECT().Close()

	err := sut.SendData(1, [][]byte{{2}}, nil, &CrawlerTaskRequest{addr: "example.com"})

	if err.Error() != errorMessage {
		t.Fail()
	}
}
