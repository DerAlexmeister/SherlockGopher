package sherlockcrawler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto"
)

const (
	amountoftasks        = 100
	staticurl            = "www.google.com"
	wrongHTTPBody string = "This should be in the body of the HTTP response."
)

/*
TestCreateTask will test the CreateTask method of SherlockCrawler
and like the Benachmark test below.
*/
func TestCreateTask(t *testing.T) {
	service := NewSherlockCrawlerService()
	for n := 0; n < amountoftasks; n++ {
		response := proto.CrawlTaskCreateResponse{}
		err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: staticurl}, &response)
		resCode := response.GetStatuscode()
		switch {
		case err != nil:
			t.Fatalf("createtask returned an error %s on number %d", err.Error(), n)
		case resCode != proto.URL_STATUS_ok:
			t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_ok, resCode)
		default:
			continue
		}
	}

	taskCount := len(*service.getQueue().GetThisQueue())
	if taskCount != 100 {
		t.Fatalf("expected 100 but has %v tasks", taskCount)
	} else {
		fmt.Println("successfully created tasks")
	}
}

/*
TestCreateTask will test the CreateTask method of SherlockCrawler
and like the Benachmark test below and fails on purpose.
*/
func TestTryToCreateTaskAndFail(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := proto.CrawlTaskCreateResponse{}
	err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: ""}, &response)
	resCode := response.GetStatuscode()
	switch {
	case err == nil:
		t.Fatal("test should fail but got no error")
	case resCode == proto.URL_STATUS_ok:
		t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_failure, resCode)
	default:
		fmt.Println("Test failed as expected.")
	}
}

/*
TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue will return an error because of a nilpointer for the queue.
*/
func TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue(t *testing.T) {
	service := NewSherlockCrawlerService()
	service.getQueue().Queue = nil
	response := proto.CrawlTaskCreateResponse{}
	err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: ""}, &response)
	resCode := response.GetStatuscode()

	switch {
	case err == nil:
		t.Fatalf("test should fail but got no error. Error returned: %s", err.Error())
	case resCode == proto.URL_STATUS_ok:
		t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_failure, resCode)
	default:
		fmt.Println("Test failed as expected.")
	}
}

/*
TestReceiveURL will test the ReceiveURL method and try to send a valid url.
*/
func TestReceiveURL(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := proto.SubmitURLResponse{}
	err := service.ReceiveURL(context.TODO(), &proto.SubmitURLRequest{URL: staticurl}, &response)

	switch {
	case err != nil:
		t.Fatalf("an error occurred while trying to submit a url. Error: %s", err.Error())
	case !response.GetRecieved():
		t.Fatalf("did not receive the url as expected. Responsefield - Expected: %t, Got, %t - Send error: %s", true, response.GetRecieved(), response.GetError())
	default:
		fmt.Println("Url received as expected")
	}
}

/*
TestReceiveURL will test the ReceiveURL method and try to send a valid url but expect it too fail.
*/
func TestReceiveURLIntendedToFail(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := proto.SubmitURLResponse{}
	err := service.ReceiveURL(context.TODO(), &proto.SubmitURLRequest{URL: "!"}, &response)
	switch {
	case err == nil:
		t.Fatalf("an error occurred while trying to submit a url. Error: %s", err.Error())
	case response.GetRecieved():
		t.Fatalf("did not receive the url as expected. Responsefield - Expected: %t, Got, %t - Send error: %s", false, response.GetRecieved(), response.GetError())
	default:
		fmt.Println("url was not submited as wanted.")
	}
}

/*
TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue will return an error because of a nil pointer for the queue.
*/
func TestStatusOfQueue(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := proto.TaskStatusResponse{}
	var tasks = []*CrawlerTaskRequest{
		{taskState: UNDONE}, {taskState: UNDONE}, {taskState: UNDONE},
		{taskState: FINISHED}, {taskState: FINISHED},
		{taskState: PROCESSING}, {taskState: PROCESSING},
		{taskState: PROCESSING}, {taskState: PROCESSING},
		{taskState: PROCESSING}, {taskState: PROCESSING},
		{taskState: FAILED},
	}
	for i, task := range tasks {
		if id := service.getQueue().AppendQueue(task); id == 0 {
			t.Fatalf("got a zero id for task at index %d", i)
		}
	}
	err := service.StatusOfTaskQueue(context.TODO(), nil, &response)
	switch {
	case err != nil && response.GetUndone() != 3:
		t.Fatalf("number of undone tasks does not match. Expected: %d, Got %d", 3, response.GetUndone())
	case response.GetFailed() != 1:
		t.Fatalf("number of failed tasks does not match. Expected: %d, Got %d", 1, response.GetFailed())
	case response.GetFinished() != 2:
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", 2, response.GetFinished())
	case response.GetProcessing() != 6:
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", 6, response.GetProcessing())
	default:
		fmt.Println("Got all states as expected.")
	}
}

