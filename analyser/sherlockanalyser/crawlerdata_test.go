package sherlockanalyser

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCrawlerDataGetterSetter(t *testing.T) {
	addr := "wAdr"
	responseBody := make([]byte, 0)
	var responseHeader *http.Header = nil
	cData := CrawlerData{}
	statusCode := 1
	taskError := fmt.Errorf("error")
	var id uint64 = 1

	cData.setAddr(addr)
	cData.setResponseBody(responseBody)
	cData.setResponseHeader(responseHeader)
	cData.setResponseTime(delaytime)
	cData.setStatusCode(statusCode)
	cData.setTaskError(taskError)
	cData.setTaskID(id)
}
