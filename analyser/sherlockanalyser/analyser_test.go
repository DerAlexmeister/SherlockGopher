package sherlockanalyser

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"net/http"
	"strconv"
	"testing"
	"time"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto"
)

func TestAnalyserGetterSetter(t *testing.T) {
	analyser := NewAnalyserServiceHandler()
	deps := NewAnalyserDependencies()
	analyser.InjectDependency(deps)

	if analyser.Watchdog() == nil {
		t.Fatal("Watchdog was not set")
	}

	dur := time.Second
	analyser.SetDelay(&dur)

	if *analyser.Delay() != dur {
		t.Fatal("Duration is wrong")
	}
}

func TestManageTasks(t *testing.T) {
	analyser := NewAnalyserServiceHandler()
	deps := NewAnalyserDependencies()
	analyser.InjectDependency(deps)

	header := http.Header{}
	cData := CrawlerData{
		taskID:            1,
		addr:              "www.error.err",
		taskError:         fmt.Errorf("test"),
		responseHeader:    &header,
		responseBodyBytes: []byte(""),
		statusCode:        200,
		responseTime:      0,
	}


	mockCtrl := gomock.NewController(t)
	mockNeoSaver := NewMockneoSaverInterface(mockCtrl)
	mockNeoSaver.EXPECT().Save(gomock.Any()).MinTimes(1)
	mockNeoSaver.EXPECT().GetSession().Return(nil)
	mockNeoSaver.EXPECT().Contains(gomock.Any()).Return(make([]string, 0)).MinTimes(1)

	task := injectDependencies(NewTask(&cData))
	(*task).SetSaver(mockNeoSaver)

	task2 := injectDependencies(NewTask(&cData))
	task2.SetState(FINISHED)

	analyser.getQueue().AppendQueue(task)
	analyser.getQueue().AppendQueue(task2)

	go analyser.ManageTasks()

	time.Sleep(3 * time.Second)

	in := proto.ChangeStateRequest{
		State: proto.AnalyserStateEnum_Stop,
	}
	out := proto.ChangeStateResponse{}

	if analyser.ChangeStateRPC(context.TODO(), &in, &out) != nil {
		t.Fatal("error while setting state")
	}

	if !analyser.getQueue().IsEmpty() {
		t.Fatal("Stop didnt work")
	}

	if analyser.State() != STOP {
		t.Fatal("RunManager failed")
	}
}

func GetExampleCData() CrawlerData {
	header := http.Header{}
	cData := CrawlerData{
		taskID:            1,
		addr:              "www.error.err",
		taskError:         fmt.Errorf("test"),
		responseHeader:    &header,
		responseBodyBytes: []byte(""),
		statusCode:        200,
		responseTime:      0,
	}
	return cData
}

/*
TestSherlockCrawlerNew will create a new instance of NewAnalyserServiceHandler
and check its getters and setters.
*/
func TestNewAnalyserServiceHandler(t *testing.T) {
	sherlock := NewAnalyserServiceHandler()
	deps := NewAnalyserDependencies()
	sherlock.InjectDependency(deps)

	if queue := sherlock.getQueue(); queue == nil {
		t.Fatal("queue was nil but should not be ")
		//} else if deps := sherlock.getDependency(); deps == nil { TODO: Alex fragen
		//	t.Fatal("dependencies was nil but should not be ")
	} else {
		fmt.Println("successfully created all things needed in the new analyser instance")
	}
}

