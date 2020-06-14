package sherlockcrawler

import (
	"fmt"
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
	wantedString  string        = "This should be in the body of the HTTP response."
)

// Testing error
var err error = errors.New("An error")

func TestNewCrawlerTaskRequest(t *testing.T) {
	header := http.Header{}
	task := NewCrawlerTaskRequest(0,
		"www.helpme.de", 0, fmt.Errorf(""),
		0, nil, &header, "html",
		[]byte("html"), 200, 2000)
	id := task.GetTaskID()
	if id != uint64(0) {
		fmt.Println(id)
		t.Fatal("Test failed")
	}
}

/*
TestCrawlerTaskRequest will test a crawler task request in general so its getters and setters.
*/
func TestCrawlerTaskRequest(t *testing.T) {
	task := NewTask()
	task.setTaskID(id)
	task.setAddr(addr)
	task.setTaskState(state)
	task.setTaskError(err)
	task.setResponseTime(ltime)
	task.setTryIfError(errorsappered)
	task.setStatusCode(statuscode)

	switch {
	case task.GetAddr() != addr:
		t.Fatalf("Address's do not match. Expected: %s, Got: %s", task.GetAddr(), addr)
	case task.GetTaskID() != id:
		t.Fatalf("IDs do not match. Expected: %d, Got: %d", task.GetTaskID(), id)
	case task.GetTaskState() != state:
		t.Fatalf("Taskstates do not match. Expected: %s, Got: %s", string(task.GetTaskState()), string(state))
	case task.GetTaskError() != err:
		t.Fatalf("Taskerror do not match. Expected: %s, Got: %s", task.GetTaskError().Error(), err.Error())
	case task.GetResponseTime() != ltime:
		t.Fatalf("Responsetime do not match. Expected: %d, Got: %d", int(task.GetResponseTime()), int(ltime))
	case task.GetStatusCode() != statuscode:
		t.Fatalf("Statuscode do not match. Expected: %d, Got: %d", task.GetStatusCode(), statuscode)
	case task.GetTryError() != errorsappered:
		t.Fatalf("Amount of errors does not match. Expected: %d, Got: %d", task.GetTryError(), errorsappered)
	default:
		fmt.Println("Getter and Setter work well (Not Tested are all response related, see later tests).")
	}
}

/* // TODO: TEST IS DEAD
TestMakeRequestForHTML will test the MakeRequestForHTML.

func TestMakeRequestForHTML(t *testing.T) {
	wanted := wantedString
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		_, err = w.Write([]byte(wanted))
		if err != nil {
			t.Fatal(err)
		}
	}))
	task := NewTask()
	task.setAddr(server.URL)
	result := task.MakeRequestAndStoreResponse(nil)

	if !result {
		t.Fatal("an error occurred during execution")
	}
	defer server.Close()
}*/

/* // TODO: TEST IS DEAD
TestMakeRequestForHTML will test the MakeRequestForHTML which returned an error.

//nolint: bodyclose
func TestMakeRequestForHTMLWithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(wantedString))
		if err != nil {
			t.Fatal(err)
		}
	}))
	task := NewTask()

	ret := task.MakeRequestAndStoreResponse(nil)

	if !ret {
		fmt.Println("An error occurred as expected.")
	} else {
		t.Fatal("No error occurred but was expected.")
	}
	defer server.Close()
}*/

/*
TestMakeRequestAndStoreResponse will make a response and store all responses.
*/
func TestMakeRequestAndStoreResponse(t *testing.T) {
	wanted := wantedString
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(wanted))
		if err != nil {
			t.Fatal(err)
		}
	}))
	task := NewTask()
	task.setAddr(server.URL)
	result := task.MakeRequestAndStoreResponse(nil)
	if result {
		body := task.GetResponseByReference().Body
		switch {
		case body == nil:
			t.Fatalf("Got a nil reference to the response Expected: not nil pointer, Got %p", body)
		case !reflect.DeepEqual([]byte(wanted), task.GetResponseBodyInBytes()):
			t.Fatal("The byte version of the body are not equal.")
		case !reflect.DeepEqual(wanted, task.GetResponseBody()):
			t.Fatalf("The body's of the response task.GetResponseBody() and wanted(This should be in the body of the HTTP response.) are different")
		case !reflect.DeepEqual(task.GetResponse().Header, task.GetResponseHeader()):
			t.Fatalf("Response.Header and GetResponseHeader returned different results. Expected: %v, Got %v", task.GetResponse().Header, task.GetResponseHeader())
		case task.GetResponse().StatusCode != 200 && task.GetResponse().StatusCode != task.GetStatusCode():
			t.Fatalf("Response status code retrieved on 2 different ways returned different results Expected: %d, Got %d", task.GetResponse().StatusCode, task.GetStatusCode())
		default:
			fmt.Println("Successfully tested missing getters compared to the test TestCrawlerTaskRequest")
		}
		defer body.Close()
	} else {
		t.Fatal("A Problem occurred while trying to call MakeRequestAndStoreResponse")
	}

	defer server.Close()
}

/*
TestMakeRequestAndStoreResponseWithEmptyAddrField will test a request with empty addr field.
*/
func TestMakeRequestAndStoreResponseWithEmptyAddrField(t *testing.T) {
	wanted := wantedString
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(wanted))
		if err != nil {
			t.Fatal(err)
		}
	}))
	task := NewTask()
	result := task.MakeRequestAndStoreResponse(nil)
	fmt.Println(result)
	fmt.Println(task.GetTaskError().Error())
	fmt.Println(task.GetTryError())
	fmt.Println(task.GetTaskState())
	if !result && task.GetTaskError().Error() == "cannot process a task with an empty address field" && task.GetTryError() > 0 && task.GetTaskState() == FAILED {
	} else {
		t.Fatal("A Problem occurred while trying to call MakeRequestAndStoreResponse should return an error but did not")
	}

	defer server.Close()
}