/*
TestCreateTask will test the CreateTask method as well as RemoveFromQueue of SherlockCrawler
and like the Benachmark test below.
*/
func TestCreateTasksAndRemoveThem(t *testing.T) {
	service := NewSherlockCrawlerService()
	for n := 0; n < amountoftasks; n++ {
		response := proto.CrawlTaskCreateResponse{}
		err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: staticurl}, &response)
		resCode := response.GetStatuscode()
		switch {
		case err != nil:
			t.Fatalf("createtask returned an error %s on number %d", err.Error(), n)
		case resCode != proto.URL_STATUS_ok:
			t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_ok, resCode)
		default:
			continue
		}
	}

	queue := service.getQueue()
	for _, n := range queue.getAllTaskIds() {
		result := queue.RemoveFromQueue(n)
		if !result && queue.ContainsTaskID(n) {
			t.Fatalf("cannot remove #%d from Queue", n)
		} else {
			continue
		}
	}

	if size := len(*queue.GetThisQueue()); size != 0 {
		t.Fatalf("queue is not empty. Excpeted: 0, Got: %d", size)
	} else {
		fmt.Println("Queue is empty.")
	}

}

/*
TestManageUndoneTasks will test the manageUndoneTasks function.
*/
func TestManageUndoneTasks(t *testing.T) {
	service := NewSherlockCrawlerService()
	var localwaitgroup sync.WaitGroup
	localwaitgroup.Add(1)
	wanted := wrongHTTPBody
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(wanted))
		if err != nil {
			t.Fatal(err)
		}
	}))
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskState: UNDONE, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}
	go service.manageUndoneTasks(&localwaitgroup)
	localwaitgroup.Wait()
	for k, v := range *service.getQueue().GetThisQueue() {
		switch {
		case v.GetTaskError() != nil:
			t.Fatalf("An error occurred but shouldnt. Error: %s", v.GetTaskError().Error())
		case v.GetStatusCode() != 200:
			t.Fatalf("Statuscode of the Task does not match. Expected: %d, Got: %d", 200, v.GetStatusCode())
		case v.GetResponseTime() < 0:
			t.Fatalf("Got a negativ or to small responsetime: Time: %d", v.GetResponseTime())
		default:
			fmt.Printf("Worked for task number %d ", k)
		}
	}

	if size := len(*service.getQueue().GetThisQueue()); size != 10 {
		t.Fatalf("queue is not empty. Excpeted: 10, Got: %d", size)
	} else {
		fmt.Println("Queue is empty.")
	}

	for _, task := range *service.getQueue().GetThisQueue() {
		if task.GetTaskState() != FINISHED {
			t.Fatalf("expected FINISHED but was %v", task.GetTaskState())
		}
	}
}

/*
BenchmarkCreateTask benchmark test for the createtask function
to see how many task can be created in a secound.
*/
func BenchmarkCreateTask(b *testing.B) {
	service := NewSherlockCrawlerService()
	for n := 0; n < b.N; n++ {
		response := proto.CrawlTaskCreateResponse{}
		err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: staticurl}, &response)
		resCode := response.GetStatuscode()
		switch {
		case err != nil:
			b.Fatalf("createtask returned an error %s on number %d", err.Error(), n)
		case resCode != proto.URL_STATUS_ok:
			b.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_ok, resCode)
		default:
			fmt.Printf("Run number %d worked fine.", n)
		}
	}
}

/*
BenchmarkCreateTask benchmark test for the createtask and RemoveFromQueue functions
to see how many task can be created and deleted in a secound.
*/
func BenchmarkCreateTaskAndClearQueue(b *testing.B) {
	service := NewSherlockCrawlerService()
	for n := 0; n < b.N; n++ {
		response := proto.CrawlTaskCreateResponse{}
		err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: staticurl}, &response)
		resCode := response.GetStatuscode()
		switch {
		case err != nil:
			b.Fatalf("createtask returned an error %s on number %d \n", err.Error(), n)
		case resCode != proto.URL_STATUS_ok:
			b.Fatalf("got unexpected status. Expected: %s, Got: %s \n", proto.URL_STATUS_ok, resCode)
		default:
			fmt.Printf("Run number %d worked fine. \n", n)
		}
	}

	queue := service.getQueue()
	for _, n := range queue.getAllTaskIds() {
		result := queue.RemoveFromQueue(n)
		if !result && queue.ContainsTaskID(n) {
			b.Fatalf("cannot remove #%d from Queue", n)
		} else {
			fmt.Printf("successfully removed #%d", n)
		}
	}

	if size := len(*queue.GetThisQueue()); size != 0 {
		b.Fatalf("queue is not empty. Excpeted: 0, Got: %d", size)
	} else {
		fmt.Println("Queue is empty.")
	}

}

