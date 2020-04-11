package sherlockcrawler

import (
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
NewCrawlerQueue will return a new Queue.
*/
func NewCrawlerQueue() CrawlerQueue {
	return CrawlerQueue{
		Queue: make(map[uint64]*CrawlerTaskRequest),
	}
}

/*
getCurrentQueue will return a pointer to the current Queue.
*/
func (que *CrawlerQueue) getCurrentQueue() *(map[uint64]*CrawlerTaskRequest) {
	return &que.Queue
}

/*
ContainsTaskID will check whether or not a id is allready in use or not.
*/
func (que *CrawlerQueue) ContainsTaskID(id uint64) bool {
	if _, contains := (*que.getCurrentQueue())[id]; !contains {
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
		(*que.getCurrentQueue())[taskid] = task
		return taskid
	}
	return 0
}

/*
RemoveFromQueue will remove a task from the queue by a given Taskid.
*/
func (que *CrawlerQueue) RemoveFromQueue(taskid uint64) bool {
	if taskid > 0 && que.ContainsTaskID(taskid) {
		delete((*que.getCurrentQueue()), taskid)
		return true
	}
	return false
}