/*
TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue will return an error because of a nilpointer for the queue.
*/
func TestWorkloadRPC(t *testing.T) {
	service := NewAnalyserServiceHandler()
	response := proto.WorkloadResponse{}
	request := proto.WorkloadRequest{}

	var undone uint64 = 1
	var processing uint64 = 2
	var crawlerError uint64 = 1
	var saving uint64 = 2
	var sendToCrawler uint64 = 1
	var finished uint64 = 2

	cData := GetExampleCData()
	tasks := make([]*AnalyserTaskRequest, 0)
	i := 0
	cData.setAddr(cData.getAddr() + strconv.Itoa(i))
	t1 := AnalyserTaskRequest{state: UNDONE, crawlerData: &cData}
	tasks = append(tasks, &t1)
	i++
	cData2 := GetExampleCData()
	cData2.setAddr(cData.getAddr() + strconv.Itoa(i))
	t2 := AnalyserTaskRequest{state: PROCESSING, crawlerData: &cData2}
	tasks = append(tasks, &t2)
	i++
	cData3 := GetExampleCData()
	cData3.setAddr(cData.getAddr() + strconv.Itoa(i))
	t3 := AnalyserTaskRequest{state: PROCESSING, crawlerData: &cData3}
	tasks = append(tasks, &t3)
	i++
	cData4 := GetExampleCData()
	cData4.setAddr(cData.getAddr() + strconv.Itoa(i))
	t4 := AnalyserTaskRequest{state: CRAWLERERROR, crawlerData: &cData4}
	tasks = append(tasks, &t4)
	i++
	cData5 := GetExampleCData()
	cData5.setAddr(cData.getAddr() + strconv.Itoa(i))
	t5 := AnalyserTaskRequest{state: SAVING, crawlerData: &cData5}
	tasks = append(tasks, &t5)
	i++
	cData6 := GetExampleCData()
	cData6.setAddr(cData.getAddr() + strconv.Itoa(i))
	t6 := AnalyserTaskRequest{state: SAVING, crawlerData: &cData6}
	tasks = append(tasks, &t6)
	i++
	cData7 := GetExampleCData()
	cData7.setAddr(cData.getAddr() + strconv.Itoa(i))
	t7 := AnalyserTaskRequest{state: SENDTOCRAWLER, crawlerData: &cData7}
	tasks = append(tasks, &t7)
	i++
	cData8 := GetExampleCData()
	cData8.setAddr(cData.getAddr() + strconv.Itoa(i))
	t8 := AnalyserTaskRequest{state: FINISHED, crawlerData: &cData8}
	tasks = append(tasks, &t8)
	i++
	cData9 := GetExampleCData()
	cData9.setAddr(cData.getAddr() + strconv.Itoa(i))
	t9 := AnalyserTaskRequest{state: FINISHED, crawlerData: &cData9}
	tasks = append(tasks, &t9)

	for i, task := range tasks {
		if id := service.getQueue().AppendQueue(task); id == 0 {
			t.Fatalf("got a zero id for task at index %d", i)
		}
	}

	err := service.WorkloadRPC(context.TODO(), &request, &response)
	switch {
	case err != nil:
		t.Fatalf("error occurred: %s", err)
	case response.GetUndone() != undone:
		t.Fatalf("number of undone tasks does not match. Expected: %d, Got %d", undone, response.GetUndone())
	case response.GetProcessing() != processing:
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", processing, response.GetProcessing())
	case response.GetCrawlerError() != crawlerError:
		t.Fatalf("number of crawlerError tasks does not match. Expected: %d, Got %d", crawlerError, response.GetCrawlerError())
	case response.GetSaving() != saving:
		t.Fatalf("number of savingtasks does not match. Expected: %d, Got %d", saving, response.GetSaving())
	case response.GetSendToCrawler() != sendToCrawler:
		t.Fatalf("number of sendToCrawler tasks does not match. Expected: %d, Got %d", sendToCrawler, response.GetSendToCrawler())
	case response.GetFinished() != 2:
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", finished, response.GetFinished())
	default:
		fmt.Println("Got all states as expected.")
	}
}

/*
TestSherlockCrawler_GetState tests GetState
*/
func TestSherlockCrawler_GetState(t *testing.T) {
	service := NewAnalyserServiceHandler()
	out := proto.StateResponse{}
	in := proto.StateRequest{}

	service.Observer().UpdateState(RUNNING)
	_ = service.StateRPC(context.TODO(), &in, &out)

	if out.State != RUNNING {
		t.Fatalf("expected running but was %v", out.State)
	}
}

/*
TestSherlockCrawler_SetState_Stop tests SetState
*/
func TestSherlockCrawler_SetState_Stop(t *testing.T) {
	service := NewAnalyserServiceHandler()
	in := proto.ChangeStateRequest{
		State: proto.AnalyserStateEnum_Stop,
	}
	out := proto.ChangeStateResponse{}

	if service.ChangeStateRPC(context.TODO(), &in, &out) != nil {
		t.Fatal("error while setting state")
	}

	if service.State() != STOP {
		t.Fatalf("expected running but was %v", service.State())
	}
}

