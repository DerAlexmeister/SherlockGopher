package sherlockanalyser

import (
	"fmt"
	"math/rand"
	"time"
)

/*
AnalyserQueue will be the queue of the current AnalyserTaskRequest.
*/
type AnalyserQueue struct {
	Queue map[uint64]*AnalyserTaskRequest
}

/*
NewAnalyserQueue will return a new Queue.
*/
func NewAnalyserQueue() AnalyserQueue {
	return AnalyserQueue{
		Queue: make(map[uint64]*AnalyserTaskRequest),
	}
}

/*
getCurrentQueue will return a pointer to the current Queue.
*/
func (que *AnalyserQueue) getCurrentQueue() *(map[uint64]*AnalyserTaskRequest) {
	return &que.Queue
}

/*
ContainsAddress will check whether or not a addr is already in use or not.
*/
func (que *AnalyserQueue) ContainsAddress(addr string) (uint64, bool, error) {
	q := que.getCurrentQueue()

	for key, ele := range *q {
		if ele.addr == addr {
			return key, true, nil
		}
	}

	return 0, false, fmt.Errorf("id is not in the queue")
}

/*
containsID will check whether the queue has already a given id in use.
*/
func (que *AnalyserQueue) ContainsID(id uint64) bool {
	_, inMap := que.Queue[id]

	return inMap
}

/*
IsEmpty returns true if queue empty
*/
func (que *AnalyserQueue) IsEmpty() bool {
	return len(que.Queue) == 0
}

/*
Function to produce a random taskid.
*/
func (que *AnalyserQueue) getRandomTaskID() uint64 {
	rand.Seed(time.Now().UnixNano())

	for {
		potantialID := rand.Uint64()
		if !que.ContainsID(potantialID) {
			return potantialID
		}
	}
}

/*
AppendQueue will append the current queue with a new AnalyserTaskRequest.
*/
func (que *AnalyserQueue) AppendQueue(task *AnalyserTaskRequest) bool {
	taskid := que.getRandomTaskID()
	if !que.ContainsID(taskid) {
		// Hier TASKID einzufügen
		(*que.getCurrentQueue())[taskid] = task
		return true
	}

	return false
}

/*
RemoveFromQueue will remove a task from the queue by a given address.
*/
func (que *AnalyserQueue) RemoveFromQueue(taskid uint64) bool {
	if !que.ContainsID(taskid) {
		delete((*que.getCurrentQueue()), taskid)
		return true
	}

	return false
}

/*
QueueStatus struct contains a status of all tasks. E.g undonetaks 10, processingtask 2343, ...
*/
type QueueStatus struct {
	undonetask      uint64
	processingtasks uint64
	finishedtasks   uint64
	failedtasks     uint64
}

/*
NewQueueStatus  will return a new queuestatus instance containing
information of all 4 states.
*/
func (que *AnalyserQueue) NewQueueStatus() QueueStatus {
	undone, processing, finished := que.getNumberOfStatus()
	return QueueStatus{
		undonetask:      undone,
		processingtasks: processing,
		finishedtasks:   finished,
	}
}

/*
getNumberOfUndoneTasks will return the number of Tasks with status undone, processing, finished, failed.
*/
func (que *AnalyserQueue) getNumberOfStatus() (uint64, uint64, uint64) {
	var undone, processing, finished uint64 = 0, 0, 0
	for _, v := range *que.getCurrentQueue() {
		if state := v.getTaskState(); state == FINISHED {
			finished++
		} else if state == PROCESSING {
			processing++
		} else {
			undone++
		}
	}
	return undone, processing, finished
}

/*
GetStatusOfQueue will return the
*/
func (que *AnalyserQueue) GetStatusOfQueue() string {
	status := que.NewQueueStatus()
	return fmt.Sprintf("States of the Task currently in the Queue. \nUndone: %d, Processing: %d, Finished: %d",
		status.getAmountOfUndoneTasks(),
		status.getAmountOfProcessedTasks(),
		status.getAmountOfFinishedTasks())
}

/*
getAmountOfUndoneTasks will return the number of all tasks with the status undone.
*/
func (status *QueueStatus) getAmountOfUndoneTasks() uint64 {
	return status.undonetask
}

/*
getAmountOfProcessedTasks will return the number of all tasks with the status processing.
*/
func (status *QueueStatus) getAmountOfProcessedTasks() uint64 {
	return status.processingtasks
}

/*
getAmountOfUndoneTasks will return the number of all tasks with the status undone.
*/
func (status *QueueStatus) getAmountOfFinishedTasks() uint64 {
	return status.finishedtasks
}
