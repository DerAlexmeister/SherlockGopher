package sherlockcrawler

import (
	"math/rand"
	"time"
)

//MAXTASKS will be the maximum amout of possible tasks in one queue.
const MAXTASKS int64 = 9223372036854775807

/*
CrawlerQueue will be the queue of the current CrawlerTaskRequest.
*/
type CrawlerQueue struct {
	Queue map[int64]*CrawlerTaskRequest
}

/*
NewCrawlerQueue will return a new Queue.
*/
func (que *CrawlerQueue) NewCrawlerQueue() CrawlerQueue {
	return CrawlerQueue{
		Queue: make(map[int64]*CrawlerTaskRequest),
	}
}

/*
getCurrentQueue will return a pointer to the current Queue.
*/
func (que *CrawlerQueue) getCurrentQueue() *(map[int64]*CrawlerTaskRequest) {
	return &que.Queue
}

/*
getMAXTASKS will return the const MAXTASKS which represents the maximum amout of current tasks.
*/
func getMAXTASKS() int64 {
	return MAXTASKS
}

/*
ContainsAddress will check whether or not a addr is allready in use or not.
*/
func (que *CrawlerQueue) ContainsAddress(addr int64) bool {
	if _, contains := (*que.getCurrentQueue())[addr]; !contains {
		return false
	}
	return true
}

/*
containsID will check whether the queue has already a given id in use.
*/
func (que *CrawlerQueue) containsID(id int64) bool {
	_, inMap := (*que.getCurrentQueue())[id]
	return inMap
}

/*
Function to produce a random taskid.
*/
func (que *CrawlerQueue) getRandomUserID(length int64) int64 {
	rand.Seed(time.Now().UnixNano())
	for {
		potantialID := rand.Int63n(length)
		if !que.containsID(potantialID) {
			return potantialID
		}
	}
}

/*
AppendQueue will append the current queue with a new CrawlerTaskRequest.
*/
func (que *CrawlerQueue) AppendQueue(task *CrawlerTaskRequest) bool {
	taskid := que.getRandomUserID(getMAXTASKS())
	if !que.ContainsAddress(taskid) {
		(*que.getCurrentQueue())[taskid] = task
		return true
	}
	return false
}

/*
RemoveFromQueue will remove a task from the queue by a given address.
*/
func (que *CrawlerQueue) RemoveFromQueue(taskid int64) bool {
	if !que.ContainsAddress(taskid) {
		delete((*que.getCurrentQueue()), taskid)
		return true
	}
	return false
}
