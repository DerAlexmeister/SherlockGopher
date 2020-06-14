package sherlockcrawler

import (
	"fmt"
	"testing"

	swd "github.com/ob-algdatii-20ss/SherlockGopher/sherlockwatchdog"
)

//=======================================================================
// Support functions
//=======================================================================

/*
containsElement will check whether a given slice has an element or not.
*/
func containsElement(s []uint64, id uint64) bool {
	for _, i := range s {
		if i == id {
			return true
		}
	}
	return false
}

//=======================================================================
// TESTS
//=======================================================================

func TestInterface(t *testing.T) {
	var i CrawlerQueueInterface

	q := NewCrawlerQueue()
	i = &q

	if i.IsEmpty() != q.IsEmpty() {
		t.Fatal("Interface error")
	}
}

/*
TestQueueHasIDAlreadyInUse tests whether a id is already in use.
*/
func TestQueueHasIDAlreadyInUse(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	fmt.Println(queue)
	injectObserver(&queue)
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("got a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id); !contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		fmt.Println("successfully found ID in taskQueue")
	}
}

func TestGetObserver(t *testing.T) {
	queue := NewCrawlerQueue()
	if queue.Observer() != nil {
		t.Fatal("Observer doesnt work")
	}
}

func TestSetQueueStateIdle(t *testing.T) {
	queue := NewCrawlerQueue()
	queue.SetQueueState(IDLE)
	if queue.State() != 4 {
		t.Fatal("Queue STATE IDLE doesnt work")
	}
}

func TestSetQueueStateDefault(t *testing.T) {
	queue := NewCrawlerQueue()
	queue.SetQueueState(8)
	if queue.State() != 0 {
		t.Fatal("Queue STATE Default doesnt work")
	}
}

/*
TestQueueContainsTaskIDFails tests tries.
*/
func TestQueueContainsTaskIDFails(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("got a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id + 1); contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		fmt.Println("successfully found ID in taskQueue")
	}

}

/*
TestQueueAppendFailesBecauseOfNilPointer will try to make the Append Method fail because of a nullpointer.
*/
func TestQueueAppendFailesBecauseOfNilPointer(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	id := queue.AppendQueue(nil)
	if id == 0 {
		fmt.Println("got a zero return code as expected", 0, id)
	} else {
		t.Fatalf("got a non-zero returncode. Expected: %d, Got: %d", 0, id)
	}

}

/*
TestQueueRemoveTask will append a task to the queue and remove it.
*/
func TestQueueRemoveTask(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("go a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id); !contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		fmt.Println("successfully added task to the taskQueue")
	}
	if works := queue.RemoveFromQueue(id); !works {
		t.Fatalf("couldnt remove task with valid id Expected: %t, Got: %t, ID: %d", true, works, id)
	} else {
		fmt.Println("successfully removed task")
	}
}

/*
TestQueueRemoveTaskFailedBecauseOfZeroID will fail because of a zero id.
*/
func TestQueueRemoveTaskFailedBecauseOfZeroID(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	task := &CrawlerTaskRequest{}
	id := queue.AppendQueue(task)
	if id == 0 {
		t.Fatal("go a 0 id from the AppendQueue-Method")
	} else if contains := queue.ContainsTaskID(id); !contains {
		t.Fatalf("id is not in the queue altought it should be. ID: %d, Contains: %t", id, contains)
	} else {
		fmt.Println("successfully added task to the taskQueue")
	}
	if works := queue.RemoveFromQueue(0); works {
		t.Fatalf("could remove task with unvalid id Expected: %t, Got: %t, ID: %d", false, works, id)
	} else {
		fmt.Println("successfully removed task")
	}
}

/*
TestGetQueueStatus will test the status of the queue and show its status.
*/
func TestGetQueueStatus(t *testing.T) {
	message := "States of the Task currently in the Queue. \nUndone: 1, Processing: 1, Finished: 1, Failed: 1"
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	task1 := &CrawlerTaskRequest{}
	task2 := &CrawlerTaskRequest{}
	task3 := &CrawlerTaskRequest{}
	task4 := &CrawlerTaskRequest{}
	task1.setTaskState(FAILED)
	task2.setTaskState(FINISHED)
	task3.setTaskState(UNDONE)
	task4.setTaskState(PROCESSING)
	id1 := queue.AppendQueue(task1)
	id2 := queue.AppendQueue(task2)
	id3 := queue.AppendQueue(task3)
	id4 := queue.AppendQueue(task4)

	if id1 == 0 || id2 == 0 || id3 == 0 || id4 == 0 {
		t.Fatalf("Missmatching ids excpeted non-zero value but got id1: %d, id2: %d, id3: %d, id4: %d", id1, id2, id3, id4)
	} else {
		fmt.Println("successfully creates 4 different tasks with different states")
	}

	status := queue.GetStatusOfQueue()

	if status != message {
		t.Fatalf("messages does not match Expected: %s, Got: %s", message, status)
	} else {
		fmt.Println("message go as expected.")
	}
}

