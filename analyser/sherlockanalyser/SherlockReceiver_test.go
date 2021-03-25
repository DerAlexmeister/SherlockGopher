package sherlockanalyser

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	crawler "github.com/DerAlexx/SherlockGopher/sherlockcrawler/sherlockcrawler"

	proto "github.com/DerAlexx/SherlockGopher/analyser/proto"
	mock "github.com/DerAlexx/SherlockGopher/sherlockcrawler/sherlockcrawler/mocks"
	"github.com/golang/mock/gomock"
)

func TestWebsiteData(t *testing.T) {
	sut := NewAnalyserServiceHandler()
	ctrl := gomock.NewController(t)
	mockDataStream := mock.NewMockAnalyser_WebsiteDataStream(ctrl)
	header := http.Header{}
	header.Add("http", "nope")
	task := crawler.NewCrawlerTaskRequest(0,
		"www.helpme.de", 0, fmt.Errorf(""),
		0, nil, &header, "html",
		[]byte("html"), 200, 2000)
	data := PrepareData(task)
	d1 := proto.CrawlerPackage{
		Id:      1,
		Content: data[0],
	}
	d2 := proto.CrawlerPackage{
		Id:      1,
		Content: data[1],
	}

	r1 := mockDataStream.EXPECT().Recv().Return(&d1, nil).Times(1)
	r2 := mockDataStream.EXPECT().Recv().Return(&d2, nil).Times(1)
	r3 := mockDataStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	gomock.InOrder(r1, r2, r3)

	mockDataStream.EXPECT().Send(gomock.Any()).Return(nil).MinTimes(1)
	mockDataStream.EXPECT().Close().MinTimes(1)
	ret := sut.WebsiteData(context.TODO(), mockDataStream)
	if ret != nil {
		t.Fatal("Test failed")
	}

	if len(sut.getQueue().Queue) != 1 {
		t.Fatal("Test failed")
	}
}

func PrepareData(task *crawler.CrawlerTaskRequest) [][]byte {
	buffer := make([][]byte, 0)
	headerPackage := make([]byte, 0)

	// TaskID
	taskIDBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(taskIDBuffer, task.GetTaskID())
	headerPackage = append(headerPackage, taskIDBuffer...)

	// Address
	byteAddress := []byte(task.GetAddr())
	addressSize := uint64(len(byteAddress))
	addressSizeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(addressSizeBuffer, addressSize)
	headerPackage = append(headerPackage, addressSizeBuffer...)
	headerPackage = append(headerPackage, byteAddress...)

	// Status Code
	statusCodeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(statusCodeBuffer, uint64(task.GetStatusCode()))
	headerPackage = append(headerPackage, statusCodeBuffer...)

	// Response Time
	responseTimeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(responseTimeBuffer, uint64(task.GetResponseTime().Nanoseconds()))
	headerPackage = append(headerPackage, responseTimeBuffer...)

	// Error
	var errorSize uint64
	errorSizeBuffer := make([]byte, 8)
	if task.GetTaskError() != nil {
		byteError := []byte(task.GetTaskError().Error())
		errorSize = uint64(len(byteAddress))
		binary.LittleEndian.PutUint64(addressSizeBuffer, errorSize)
		headerPackage = append(headerPackage, errorSizeBuffer...)
		headerPackage = append(headerPackage, byteError...)
	} else {
		binary.LittleEndian.PutUint64(errorSizeBuffer, 0)
		headerPackage = append(headerPackage, errorSizeBuffer...)
	}

	// Header Fields Count
	taskHeaderFieldsCount := make([]byte, 8)
	binary.LittleEndian.PutUint64(taskHeaderFieldsCount, uint64(len(task.GetResponseHeader())))
	headerPackage = append(headerPackage, taskHeaderFieldsCount...)

	// Header Fields
	for k, v := range task.GetResponseHeader() {
		// Header Field Size
		headerValues := strings.Join(v, "J_W")
		headerData := k + "W_J" + headerValues
		headerFieldSize := make([]byte, 8)

		binary.LittleEndian.PutUint64(headerFieldSize, uint64(len(headerData)))
		headerPackage = append(headerPackage, headerFieldSize...)
		headerPackage = append(headerPackage, []byte(headerData)...)
	}

	buffer = append(buffer, headerPackage)
	// Body package count
	buffer = append(buffer, []byte(task.GetResponseBody()))

	return buffer
}
