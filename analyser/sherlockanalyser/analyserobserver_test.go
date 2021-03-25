package sherlockanalyser

import (
	"testing"

	"github.com/DerAlexx/SherlockGopher/sherlockwatchdog"
)

//=======================================================================
// TESTS
//=======================================================================

func TestSetAnalyser(t *testing.T) {
	observer := NewAnalyserObserver()
	analyser := NewAnalyserServiceHandler()
	observer.SetAnalyser(analyser)

	if observer.Analyser() == nil {
		t.Fatal("Analyser is nil")
	}
}

func TestSetQueue(t *testing.T) {
	observer := NewAnalyserObserver()
	queue := NewAnalyserQueue()
	observer.SetQueue(&queue)

	if observer.Queue() == nil {
		t.Fatal("Queue is nil")
	}
}

func TestUpdateState(t *testing.T) {
	observer := NewAnalyserObserver()
	queue := NewAnalyserQueue()
	observer.SetQueue(&queue)
	analyser := NewAnalyserServiceHandler()
	observer.SetAnalyser(analyser)

	observer.UpdateState(1)

	if analyser.State() != 1 {
		t.Fatal("Analyser has wrong state")
	}
	if queue.State() != 1 {
		t.Fatal("Queue has wrong state")
	}
}

func TestCleanCommand(t *testing.T) {
	observer := NewAnalyserObserver()
	queue := NewAnalyserQueue()
	cData := GetExampleCData()
	task := AnalyserTaskRequest{crawlerData: &cData}
	observer.SetQueue(&queue)
	wd := sherlockwatchdog.NewSherlockWatchdog()
	queue.SetWatchdog(&wd)
	queue.AppendQueue(&task)

	observer.CleanCommand()

	if !queue.IsEmpty() {
		t.Fatal("Queue is not cleaned")
	}
}
