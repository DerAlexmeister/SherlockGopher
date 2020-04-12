package sherlockcrawler

import (
	"fmt"
	"math/rand"
	"time"
)

/*
CrawlerQueue will be the queue of the current CrawlerTaskRequest.
*/
type CrawlerQueue struct {
	Queue map[uint64]*CrawlerTaskRequest
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
NewCrawlerQueue will return a new Queue.
*/
func NewCrawlerQueue() CrawlerQueue {
	return CrawlerQueue{
		Queue: make(map[uint64]*CrawlerTaskRequest),
	}
}

/*
getThisQueue will return a pointer to the current Queue.
*/
func (que *CrawlerQueue) getThisQueue() *(map[uint64]*CrawlerTaskRequest) {
	return &que.Queue
}

/*
ContainsTaskID will check whether or not a id is allready in use or not.
*/
func (que *CrawlerQueue) ContainsTaskID(id uint64) bool {
	if _, contains := (*que.getThisQueue())[id]; !contains {
		return false
	}
	return true
}

/*
Function to produce a random taskid.
*/
func (que *CrawlerQueue) getRandomUserID() uint64 {
	rand.Seed(time.Now().UnixNano())
	for {
		potantialID := rand.Uint64()
		if !que.ContainsTaskID(potantialID) && potantialID > 0 {
			return potantialID
		}
	}
}

/*
AppendQueue will append the current queue with a new CrawlerTaskRequest.
Will return the taskid. Incase an error occurred it will return 0.
The Taskid is a uint64.
*/
func (que *CrawlerQueue) AppendQueue(task *CrawlerTaskRequest) uint64 {
	taskid := que.getRandomUserID()
	if !que.ContainsTaskID(taskid) && task != nil {
		task.setTaskID(taskid)
		(*que.getThisQueue())[taskid] = task
		return taskid
	}
	return 0
}

/*
RemoveFromQueue will remove a task from the queue by a given Taskid.
*/
func (que *CrawlerQueue) RemoveFromQueue(taskid uint64) bool {
	if taskid > 0 && que.ContainsTaskID(taskid) {
		delete((*que.getThisQueue()), taskid)
		return true
	}
	return false
}

/*
getAllTaskIds will return the ids of all tasks in a slice.
*/
func (que *CrawlerQueue) getAllTaskIds() []uint64 {
	var ids []uint64
	for k := range *que.getThisQueue() {
		ids = append(ids, k)
	}
	return ids
}

/*
NewQueueStatus  will return a new queuestatus instance containing
information of all 4 states.
*/
func (que *CrawlerQueue) NewQueueStatus() QueueStatus {
	undone, processing, finished, failed := que.getNumberOfStatus()
	return QueueStatus{
		undonetask:      undone,
		processingtasks: processing,
		finishedtasks:   finished,
		failedtasks:     failed,
	}
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

/*
getAmountOfUndoneTasks will return the number of all tasks with the status undone.
*/
func (status *QueueStatus) getAmountOfFailedTasks() uint64 {
	return status.failedtasks
}

/*
getNumberOfUndoneTasks will return the number of Tasks with status undone, processing, finished, failed.
*/
func (que *CrawlerQueue) getNumberOfStatus() (uint64, uint64, uint64, uint64) {
	var undone, processing, finished, failed uint64 = 0, 0, 0, 0
	for _, v := range *que.getThisQueue() {
		if state := v.getTaskState(); state == FINISHED {
			finished++
		} else if state == PROCESSING {
			processing++
		} else if state == FAILED {
			failed++
		} else {
			undone++
		}
	}
	return undone, processing, finished, failed
}

/*
GetStatusOfQueue will return the
*/
func (que *CrawlerQueue) GetStatusOfQueue() string {
	status := que.NewQueueStatus()
	return fmt.Sprintf("States of the Task currently in the Queue. \nUndone: %d, Processing: %d, Finished: %d, Failed: %d",
		status.getAmountOfUndoneTasks(),
		status.getAmountOfProcessedTasks(),
		status.getAmountOfFinishedTasks(),
		status.getAmountOfFailedTasks())
}
