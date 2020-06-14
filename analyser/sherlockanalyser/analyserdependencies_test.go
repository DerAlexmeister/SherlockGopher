package sherlockanalyser

import (
	"testing"
)

func TestNewAnalyserDependencies(t *testing.T) {
	dep := NewAnalyserDependencies()

	if dep == nil {
		t.Fatal("Deps is nil")
	}
}
