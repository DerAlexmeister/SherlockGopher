package sherlockanalyser

import (
	"context"
	"testing"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto"
)

/*
TestSherlockCrawlerNew will create a new instance of NewSherlockCrawlerService
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
		t.Log("successfully created all things needed in the new analyser instance")
	}
}

/*
TestTryToCreateTaskAndFailOnCreatingTaskBecauseOfNilQueue will return an error because of a nilpointer for the queue.
*/
func TestStatusOfQueue(t *testing.T) {
	service := NewAnalyserServiceHandler()
	response := proto.TaskStatusResponse{}

	var undone uint64 = 1
	var processing uint64 = 2
	var crawlerError uint64 = 1
	var saving uint64 = 2
	var sendToCrawler uint64 = 1
	var finished uint64 = 2

	var tasks = []*analyserTaskRequest{
		&analyserTaskRequest{state: UNDONE},
		&analyserTaskRequest{state: PROCESSING}, &analyserTaskRequest{state: PROCESSING},
		&analyserTaskRequest{state: CRAWLERERROR},
		&analyserTaskRequest{state: SAVING}, &analyserTaskRequest{state: SAVING},
		&analyserTaskRequest{state: SENDTOCRAWLER},
		&analyserTaskRequest{state: FINISHED}, &analyserTaskRequest{state: FINISHED},
	}

	for i, task := range tasks {
		if id := service.getQueue().AppendQueue(task); id == 0 {
			t.Fatalf("got a zero id for task at index %d", i)
		}
	}
	if service.StatusOfTaskQueue(context.TODO(), nil, &response); response.GetUndone() != undone {
		t.Fatalf("number of undone tasks does not match. Expected: %d, Got %d", undone, response.GetUndone())
	} else if response.GetProcessing() != processing {
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", processing, response.GetProcessing())
	} else if response.GetCrawlerError() != crawlerError {
		t.Fatalf("number of crawlerError tasks does not match. Expected: %d, Got %d", crawlerError, response.GetCrawlerError())
	} else if response.GetSaving() != saving {
		t.Fatalf("number of savingtasks does not match. Expected: %d, Got %d", saving, response.GetSaving())
	} else if response.GetSendToCrawler() != sendToCrawler {
		t.Fatalf("number of sendToCrawler tasks does not match. Expected: %d, Got %d", sendToCrawler, response.GetSendToCrawler())
	} else if response.GetFinished() != 2 {
		t.Fatalf("number of finished tasks does not match. Expected: %d, Got %d", finished, response.GetFinished())
	} else {
		t.Log("Got all states as expected.")
	}
}
