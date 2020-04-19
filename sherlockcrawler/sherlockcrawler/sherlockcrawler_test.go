package sherlockcrawler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
	protoweb "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertowebserver"
)

const (
	amountoftasks = 100
	staticurl     = "www.google.com"
)

/*
TestSherlockCrawlerNew will create a new instance of NewSherlockCrawlerService
and check its getters and setters.
*/
func TestSherlockCrawlerNew(t *testing.T) {
	sherlock := NewSherlockCrawlerService()
	deps := NewSherlockDependencies()
	streamingserver := NewStreamingServer()

	sherlock.InjectDependency(deps)
	sherlock.SetSherlockStreamer(&streamingserver)
	if queue := sherlock.getQueue(); queue == nil {
		t.Fatal("queue was nil but should not be ")
	} else if deps := sherlock.getDependency(); deps == nil {
		t.Fatal("dependencies was nil but should not be ")
	} else if sts := sherlock.getSherlockStreamer(); sts == nil {
		t.Fatal("streamingserver was nil but should not be ")
	} else {
		t.Log("successfully created all things needed in the new crawler instance")
	}
}

/*
TestCreateTask will test the CreateTask method of Sherlockcrawler
and like the Benachmark test below.
*/
func TestCreateTask(t *testing.T) {
	service := NewSherlockCrawlerService()
	for n := 0; n < amountoftasks; n++ {
		response := proto.CrawlTaskCreateResponse{}
		if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%s", staticurl)}, &response); err != nil {
			t.Fatalf("createtask returned an error %s on number %d", err.Error(), n)
		} else if rescode := response.GetStatuscode(); rescode != proto.URL_STATUS_ok {
			t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_ok, rescode)
		} else {
			t.Logf("Run number %d worked fine.", n)
		}
	}
}

/*
TestCreateTask will test the CreateTask method of Sherlockcrawler
and like the Benachmark test below and fails on purpose.
*/
func TestTryToCreateTaskAndFail(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := proto.CrawlTaskCreateResponse{}
	if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%s", "")}, &response); err == nil {
		t.Fatal("test should fail but got no error")
	} else if rescode := response.GetStatuscode(); rescode == proto.URL_STATUS_ok {
		t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_failure, rescode)
	} else {
		t.Log("Test failed as expected.")
	}
}

/*
TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue will return an error because of a nilpointer for the queue.
*/
func TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue(t *testing.T) {
	service := NewSherlockCrawlerService()
	service.getQueue().Queue = nil
	response := proto.CrawlTaskCreateResponse{}
	if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%s", "")}, &response); err == nil {
		t.Fatalf("test should fail but got no error. Error returned: %s", err.Error())
	} else if rescode := response.GetStatuscode(); rescode == proto.URL_STATUS_ok {
		t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_failure, rescode)
	} else {
		t.Log("Test failed as expected.")
	}
}

/*
TestReceiveURL will test the ReceiveURL method and try to send a valid url.
*/
func TestReceiveURL(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := protoweb.SubmitURLResponse{}
	if err := service.ReceiveURL(context.TODO(), &protoweb.SubmitURLRequest{URL: staticurl}, &response); err != nil {
		t.Fatalf("an error occurred while trying to submit a url. Error: %s", err.Error())
	} else if !response.GetRecieved() {
		t.Fatalf("did not receive the url as expected. Responsefield - Expected: %t, Got, %t - Send error: %s", true, response.GetRecieved(), response.GetError())
	} else {
		t.Log("Url received as expected")
	}
}

/*
TestReceiveURL will test the ReceiveURL method and try to send a valid url but expect it too fail.
*/
func TestReceiveURLIntendedToFail(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := protoweb.SubmitURLResponse{}
	if err := service.ReceiveURL(context.TODO(), &protoweb.SubmitURLRequest{URL: "!"}, &response); err == nil {
		t.Fatalf("an error occurred while trying to submit a url. Error: %s", err.Error())
	} else if response.GetRecieved() {
		t.Fatalf("did not receive the url as expected. Responsefield - Expected: %t, Got, %t - Send error: %s", false, response.GetRecieved(), response.GetError())
	} else {
		t.Log("url was not submited as wanted.")
	}
}

