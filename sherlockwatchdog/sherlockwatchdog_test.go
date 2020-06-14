package sherlockwatchdog

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

type testError struct{}

func (err *testError) Error() string {
	return "expected test-error"
}

func Test_Watch(t *testing.T) {
	sut := NewSherlockWatchdog()
	sut.Watch()
	assert.Equal(t, sut.hits, 1)
}

func Test_Fass(t *testing.T) {
	sut := NewSherlockWatchdog()
	sut.Fass()
	assert.Equal(t, sut.hits, 1)
}

func Test_Aus(t *testing.T) {
	sut := NewSherlockWatchdog()
	sut.Fass()
	sut.Aus()
	assert.Equal(t, sut.hits, 0)
}

func Test_WatchWait(t *testing.T) {
	sut := NewSherlockWatchdog()
	sut.hits = 101
	sut.Watch()
	assert.Equal(t, sut.hits, 102)
}
