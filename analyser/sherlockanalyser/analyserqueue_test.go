package sherlockanalyser

import (
	"strconv"
	"testing"
)

func TestNewAnalyserQueue(t *testing.T) {
	que := getTestQueue()

	if que.IsEmpty() {
		t.Error("que failed")
	}
}

func TestContainsAddress(t *testing.T) {
	que := getTestQueue()

	if _, ok, err := que.ContainsAddress("www.1234.de"); err != nil && !ok {
		t.Error("element is missing")
	}
}

func TestContainsID(t *testing.T) {
	que := getTestQueue()

	if id, ok, err := que.ContainsAddress("www.1234.de"); err == nil && ok {
		if !que.ContainsID(id) {
			t.Error("element is missing")
		}
	} else {
		t.Error("element is missing")
	}
}

func TestAppendQueue(t *testing.T) {
	que := getTestQueue()

	task := NewTask(getTestData(9999999))
	if !que.AppendQueue(&task) {
		t.Error("append alement failed")
	}
}

func TestRemoveFromQueue(t *testing.T) {
	que := getTestQueue()

	if !que.RemoveFromQueue(1234) {
		t.Error("remove alement failed")
	}
}

func getTestQueue() AnalyserQueue {
	aq := NewAnalyserQueue()

	var i uint64
	for i = 0; i < 10000; i++ {
		task := NewTask(getTestData(i))
		aq.AppendQueue(&task)
	}

	return aq
}

func getTestData(id uint64) CrawlerData {
	return CrawlerData{
		taskid:            id,
		addr:              "www." + strconv.FormatUint(id, 10) + ".de",
		taskstate:         PROCESSING,
		taskerror:         nil,
		taskerrortry:      0,
		responseHeader:    nil,
		responseBodyBytes: []byte("html"),
		statuscode:        200,
		responseTime:      0,
	}
}
