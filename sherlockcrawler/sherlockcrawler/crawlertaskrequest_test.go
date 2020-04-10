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
	id                          = 1
	addr                        = "localhost:8080"
	state                       = FINISHED
	ltime         time.Duration = 10
	errorsappered               = 3
	statuscode                  = 200
)

// Testing error
var err error = errors.New("An error")

/*
TestCrawlerTaskRequest will test a crawlertaskrequest in general so its getters and setters.
*/
func TestCrawlerTaskRequest(t *testing.T) {
	task := NewTask()
	task.setTaskID(id)
	task.setAddr(addr)
	task.setTaskState(state)
	task.setTaskError(err)
	task.setResponseTime(ltime)
	task.setTrysIfError(errorsappered)
	task.setStatusCode(statuscode)

	if task.getAddr() != addr {
		t.Fatalf("Address's do not match. Expected: %s, Got: %s", task.getAddr(), addr)
	} else if task.getTaskID() != id {
		t.Fatalf("IDs do not match. Expected: %d, Got: %d", task.getTaskID(), id)
	} else if task.getTaskState() != state {
		t.Fatalf("Taskstates do not match. Expected: %s, Got: %s", string(task.getTaskState()), string(state))
	} else if task.getTaskError() != err {
		t.Fatalf("Taskerror do not match. Expected: %s, Got: %s", string(task.getTaskError().Error()), string(err.Error()))
	} else if task.getResponseTime() != ltime {
		t.Fatalf("Responsetime do not match. Expected: %d, Got: %d", int(task.getResponseTime()), int(ltime))
	} else if task.getStatusCode() != statuscode {
		t.Fatalf("Statuscode do not match. Expected: %d, Got: %d", task.getStatusCode(), statuscode)
	} else if task.getTrysError() != errorsappered {
		t.Fatalf("Amount of errors does not match. Expected: %d, Got: %d", task.getTrysError(), errorsappered)
	} else {
		t.Log("Getter and Setter work well (Not Tested are all response related, see later tests).")
	}
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

func TestMakeRequestAndStoreResponse(t *testing.T) { //TODO
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	task := NewTask()
	task.setAddr(server.URL)
	result := task.MakeRequestAndStoreResponse()
	if result {
		if task.getResponseByReferenz() == nil {

		} else if !reflect.DeepEqual([]byte(wanted), task.getResponseBodyInBytes()) {

		} else if !reflect.DeepEqual(wanted, task.getResponseBody()) {

		}
	} else {
		t.Fatal("A Problem occrued while trying to call MakeRequestAndStoreResponse")
	}

	defer server.Close()
}

/*
TestMakeRequestAndStoreResponseWithEmptyAddrField will test a request with empty addr field.
*/
func TestMakeRequestAndStoreResponseWithEmptyAddrField(t *testing.T) {
	wanted := "This should be in the body of the HTTP response."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(wanted))
	}))
	task := NewTask()
	result := task.MakeRequestAndStoreResponse()
	if !result && task.getTaskError().Error() == "cannot process a task with an empty address field" && task.getTrysError() > 0 && task.getTaskState() == FAILED {
		t.Log("An error occurred as expected.")
	} else {
		t.Fatal("A Problem occrued while trying to call MakeRequestAndStoreResponse should return an error but did not")
	}

	defer server.Close()
}
