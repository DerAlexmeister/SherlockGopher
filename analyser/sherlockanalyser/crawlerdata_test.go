package sherlockanalyser

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

/*
TestCrawlerDataGetterSetter.
*/
func TestCrawlerDataGetterSetter(t *testing.T) {
	addr := "wAdr"
	responseBody := make([]byte, 0)
	responseHeader := http.Header{}
	responseHeader.Add("Content-Language", "de")
	cData := CrawlerData{}
	statusCode := 1
	var startDelayTime time.Duration = 20
	taskError := fmt.Errorf("error")
	var id uint64 = 1

	cData.setAddr(addr)
	cData.setResponseBody(responseBody)
	cData.setResponseHeader(&responseHeader)
	cData.setResponseTime(startDelayTime)
	cData.setStatusCode(statusCode)
	cData.setTaskError(taskError)
	cData.setTaskID(id)

	if cData.getTaskID() != id {
		t.Errorf("id ist wrong")
	}
	if cData.getTaskError() != taskError {
		t.Errorf("taskError ist wrong")
	}
	if cData.getStatusCode() != statusCode {
		t.Errorf("statusCode ist wrong")
	}
	if cData.getResponseTime() != startDelayTime {
		t.Errorf("responseTime ist wrong")
	}
	if cData.getResponseHeader().Get("Content-Language") != responseHeader.Get("Content-Language") {
		t.Errorf("responseHeader ist wrong")
	}
	if len(cData.getResponseBody()) != len(responseBody) {
		t.Errorf("responseBody ist wrong")
	}
	if cData.getAddr() != addr {
		t.Errorf("addr ist wrong")
	}
}
