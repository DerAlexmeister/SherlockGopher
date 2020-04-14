package sherlockcrawler

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

const (
	testid                       = 1
	testaddr                     = "localhost:8080"
	testtime       time.Duration = 10
	teststatuscode               = 200
)

// Testing error
var erro error = errors.New("testerror")

func TestMin(t *testing.T) {
	a := 1
	b := 2
	ret := min(a, b)
	if ret != a {
		t.Fatalf("Min function: expected %d but was %d", a, ret)
	} else {
		t.Log("Min fuction works")
	}
}

func TestGetChunkSize(t *testing.T) {
	ret := getChunkSize()
	if ret != 1024 {
		t.Fatalf("Min function: expected %d but was %d", 1024, ret)
	} else {
		t.Log("Min fuction works")
	}
}