/*
TestGetAllTaskIds if GetAllTaskIds will return all ids.
*/
func TestGetAllTaskIds(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	task1 := &CrawlerTaskRequest{}
	task2 := &CrawlerTaskRequest{}
	id1 := queue.AppendQueue(task1)
	id2 := queue.AppendQueue(task2)

	if id1 == 0 || id2 == 0 {
		t.Fatalf("Missmatching ids excpeted non-zero value but got id1: %d, id2: %d", id1, id2)
	} else {
		fmt.Println("successfully creates 4 different tasks with different states")
	}

	ids := queue.getAllTaskIds()

	if length := len(ids); length != 2 {
		t.Fatalf("The length of the slice of ids does not be match the real one. Expected: %d, Got: %d", 2, length)
	} else if contains := containsElement(ids, id1); !contains {
		t.Fatalf("element in should be in map but wasnt there. Couldnt find: %d", id1)
	} else if contains := containsElement(ids, id2); !contains {
		t.Fatalf("element in should be in map but wasnt there. Couldnt find: %d", id2)
	} else {
		fmt.Println("successfully found all elements")
	}
}

/*
TestCrawlerQueue_QueueState tests QueueState
*/
func TestCrawlerQueue_QueueState(t *testing.T) {
	queue := NewCrawlerQueue()
	injectObserver(&queue)
	queue.SetQueueState(RUNNING)

	if queue.QueueState() != RUNNING {
		t.Fatalf("expected running but was %v", queue.QueueState())
	}
}

/*
TestCrawlerQueue_SetQueueState tests SetQueueState
*/
func TestCrawlerQueue_SetQueueState(t *testing.T) {
	queue := NewCrawlerQueue()
	injectObserver(&queue)
	queue.SetQueueState(STOP)

	if queue.QueueState() != STOP {
		t.Fatalf("expected stop but was %v", queue.QueueState())
	}

	queue.SetQueueState(PAUSE)

	if queue.QueueState() != PAUSE {
		t.Fatalf("expected pause but was %v", queue.QueueState())
	}

	queue.SetQueueState(RUNNING)

	if queue.QueueState() != RUNNING {
		t.Fatalf("expected running but was %v", queue.QueueState())
	}
}

/*
TestCrawlerQueue_StopQueue tests StopQueue
*/
func TestCrawlerQueue_StopQueue(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	task1 := &CrawlerTaskRequest{}
	task2 := &CrawlerTaskRequest{}
	queue.AppendQueue(task1)
	queue.AppendQueue(task2)

	queue.StopQueue()

	if len(queue.getAllTaskIds()) != 0 {
		t.Fatalf("expected empty but was %v", len(queue.getAllTaskIds()))
	}
}

/*
TestCrawlerQueue_CleanQueue tests CleanQueue
*/
func TestCrawlerQueue_CleanQueue(t *testing.T) {
	queue := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	injectObserver(&queue)
	task1 := &CrawlerTaskRequest{}
	task2 := &CrawlerTaskRequest{}
	queue.AppendQueue(task1)
	queue.AppendQueue(task2)

	queue.CleanQueue()

	if len(queue.getAllTaskIds()) != 0 {
		t.Fatalf("expected empty but was %v", len(queue.getAllTaskIds()))
	}
}

func injectObserver(que *CrawlerQueue) {
	observer := CrawlerObserver{}

	observer.SetQueue(que)
}

func TestCleanCommand(t *testing.T) {
	observer := CrawlerObserver{}
	queue := NewCrawlerQueue()
	queue2 := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	queue.SetWatchdog(&watchdog)
	observer.SetQueue(&queue)
	observer.SetStreamerQue(&queue2)
	task1 := &CrawlerTaskRequest{}
	task2 := &CrawlerTaskRequest{}
	queue.AppendQueue(task1)
	queue.AppendQueue(task2)

	queue.Observer().CleanCommand()

	if len(queue.getAllTaskIds()) != 0 {
		t.Fatalf("expected empty but was %v", len(queue.getAllTaskIds()))
	}
}

func TestObserverUpdateStateDefault(t *testing.T) {
	observer := CrawlerObserver{}
	observer.SetCrawler(NewSherlockCrawlerService())

	observer.UpdateState(5)

	if observer.Crawler().State() != 2 {
		t.Fatalf("Observer Update STATE Default doesnt work")
	}
}
