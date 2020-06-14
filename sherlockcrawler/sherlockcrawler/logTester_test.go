package sherlockcrawler

import (
	"testing"
)

func Test_contains(t *testing.T) {
	hf := logTester{}

	hf.Run()

}

func TestContains(t *testing.T) {
	letters := []string{"a", "b"}

	if !contains(letters, "a") {
		t.Fatal("contains doesnt work")
	}

	if contains(letters, "c") {
		t.Fatal("contains doesnt work")
	}

}
