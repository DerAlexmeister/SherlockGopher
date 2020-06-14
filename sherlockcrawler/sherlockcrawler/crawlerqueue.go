package sherlockcrawler

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	swd "github.com/ob-algdatii-20ss/SherlockGopher/sherlockwatchdog"
	log "github.com/sirupsen/logrus"
)

/*
CrawlerQueueInterface is a interface for the crawler queue.
*/
type CrawlerQueueInterface interface {
	Watchdog() swd.WatchdogInterface
	SetWatchdog(watchdog swd.WatchdogInterface)
	SetObserver(observer CrawlerObserverInterface)
	Observer() CrawlerObserverInterface
	State() int
	SetState(state int)
	IsEmpty() bool
	GetThisQueue() *map[uint64]*CrawlerTaskRequest
	ContainsTaskID(id uint64) bool
	AppendQueue(task *CrawlerTaskRequest) uint64
	RemoveFromQueue(taskID uint64) bool
	NewQueueStatus() QueueStatus
	GetStatusOfQueue() string
	QueueState() int
	SetQueueState(state int)
	StopQueue()
	CleanQueue()
}

/*
CrawlerQueue will be the queue of the current CrawlerTaskRequest.
*/
type CrawlerQueue struct {
	Queue    map[uint64]*CrawlerTaskRequest
	state    int
	observer CrawlerObserverInterface
	watchdog swd.WatchdogInterface
	mutex    *sync.Mutex
}

/*
QueueStatus struct contains a status of all tasks. E.g undonetaks 10, processingtask 2343, etc.
*/
type QueueStatus struct {
	undoneTask      uint64
	processingTasks uint64
	finishedTasks   uint64
	failedTasks     uint64
}

//TODO JW comments
/*
Watchdog returns the watchdog of the CrawlerQueue.
*/
func (que CrawlerQueue) Watchdog() swd.WatchdogInterface {
	return que.watchdog
}

/*
SetWatchdog is a setter for the watchdog of the CrawlerQueue.
*/
func (que *CrawlerQueue) SetWatchdog(watchdog swd.WatchdogInterface) {
	que.watchdog = watchdog
}

/*
SetObserver is a setter for the observer of the CrawlerQueue.
*/
func (que *CrawlerQueue) SetObserver(observer CrawlerObserverInterface) {
	que.observer = observer
}

/*
Observer returns the observer of the CrawlerQueue.
*/
func (que CrawlerQueue) Observer() CrawlerObserverInterface {
	return que.observer
}

/*
State returns the state of the que.
*/
func (que CrawlerQueue) State() int {
	return que.state
}

/*
SetState sets the state of the que.
*/
func (que *CrawlerQueue) SetState(state int) {
	que.state = state
}

/*
IsEmpty returns true if queue empty.
*/
func (que CrawlerQueue) IsEmpty() bool {
	return len(que.Queue) == 0
}

/*
NewCrawlerQueue will return a new Queue.
*/
func NewCrawlerQueue() CrawlerQueue {
	return CrawlerQueue{
		Queue: map[uint64]*CrawlerTaskRequest{},
		mutex: &sync.Mutex{},
	}
}

/*
GetThisQueue will return a pointer to the current Queue.
*/
func (que CrawlerQueue) GetThisQueue() *map[uint64]*CrawlerTaskRequest {
	return &que.Queue
}

/*
ContainsTaskID will check whether or not a id is already in use or not.
*/
func (que CrawlerQueue) ContainsTaskID(id uint64) bool {
	if _, contains := (*que.GetThisQueue())[id]; !contains {
		return false
	}
	return true
}

/*
Function to produce a random taskID.
*/
func (que *CrawlerQueue) getRandomUserID() uint64 {
	rand.Seed(time.Now().UnixNano())
	for {
		potentialID := rand.Uint64()
		if !que.ContainsTaskID(potentialID) && potentialID > 0 {
			return potentialID
		}
	}
}

/*
AppendQueue will append the current queue with a new CrawlerTaskRequest.
Will return the taskID. In case an error occurred it will return 0.
The taskID is a uint64.
*/
func (que CrawlerQueue) AppendQueue(task *CrawlerTaskRequest) uint64 {
	que.mutex.Lock()

	taskID := que.getRandomUserID()
	var ret uint64 = 0
	if !que.ContainsTaskID(taskID) && task != nil {
		task.setTaskID(taskID)
		(*que.GetThisQueue())[taskID] = task

		que.Watchdog().Aus()

		ret = taskID
	}

	que.mutex.Unlock()
	return ret
}

