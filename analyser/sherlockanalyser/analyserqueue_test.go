package sherlockanalyser

import (
	"github.com/ob-algdatii-20ss/SherlockGopher/sherlockwatchdog"
	"testing"
)

func getTestQueue() AnalyserQueue {
	watchDog := sherlockwatchdog.NewSherlockWatchdog()
	observer := NewAnalyserObserver()
	que := NewAnalyserQueue()
	que.SetWatchdog(&watchDog)
	que.SetObserver(observer)

	return que
}

//=======================================================================
// TESTS
//=======================================================================

func TestObserver(t *testing.T) {
	que := getTestQueue()

	if que.Observer() == nil {
		t.Fatal("Observer is nil")
	}
}

func TestSetState_STOP(t *testing.T) {
	que := getTestQueue()
	cData := GetExampleCData()
	task := AnalyserTaskRequest{crawlerData: &cData}
	que.AppendQueue(&task)
	que.SetState(STOP)

	if !que.IsEmpty() {
		t.Fatal("Queue is expected to be empty but is not")
	}

	if que.State() != STOP {
		t.Fatal("State is expected to be STOP but is not")
	}
}

func TestSetState_PAUSE(t *testing.T) {
	que := getTestQueue()
	que.SetState(PAUSE)

	if que.State() != PAUSE {
		t.Fatal("State is expected to be PAUSE but is not")
	}
}

func TestSetState_RUNNING(t *testing.T) {
	que := getTestQueue()
	que.SetState(RUNNING)

	if que.State() != RUNNING {
		t.Fatal("State is expected to be RUNNING but is not")
	}
}

func TestSetState_IDLE(t *testing.T) {
	que := getTestQueue()
	que.SetState(IDLE)

	if que.State() != IDLE {
		t.Fatal("State is expected to be IDLE but is not")
	}
}

func TestRemoveFromQueue(t *testing.T) {
	que := getTestQueue()

	cData := GetExampleCData()
	task := AnalyserTaskRequest{crawlerData: &cData}
	que.AppendQueue(&task)

	if !que.RemoveFromQueue(task.getID()) {
		t.Fatal("Remove failed")
	}

	if !que.IsEmpty() {
		t.Fatal("Remove didnt remove")
	}
}

func TestGetAllTaskIds(t *testing.T) {
	que := getTestQueue()

	cData := GetExampleCData()
	task := AnalyserTaskRequest{crawlerData: &cData}
	que.AppendQueue(&task)

	ids := que.getAllTaskIds()

	if len(ids) != 1 {
		t.Fatal("Wrong amount of ids")
	}

	if ids[0] != task.getID() {
		t.Fatal("Wrong id returned")
	}
}

func TestGetStatusOfQueue(t *testing.T) {
	que := getTestQueue()


	cData := GetExampleCData()
	cData.setAddr(cData.addr + "1")
	task1 := AnalyserTaskRequest{crawlerData: &cData}
	cData2 := GetExampleCData()
	cData2.setAddr(cData2.addr + "2")
	task2 := AnalyserTaskRequest{crawlerData: &cData2}
	cData3 := GetExampleCData()
	cData3.setAddr(cData3.addr + "3")
	task3 := AnalyserTaskRequest{crawlerData: &cData3}
	cData4 := GetExampleCData()
	cData4.setAddr(cData4.addr + "4")
	task4 := AnalyserTaskRequest{crawlerData: &cData4}
	cData5 := GetExampleCData()
	cData5.setAddr(cData5.addr + "5")
	task5 := AnalyserTaskRequest{crawlerData: &cData5}
	cData6 := GetExampleCData()
	cData6.setAddr(cData6.addr + "6")
	task6 := AnalyserTaskRequest{crawlerData: &cData6}
	task1.SetState(0)
	task2.SetState(1)
	task3.SetState(2)
	task4.SetState(3)
	task5.SetState(4)
	task6.SetState(5)
	que.AppendQueue(&task1)
	que.AppendQueue(&task2)
	que.AppendQueue(&task3)
	que.AppendQueue(&task4)
	que.AppendQueue(&task5)
	que.AppendQueue(&task6)

	queStatus := que.GetStatusOfQueue()
	expected := "States of the Task currently in the Queue. \nUndone: 1, Processing: 1, CrawlerError: 1, Saving: 1, SendToCrawler: 1, Finished: 1"

	if queStatus != expected {
		t.Fatal("Returned status does not equal expected")
	}
}

func TestStopQueue(t *testing.T) {
	que := getTestQueue()

	cData := GetExampleCData()
	task := AnalyserTaskRequest{crawlerData: &cData}
	que.AppendQueue(&task)

	que.StopQueue()

	if !que.IsEmpty() {
		t.Fatal("Remove didnt remove")
	}
}
