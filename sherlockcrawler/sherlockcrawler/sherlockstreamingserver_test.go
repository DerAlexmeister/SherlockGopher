package sherlockcrawler

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/micro/go-micro"
	"github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockanalyser"
	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyserfiletransfer"
	"github.com/pkg/errors"
)

const (
	testid                       = 1
	testaddr                     = "localhost:8080"
	testtime       time.Duration = 10
	teststatuscode               = 200
)

// Testing error
var erro error = errors.New("testerror")

/*
TestMin compares two values and should return the smaller one
*/
func TestMin(t *testing.T) {
	a := 1
	b := 2
	ret := min(a, b)
	if ret != a {
		t.Fatalf("Min function: expected %d but was %d", a, ret)
	} else {
		t.Log("Min fuction works")
	}
}

/*
TestGetChunkSize should return the chunksize of 1024
*/
func TestGetChunkSize(t *testing.T) {
	ret := getChunkSize()
	if ret != 1024 {
		t.Fatalf("Min function: expected %d but was %d", 1024, ret)
	} else {
		t.Log("Min fuction works")
	}
}

/*
TestNewServer should return a new StreamingServer
*/
func TestNewServer(t *testing.T) {
	server := NewStreamingServer()
	if server.getQueue() == nil {
		t.Fatal("Creating new Streaming Server failed")
	} else {
		t.Log("NewStreamingServer works")
	}
}

/*
TestSendAndReceive tests the Upload function of the receivefromcrawler file and
helpInfos, helpSend, helpSendFileToAnalyser and sendFileToAnalyser function of the sherlockstreamingserver file.
Tests the grpc stream functionality between the crawler and analyser
*/
func TestSendAndReceive(t *testing.T) {
	ltask := createTask()

	streamServer := NewStreamingServer()
	streamServer.getQueue().AppendQueue(&ltask)
	cservice := micro.NewService()
	cservice.Init()

	aservice := micro.NewService(
		micro.Name("go.micro.srv.stream"),
	)
	aservice.Init()
	createAQ := sherlockanalyser.NewAnalyserQueue()
	serverGrpc := sherlockanalyser.NewServerGRPC(&createAQ)
	proto.RegisterSenderHandler(aservice.Server(), serverGrpc)
	if err := aservice.Run(); err != nil {
		log.Fatal(err)
	}

	// create client
	cl := proto.NewSenderService("go.micro.srv.stream", cservice.Client())
	err := streamServer.sendFileToAnalyser(cl)

	time.Sleep(200 * time.Millisecond)
	aTestQ := serverGrpc.Queue
	if err != nil || aTestQ.IsEmpty() {
		t.Fatal("Streaming failed")
	} else {
		t.Log("Streaming works")
	}
}

/*
TestSendAndReceiveErrorCase tests the Upload function of the receivefromcrawler file and
helpErrorCase, helpSendFileToAnalyser and sendFileToAnalyser function of the sherlockstreamingserver file.
Tests the grpc stream functionality between the crawler and analyser
*/
func TestSendAndReceiveErrorCase(t *testing.T) {
	ltask := createTaskErrorCase()

	streamServer := NewStreamingServer()
	streamServer.getQueue().AppendQueue(&ltask)
	cservice := micro.NewService()
	cservice.Init()

	aservice := micro.NewService(
		micro.Name("go.micro.srv.stream"),
	)
	aservice.Init()
	createAQ := sherlockanalyser.NewAnalyserQueue()
	serverGrpc := sherlockanalyser.NewServerGRPC(&createAQ)
	proto.RegisterSenderHandler(aservice.Server(), serverGrpc)
	if err := aservice.Run(); err != nil {
		log.Fatal(err)
	}

	// create client
	cl := proto.NewSenderService("go.micro.srv.stream", cservice.Client())
	err := streamServer.sendFileToAnalyser(cl)

	time.Sleep(200 * time.Millisecond)
	aTestQ := serverGrpc.Queue
	if err != nil || aTestQ.IsEmpty() {
		t.Fatal("Streaming ErrorCase failed")
	} else {
		t.Log("Streaming ErrorCase works")
	}
}

/*
createTask() creates a dummy task for the TestSendAndReceive function
*/
func createTask() CrawlerTaskRequest {
	ltask := NewTask()

	header := http.Header{}
	header.Add("Accept-Ranges", "bytes")
	header.Add("Age", "12")
	header.Add("Allow", "GET, HEAD")
	header.Add("Cache-Control", "max-age=3600")
	header.Add("Connection", "close")
	header.Add("Content-Encoding", "gzip")
	header.Add("Content-Language", "de")
	header.Add("Content-Location", "/foo.html.de")
	header.Add("Content-MD5", "Q2hlY2sgSW50ZWdyaXR5IQ==")
	header.Add("Content-Type", "text/html; charset=utf-8")
	header.Add("Expires", "Thu, 01 Dec 1994 16:00:00 GMT")
	header.Add("Proxy-Authenticate", "Basic")

	bytearr := []byte{}
	for i := 0; i < 5000; i++ {
		bytearr = append(bytearr, byte(i))
	}

	ltask.setTaskID(uint64(10))
	ltask.setAddr("127.0.0.1")
	ltask.setResponseHeader(&header)
	ltask.setResponseBodyInBytes(bytearr)
	ltask.setStatusCode(200)
	ltask.setResponseTime(time.Duration(5))

	return ltask
}

/*
createTaskErrorCase() creates a dummy task for the TestSendAndReceiveErrorCase function
*/
func createTaskErrorCase() CrawlerTaskRequest {
	ltask := NewTask()

	ltask.setTaskID(uint64(10))
	ltask.setAddr("127.0.0.1")
	ltask.setTaskError(errors.New("test"))
	ltask.setResponseTime(time.Duration(5))

	return ltask
}