func TestSherlockCrawler_GetState(t *testing.T) {
	service := NewSherlockCrawlerService()
	out := proto.StateGetResponse{}

	service.SetCrawlerState(RUNNING)
	_ = service.GetState(context.TODO(), nil, &out)

	if out.State != RUNNING {
		t.Fatalf("expected running but was %v", out.State)
	}
}

func TestSherlockCrawler_SetState_Stop(t *testing.T) {
	service := NewSherlockCrawlerService()
	in := proto.StateRequest{
		State: proto.CurrentState_Stop,
	}
	out := proto.StateResponse{}

	if service.SetState(context.TODO(), &in, &out) != nil {
		t.Fatal("setting state failed")
	}

	if service.CrawlerState() != STOP {
		t.Fatalf("expected running but was %v", service.CrawlerState())
	}
}

func TestSherlockCrawler_SetState_Pause(t *testing.T) {
	service := NewSherlockCrawlerService()
	in := proto.StateRequest{
		State: proto.CurrentState_Pause,
	}
	out := proto.StateResponse{}

	if service.SetState(context.TODO(), &in, &out) != nil {
		t.Fatal("setting state failed")
	}

	if service.CrawlerState() != PAUSE {
		t.Fatalf("expected running but was %v", service.CrawlerState())
	}
}

func TestSherlockCrawler_SetState_Running(t *testing.T) {
	service := NewSherlockCrawlerService()
	in := proto.StateRequest{
		State: proto.CurrentState_Running,
	}
	out := proto.StateResponse{}

	if service.SetState(context.TODO(), &in, &out) != nil {
		t.Fatal("setting state failed")
	}

	if service.CrawlerState() != RUNNING {
		t.Fatalf("expected running but was %v", service.CrawlerState())
	}
}

func TestSherlockCrawler_SetState_Clean(t *testing.T) {
	service := NewSherlockCrawlerService()
	task := NewTask()
	service.getQueue().AppendQueue(&task)

	in := proto.StateRequest{
		State: proto.CurrentState_Clean,
	}
	out := proto.StateResponse{}

	if service.SetState(context.TODO(), &in, &out) != nil {
		t.Fatal("setting state failed")
	}

	if service.CrawlerState() != RUNNING {
		t.Fatalf("expected running but was %v", service.CrawlerState())
	}

	if len(service.Queue.getAllTaskIds()) != 0 {
		t.Fatalf("expected 0 but was %v", len(service.Queue.getAllTaskIds()))
	}
}

func TestSherlockCrawler_SetState_Default(t *testing.T) {
	service := NewSherlockCrawlerService()
	in := proto.StateRequest{
		State: 124,
	}

	out := proto.StateResponse{}

	if service.SetState(context.TODO(), &in, &out) == nil {
		t.Fatal("setting state failed")
	}

	if out.Received != false {
		t.Fatalf("expected false but was %v", out.Received)
	}
}

func TestSherlockCrawler_CrawlerState(t *testing.T) {
	service := NewSherlockCrawlerService()

	if service.CrawlerState() != RUNNING {
		t.Fatalf("expected running but was %v", service.CrawlerState())
	}
}

func TestSherlockCrawler_SetCrawlerState(t *testing.T) {
	service := NewSherlockCrawlerService()
	fmt.Println(service.Observer())
	fmt.Println(service.getQueue())
	fmt.Println(service.getStreamQueue())

	service.Observer().UpdateState(PAUSE)

	if service.getQueue().QueueState() != PAUSE {
		t.Fatalf("expected pause but was %v", service.getQueue().QueueState())
	}

	if service.CrawlerState() != PAUSE {
		t.Fatalf("expected pause but was %v", service.CrawlerState())
	}
}

func TestGetterAndSetter(t *testing.T) {
	service := NewSherlockCrawlerService()
	var x time.Duration = 5
	service.SetDelay(&x)
	if service.Delay() != &x {
		t.Fatal("Delay Getter and Setter dont work")
	}
	if service.Watchdog() == nil {
		t.Fatal("Watchdog Getter and Setter dont work")
	}
	service.Dependencies = NewSherlockDependencies()
	if service.getDependency() == nil {
		t.Fatal("Dependencies dont work")
	}
}

func TestManageFinishedTasks(t *testing.T){
	sut := NewSherlockCrawlerService()
	wg := sync.WaitGroup{}
	wg.Add(1)
	sut.manageFinishedTasks(&wg)
}


