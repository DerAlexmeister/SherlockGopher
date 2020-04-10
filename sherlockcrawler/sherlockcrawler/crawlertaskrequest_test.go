package sherlockcrawler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

const (
	id                  = 1
	addr                = "localhost:8080"
	state               = FINISHED
	ltime time.Duration = 10
)

// Testing error
var err error = errors.New("An error")

func TestCrawlerTaskRequest(t *testing.T) {
	task := NewTask()
	task.setTaskID(id)
	task.setAddr(addr)
	task.setTaskState(state)
	task.setTaskError(err)
	task.setResponseTime(ltime)
}

/*
TestMakeRequestForHTML will test the MakeRequestForHTML.
*/
func TestMakeRequestForHTML(t *testing.T) {
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	task := NewTask()
	task.setAddr(server.URL)
	result, err := task.MakeRequestForHTML()
	if err != nil {
		t.Fatal(err)
	}
	defer result.Body.Close()
	bodyBytes, err := ioutil.ReadAll(result.Body)
	if !reflect.DeepEqual([]byte(wanted), bodyBytes) {
		t.Fatal(wanted, string(bodyBytes))
	}
	defer server.Close()
}

/*
TestMakeRequestForHTML will test the MakeRequestForHTML which returned an error.
*/
func TestMakeRequestForHTMLWithError(t *testing.T) {
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	task := NewTask()
	_, err := task.MakeRequestForHTML()
	if err != nil {
		t.Log("An error occured as expected.", err)
	} else {
		t.Fatal("No error occured but was expected.")
	}
	defer server.Close()
}

func TestMakeRequestAndStoreResponse(t *testing.T) {

}