/*
TestSherlockCrawler_SetState_Pause sets SetState
*/
func TestSherlockCrawler_SetState_Pause(t *testing.T) {
	service := NewAnalyserServiceHandler()
	in := proto.ChangeStateRequest{
		State: proto.AnalyserStateEnum_Pause,
	}
	out := proto.ChangeStateResponse{}

	if service.ChangeStateRPC(context.TODO(), &in, &out) != nil {
		t.Fatal("error while setting state")
	}

	if service.State() != PAUSE {
		t.Fatalf("expected running but was %v", service.State())
	}
}

/*
TestSherlockCrawler_SetState_Running tests SetState
*/
func TestSherlockCrawler_SetState_Running(t *testing.T) {
	service := NewAnalyserServiceHandler()
	in := proto.ChangeStateRequest{
		State: proto.AnalyserStateEnum_Running,
	}
	out := proto.ChangeStateResponse{}

	if service.ChangeStateRPC(context.TODO(), &in, &out) != nil {
		t.Fatal("error while setting state")
	}

	if service.State() != RUNNING {
		t.Fatalf("expected running but was %v", service.State())
	}
}

func TestSherlockCrawler_SetState_Idle(t *testing.T) {
	service := NewAnalyserServiceHandler()
	in := proto.ChangeStateRequest{
		State: proto.AnalyserStateEnum_Idle,
	}
	out := proto.ChangeStateResponse{}

	if service.ChangeStateRPC(context.TODO(), &in, &out) == nil {
		t.Fatal("error while setting state")
	}
}

/*
TestSherlockCrawler_SetState_Clean tests SetState
*/
func TestSherlockCrawler_SetState_Clean(t *testing.T) {
	service := NewAnalyserServiceHandler()
	task := NewTask(&CrawlerData{
		taskID:            1,
		addr:              "123",
		taskError:         nil,
		responseHeader:    &http.Header{},
		responseBodyBytes: []byte("123"),
		statusCode:        0,
		responseTime:      0,
	})
	service.getQueue().AppendQueue(task)

	in := proto.ChangeStateRequest{
		State: proto.AnalyserStateEnum_Clean,
	}
	out := proto.ChangeStateResponse{}

	if service.ChangeStateRPC(context.TODO(), &in, &out) != nil {
		t.Fatal("error while setting state")
	}

	if service.State() != RUNNING {
		t.Fatalf("expected running but was %v", service.State())
	}

	if len(service.getQueue().getAllTaskIds()) != 0 {
		t.Fatalf("expected 0 but was %v", len(service.getQueue().getAllTaskIds()))
	}
}

/*
TestSherlockCrawler_SetState_Default tests SetState
*/
func TestSherlockCrawler_SetState_Default(t *testing.T) {
	service := NewAnalyserServiceHandler()
	in := proto.ChangeStateRequest{
		State: 123,
	}
	out := proto.ChangeStateResponse{}

	if service.ChangeStateRPC(context.TODO(), &in, &out) == nil {
		t.Fatal("expected error while setting state but was successful")
	}

	if out.GetSuccess() != false {
		t.Fatalf("expected false but was %v", out.GetSuccess())
	}
}

/*
TestSherlockCrawler_AnalyserState tests AnalyserState
*/
func TestSherlockCrawler_AnalyserState(t *testing.T) {
	service := NewAnalyserServiceHandler()

	if service.State() != RUNNING {
		t.Fatalf("expected running but was %v", service.State())
	}
}

/*
TestSherlockCrawler_SetAnalyserState tests SetAnalyserState
*/
func TestSherlockCrawler_SetAnalyserState(t *testing.T) {
	service := NewAnalyserServiceHandler()
	service.Observer().UpdateState(PAUSE)

	if service.getQueue().State() != PAUSE {
		t.Fatalf("expected pause but was %v", service.getQueue().State())
	}

	if service.State() != PAUSE {
		t.Fatalf("expected pause but was %v", service.State())
	}
}
