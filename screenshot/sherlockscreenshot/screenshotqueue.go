package sherlockscreenshot

import (
	"math/rand"
	"sync"
	"time"
)

/*
ScreenshotQueue will store all urls that weren't screenshoted yet.
*/
type ScreenshotQueue struct {
	Queue map[uint64]*ScreenshotTaskRequest
	mutex *sync.Mutex
}

/*
NewScreenshotQueue will return a new Queue.
*/
func NewScreenshotQueue() ScreenshotQueue {
	return ScreenshotQueue{
		Queue: map[uint64]*ScreenshotTaskRequest{},
		mutex: &sync.Mutex{},
	}
}

/*
getThisQueue will return a pointer to the current Queue.
*/
func (que *ScreenshotQueue) getThisQueue() *map[uint64]*ScreenshotTaskRequest {
	return &que.Queue
}

/*
ContainsTaskID will check whether or not a id is already in use or not.
*/
func (que *ScreenshotQueue) ContainsTaskID(id uint64) bool {
	if _, contains := (*que.getThisQueue())[id]; !contains {
		return false
	}
	return true
}

/*
getRandomTaskID is a function that produces a random taskID.
*/
func (que *ScreenshotQueue) getRandomTaskID() uint64 {
	rand.Seed(time.Now().UnixNano())
	for {
		potentialID := rand.Uint64()
		if !que.ContainsTaskID(potentialID) && potentialID > 0 {
			return potentialID
		}
	}
}

/*
AppendQueue will append the current queue with a new ScreenshotTaskRequest.
Will return the taskID. In case an error occurred it will return 0.
The taskID is a uint64.
*/
func (que *ScreenshotQueue) AppendQueue(task *ScreenshotTaskRequest) uint64 {
	que.mutex.Lock()
	taskID := que.getRandomTaskID()
	var ret uint64 = 0
	if !que.ContainsTaskID(taskID) && task != nil {
		task.SetID(taskID)
		(*que.getThisQueue())[taskID] = task
		ret = taskID
	}
	que.mutex.Unlock()
	return ret
}

/*
RemoveFromQueue will remove a task from the queue by a given taskID.
*/
func (que *ScreenshotQueue) RemoveFromQueue(taskID uint64) bool {
	que.mutex.Lock()
	ret := false
	if taskID > 0 && que.ContainsTaskID(taskID) {
		delete(*que.getThisQueue(), taskID)
		ret = true
	}
	que.mutex.Unlock()
	return ret
}