/*
TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue will return an error because of a nilpointer for the queue.
*/
func TestStatusOfQueue(t *testing.T) {
	service := NewSherlockCrawlerService()
	response := protoweb.TaskStatusResponse{}
	var tasks = []*CrawlerTaskRequest{
		&CrawlerTaskRequest{taskstate: UNDONE}, &CrawlerTaskRequest{taskstate: UNDONE}, &CrawlerTaskRequest{taskstate: UNDONE},
		&CrawlerTaskRequest{taskstate: FINISHED}, &CrawlerTaskRequest{taskstate: FINISHED},
		&CrawlerTaskRequest{taskstate: PROCESSING}, &CrawlerTaskRequest{taskstate: PROCESSING},
		&CrawlerTaskRequest{taskstate: PROCESSING}, &CrawlerTaskRequest{taskstate: PROCESSING},
		&CrawlerTaskRequest{taskstate: PROCESSING}, &CrawlerTaskRequest{taskstate: PROCESSING},
		&CrawlerTaskRequest{taskstate: FAILED},
	}
	for i, task := range tasks {
		if id := service.getQueue().AppendQueue(task); id == 0 {
			t.Fatalf("got a zero id for task at index %d", i)
		}
	}
	if service.StatusOfTaskQueue(context.TODO(), nil, &response); response.GetUndone() != 3 {
		t.Fatalf("number of undone tasks does not match. Expected: %d, Got %d", 3, response.GetUndone())
	} else if response.GetFailed() != 1 {
		t.Fatalf("number of failed tasks does not match. Expected: %d, Got %d", 1, response.GetFailed())
	} else if response.GetFinished() != 2 {
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", 2, response.GetFinished())
	} else if response.GetProcessing() != 6 {
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", 6, response.GetProcessing())
	} else {
		t.Log("Got all states as expected.")
	}
}

/*
TestCreateTask will test the CreateTask method as well as RemoveFromQueue of Sherlockcrawler
and like the Benachmark test below.
*/
func TestCreateTasksAndRemoveThem(t *testing.T) {
	service := NewSherlockCrawlerService()
	for n := 0; n < amountoftasks; n++ {
		response := proto.CrawlTaskCreateResponse{}
		if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%s", staticurl)}, &response); err != nil {
			t.Fatalf("createtask returned an error %s on number %d", err.Error(), n)
		} else if rescode := response.GetStatuscode(); rescode != proto.URL_STATUS_ok {
			t.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_ok, rescode)
		} else {
			t.Logf("Run number %d worked fine.", n)
		}
	}

	queue := service.getQueue()
	for _, n := range queue.getAllTaskIds() {
		result := queue.RemoveFromQueue(n)
		if !result && queue.ContainsTaskID(n) {
			t.Fatalf("cannot remove #%d from Queue", n)
		} else {
			t.Logf("successfully removed #%d", n)
		}
	}

	if len := len(*queue.getThisQueue()); len != 0 {
		t.Fatalf("queue is not empty. Excpeted: 0, Got: %d", len)
	} else {
		t.Log("Queue is empty.")
	}

}

/*
TestManageUndoneTasks will test the manageUndoneTasks function.
*/
func TestManageUndoneTasks(t *testing.T) {
	service := NewSherlockCrawlerService()
	var localwaitgroup sync.WaitGroup
	localwaitgroup.Add(1)
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskstate: UNDONE, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}
	go service.manageUndoneTasks(&localwaitgroup)
	localwaitgroup.Wait()
	for k, v := range *(*service.getQueue()).getThisQueue() {
		if v.getTaskError() != nil {
			t.Fatalf("An error occured but shouldnt. Error: %s", v.getTaskError().Error())
		} else if v.getStatusCode() != 200 {
			t.Fatalf("Statuscode of the Task does not match. Expected: %d, Got: %d", 200, v.getStatusCode())
		} else if v.getResponseTime() < 0 {
			t.Fatalf("Got a negativ or to small responsetime: Time: %d", v.getResponseTime())
		} else {
			t.Logf("Worked for task number %d ", k)
		}
	}
}

/*
TestManageFinishedTasks will test the manageFinishedTasks function
by adding undone tasks and later on send the finished tasks.
*/
func TestManageFinishedTasks(t *testing.T) {
	service := NewSherlockCrawlerService()
	deps := NewSherlockDependencies()
	streamingserver := NewStreamingServer()
	service.InjectDependency(deps)
	service.SetSherlockStreamer(&streamingserver)
	var localwaitgroup sync.WaitGroup
	localwaitgroup.Add(1)
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskstate: UNDONE, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}

	go service.manageUndoneTasks(&localwaitgroup)
	localwaitgroup.Wait()
	localwaitgroup.Add(1)
	go service.manageFinishedTasks(&localwaitgroup)
	localwaitgroup.Wait()
	if length := len((*(*service.getSherlockStreamer()).getQueue()).getAllTaskIds()); length != 10 {
		t.Fatalf("The queue should have %d elements but got %d", 10, length)
	} else if lengthOfServiceQueue := len((*service.getQueue()).getAllTaskIds()); lengthOfServiceQueue != 0 {
		t.Fatalf("The queue should have %d elements but got %d", 0, lengthOfServiceQueue)
	} else {
		t.Log("Successfully created all tasks and finished them of.")
	}

}

