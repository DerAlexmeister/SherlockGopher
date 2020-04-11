package sherlockcrawler

import (
	"context"
	"fmt"
	"testing"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
)

/*
BenchmarkCreateTask benchmark test for the createtask function
to see how many task can be created in a secound.
*/
func BenchmarkCreateTask(b *testing.B) {
	service := NewSherlockCrawlerService()
	for n := 0; n < b.N; n++ {
		response := proto.CrawlTaskCreateResponse{}
		if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%d", n)}, &response); err != nil {
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
		if err := service.CreateTask(context.TODO(), &proto.CrawlTaskCreateRequest{Url: fmt.Sprintf("%d", n)}, &response); err != nil {
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

	if len := len(*queue.getCurrentQueue()); len != 0 {
		b.Fatalf("queue is not empty. Excpeted: 0, Got: %d", len)
	} else {
		b.Log("Queue is empty.")
	}

}