/*
RemoveFromQueue will remove a task from the queue by a given taskID.
*/
func (que CrawlerQueue) RemoveFromQueue(taskID uint64) bool {
	que.mutex.Lock()
	ret := false
	if taskID > 0 && que.ContainsTaskID(taskID) {
		delete(*que.GetThisQueue(), taskID)
		ret = true
	}
	que.mutex.Unlock()
	return ret
}

/*
getAllTaskIds will return the ids of all tasks in a slice.
*/
func (que CrawlerQueue) getAllTaskIds() []uint64 {
	var ids []uint64
	for k := range *que.GetThisQueue() {
		ids = append(ids, k)
	}
	return ids
}

/*
NewQueueStatus  will return a new queueStatus instance containing
information of all 4 states.
*/
func (que CrawlerQueue) NewQueueStatus() QueueStatus {
	undone, processing, finished, failed := que.getNumberOfStatus()
	return QueueStatus{
		undoneTask:      undone,
		processingTasks: processing,
		finishedTasks:   finished,
		failedTasks:     failed,
	}
}

/*
getAmountOfUndoneTasks will return the number of all tasks with the status undone.
*/
func (status *QueueStatus) getAmountOfUndoneTasks() uint64 {
	return status.undoneTask
}

/*
getAmountOfProcessedTasks will return the number of all tasks with the status processing.
*/
func (status *QueueStatus) getAmountOfProcessedTasks() uint64 {
	return status.processingTasks
}

/*
getAmountOfUndoneTasks will return the number of all tasks with the status undone.
*/
func (status *QueueStatus) getAmountOfFinishedTasks() uint64 {
	return status.finishedTasks
}

/*
getAmountOfUndoneTasks will return the number of all tasks with the status undone.
*/
func (status *QueueStatus) getAmountOfFailedTasks() uint64 {
	return status.failedTasks
}

/*
getNumberOfUndoneTasks will return the number of Tasks with status undone, processing, finished, failed.
*/
func (que *CrawlerQueue) getNumberOfStatus() (uint64, uint64, uint64, uint64) {
	que.mutex.Lock()
	var undone, processing, finished, failed uint64 = 0, 0, 0, 0
	for _, v := range *que.GetThisQueue() {
		state := v.GetTaskState()
		switch state {
		case FINISHED:
			finished++
		case PROCESSING:
			processing++
		case FAILED:
			failed++
		default:
			undone++
		}
	}
	que.mutex.Unlock()
	return undone, processing, finished, failed
}

/*
GetStatusOfQueue will return the status of the queue.
*/
func (que *CrawlerQueue) GetStatusOfQueue() string {
	status := que.NewQueueStatus()
	que.mutex.Lock()
	ret := fmt.Sprintf("States of the Task currently in the Queue. \nUndone: %d, Processing: %d, Finished: %d, Failed: %d",
		status.getAmountOfUndoneTasks(),
		status.getAmountOfProcessedTasks(),
		status.getAmountOfFinishedTasks(),
		status.getAmountOfFailedTasks())

	que.mutex.Unlock()
	return ret
}

/*
QueueState will return the current state of the queue.
*/
func (que CrawlerQueue) QueueState() int {
	return que.state
}

/*
SetQueueState will set the .
*/
func (que *CrawlerQueue) SetQueueState(state int) {
	log.Info("Queue changed state from  to ", que.State(), state)

	switch state {
	case STOP:
		que.state = state
		que.StopQueue()
	case PAUSE:
		que.state = state
	case RUNNING:
		que.state = state
	case IDLE:
		que.state = state
	default:
		log.Warn("Queue got unknown state ", state)
	}
}

/*
StopQueue will stop the queue.
*/
func (que *CrawlerQueue) StopQueue() {
	que.deleteAllTasks()
}

/*
CleanQueue will clean the queue.
*/
func (que *CrawlerQueue) CleanQueue() {
	que.deleteAllTasks()
}

/*
deleteAllTasks will empty the queue.
*/
func (que *CrawlerQueue) deleteAllTasks() {
	log.Info("Queue was emptied")
	que.mutex.Lock()
	que.Queue = make(map[uint64]*CrawlerTaskRequest)
	que.mutex.Unlock()
}
