package sherlockanalyser

import (
	"math/rand"
	"time"
)

//MAXTASKS will be the maximum amout of possible tasks in one queue.
const MAXTASKS int64 = 9223372036854775807

/*
AnalyserQueue will be the queue of the current AnalyserTaskRequest.
*/
type AnalyserQueue struct {
	Queue map[int64]*AnalyserTaskRequest
}

/*
NewAnalyserQueue will return a new Queue.
*/
func NewAnalyserQueue() AnalyserQueue {
	return AnalyserQueue{
		Queue: make(map[int64]*AnalyserTaskRequest),
	}
}

/*
getCurrentQueue will return a pointer to the current Queue.
*/
func (que *AnalyserQueue) getCurrentQueue() *(map[int64]*AnalyserTaskRequest) {
	return &que.Queue
}

/*
getMAXTASKS will return the const MAXTASKS which represents the maximum amout of current tasks.
*/
func getMAXTASKS() int64 {
	return MAXTASKS
}

/*
ContainsAddress will check whether or not a addr is already in use or not.
*/
func (que *AnalyserQueue) ContainsAddress(addr string) (int64, bool) {
	q := que.getCurrentQueue()

	for key, ele := range *q {
		if ele.addr == addr {
			return key, true
		}
	}

	return -1, false
}

/*
containsID will check whether the queue has already a given id in use.
*/
func (que *AnalyserQueue) ContainsID(id int64) bool {
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
func (que *AnalyserQueue) getRandomTaskID(length int64) int64 {
	rand.Seed(time.Now().UnixNano())

	for {
		potantialID := rand.Int63n(length)
		if !que.ContainsID(potantialID) {
			return potantialID
		}
	}
}

/*
AppendQueue will append the current queue with a new AnalyserTaskRequest.
*/
func (que *AnalyserQueue) AppendQueue(task *AnalyserTaskRequest) bool {
	taskid := que.getRandomTaskID(getMAXTASKS())
	if !que.ContainsID(taskid) {
		(*que.getCurrentQueue())[taskid] = task
		return true
	}

	return false
}

/*
RemoveFromQueue will remove a task from the queue by a given address.
*/
func (que *AnalyserQueue) RemoveFromQueue(taskid int64) bool {
	if !que.ContainsID(taskid) {
		delete((*que.getCurrentQueue()), taskid)
		return true
	}

	return false
}