/*
TestManageFailedTasks will test the manageFailedTasks function
by adding 10 failed tasks with 3 trys and 10 failed with 1 try so
in the end the streaming server should have 10 task and the crawler
should have 10 tasks.
*/
func TestManageFailedTasks(t *testing.T) {
	service := NewSherlockCrawlerService()
	deps := NewSherlockDependencies()
	streamingserver := NewStreamingServer()
	service.InjectDependency(deps)
	service.SetSherlockStreamer(&streamingserver)
	var localwaitgroup sync.WaitGroup
	localwaitgroup.Add(1)
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskstate: FAILED, taskerrortry: 1, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskstate: FAILED, taskerrortry: 3, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}

	go service.manageFailedTasks(&localwaitgroup)
	localwaitgroup.Wait()
	if length := len((*(*service.getSherlockStreamer()).getQueue()).getAllTaskIds()); length != 10 {
		t.Fatalf("The queue should have %d elements but got %d", 10, length)
	} else if lengthOfServiceQueue := len((*service.getQueue()).getAllTaskIds()); lengthOfServiceQueue != 10 {
		t.Fatalf("The queue should have %d elements but got %d", 10, lengthOfServiceQueue)
	} else {
		t.Log("Successfully created all tasks and finished them of.")
	}

}

/*
TestRunManager will test the runManager function which will run over all
types of tasks an start a goroutine for the function ManageFailedTasks, manageFinishedTasks
and manageUndoneTasks. In the end there should be 10 failed task for the crawler and 10 finished
task which are previously undone tasks and 10 failed task to send to the analyser.
*/
func TestRunManager(t *testing.T) {
	service := NewSherlockCrawlerService()
	deps := NewSherlockDependencies()
	streamingserver := NewStreamingServer()
	service.InjectDependency(deps)
	service.SetSherlockStreamer(&streamingserver)
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskstate: FAILED, taskerrortry: 1, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskstate: FAILED, taskerrortry: 3, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}
	for i := 0; i < 10; i++ {
		if id := service.getQueue().AppendQueue(&CrawlerTaskRequest{taskstate: UNDONE, addr: server.URL}); id == 0 {
			t.Fatalf("got a zero id for task at index %d ", i)
		}
	}

	service.runManager()

	if length := len((*(*service.getSherlockStreamer()).getQueue()).getAllTaskIds()); length != 10 {
		t.Fatalf("The queue should have %d elements but got %d", 10, length)
	} else if lengthOfServiceQueue := len((*service.getQueue()).getAllTaskIds()); lengthOfServiceQueue != 20 {
		t.Fatalf("The queue should have %d elements but got %d", 20, lengthOfServiceQueue)
	} else {
		t.Log("Successfully created all tasks and finished them of.")
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
		if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%s", staticurl)}, &response); err != nil {
			b.Fatalf("createtask returned an error %s on number %d", err.Error(), n)
		} else if rescode := response.GetStatuscode(); rescode != proto.URL_STATUS_ok {
			b.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_ok, rescode)
		} else {
			b.Logf("Run number %d worked fine.", n)
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
		if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%s", staticurl)}, &response); err != nil {
			b.Fatalf("createtask returned an error %s on number %d", err.Error(), n)
		} else if rescode := response.GetStatuscode(); rescode != proto.URL_STATUS_ok {
			b.Fatalf("got unexpected status. Expected: %s, Got: %s", proto.URL_STATUS_ok, rescode)
		} else {
			b.Logf("Run number %d worked fine.", n)
		}
	}

	queue := service.getQueue()
	for _, n := range queue.getAllTaskIds() {
		result := queue.RemoveFromQueue(n)
		if !result && queue.ContainsTaskID(n) {
			b.Fatalf("cannot remove #%d from Queue", n)
		} else {
			b.Logf("successfully removed #%d", n)
		}
	}

	if len := len(*queue.getThisQueue()); len != 0 {
		b.Fatalf("queue is not empty. Excpeted: 0, Got: %d", len)
	} else {
		b.Log("Queue is empty.")
	}

}
