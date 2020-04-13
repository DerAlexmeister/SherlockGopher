package sherlockcrawler

import (
	"context"
	"fmt"
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
