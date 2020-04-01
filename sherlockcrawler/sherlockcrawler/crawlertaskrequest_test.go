package sherlockcrawler

import (
	"testing"
)

const (
	ADDRESS = "www.hm.edu"
	DONE    = true //TODO
)

func TestCrawlertaskRequestGetterAndSetter(t *testing.T) {
	instance := CrawlerTaskRequest{}

	instance.setAddr(ADDRESS)
	//instance.setDone(DONE) TODO
	instance.setResponse(nil)

	if instance.getAddr() != ADDRESS {
		t.Errorf("Address of the instance does not match the given one. GOT: %s, WANTED: %s", instance.getAddr(), ADDRESS)
		//TODO
		//} else if instance.isDone() != DONE {
		//	t.Errorf("Status of the instance does not match the given one. GOT: %t, WANTED: %t", instance.isDone(), DONE)
	} else if instance.getResponseByReferenz() != nil {
		t.Errorf("Address of the instance does not match the given one. GOT: %p, WANTED: %s", instance.getResponseByReferenz(), "nil")
	} else {
		t.Log("Finished with success.")
	}
}
