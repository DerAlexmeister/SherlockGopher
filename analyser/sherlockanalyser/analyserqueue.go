package sherlockanalyser

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	swd "github.com/ob-algdatii-20ss/SherlockGopher/sherlockwatchdog"
	log "github.com/sirupsen/logrus"
)

/*
AnalyserQueue will be the queue of the current CrawlerTaskRequest.
*/
type AnalyserQueue struct {
	Queue    map[uint64]*AnalyserTaskRequest
	observer AnalyserObserverInterface
	watchdog swd.WatchdogInterface
	state    int
	mutex    *sync.Mutex
	saver    neoSaverInterface
	cache 	AnalyserCacheInterface
}

/*
Watchdog returns the watchdog stored in the AnalyserQueue.
*/
func (que *AnalyserQueue) Watchdog() swd.WatchdogInterface {
	return que.watchdog
}

/*
SetWatchdog is a setter for the watchdog stored in the AnalyserQueue.
*/
func (que *AnalyserQueue) SetWatchdog(watchdog swd.WatchdogInterface) {
	que.watchdog = watchdog
}

/*
SetObserver is a setter for the observer stored in the AnalyserQueue.
*/
func (que *AnalyserQueue) SetObserver(observer AnalyserObserverInterface) {
	que.observer = observer
}

/*
Observer returns the observer stored in the AnalyserQueue.
*/
func (que *AnalyserQueue) Observer() AnalyserObserverInterface {
	return que.observer
}

/*
State returns the state of the queue.
*/
func (que *AnalyserQueue) State() int {
	return que.state
}

/*
SetState sets the state of the queue.
*/
func (que *AnalyserQueue) SetState(state int) {
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
	}
}

/*
QueueStatus struct contains a status of all tasks. E.g undonetaks 10, processingtask 2343, etc.
*/
type QueueStatus struct {
	undoneTasks        uint64
	processingTasks    uint64
	crawlerErrorTasks  uint64
	savingTasks        uint64
	sendToCrawlerTasks uint64
	finishedTasks      uint64
}

/*
NewAnalyserQueue will return a new Queue.
*/
func NewAnalyserQueue() AnalyserQueue {
	saver := newNeoSaver()
	cache := NewAnalyserCache()
	return AnalyserQueue{
		Queue: map[uint64]*AnalyserTaskRequest{},
		state: RUNNING,
		mutex: &sync.Mutex{},
		saver: &saver,
		cache: &cache,
	}
}

/*
getThisQueue will return a pointer to the current Queue.
*/
func (que *AnalyserQueue) getThisQueue() *map[uint64]*AnalyserTaskRequest {
	return &que.Queue
}

/*
ContainsTaskID will check whether or not a id is already in use or not.
*/
func (que *AnalyserQueue) ContainsTaskID(id uint64) bool {
	if _, contains := (*que.getThisQueue())[id]; !contains {
		return false
	}
	return true
}

/*
getRandomTaskID is a function that produces a random taskID.
*/
func (que *AnalyserQueue) getRandomTaskID() uint64 {
	rand.Seed(time.Now().UnixNano())
	for {
		potentialID := rand.Uint64()
		if !que.ContainsTaskID(potentialID) && potentialID > 0 {
			return potentialID
		}
	}
}

/*
AppendQueue will append the current queue with a new AnalyserTaskRequest.
Will return the taskID. In case an error occurred it will return 0.
The taskID is a uint64.
*/
func (que *AnalyserQueue) AppendQueue(task *AnalyserTaskRequest) uint64 {
	que.mutex.Lock()
	taskID := que.getRandomTaskID()
	var ret uint64 = 0
	if !que.ContainsTaskID(taskID) && task != nil {
		if que.cache.Request(task.crawlerData.addr) {
			que.mutex.Unlock()
			return 0 //kann komische dinge verursachen
		}
		que.cache.Register(task.crawlerData.addr)
		task.cache = que.cache
		task.SetID(taskID)
		task.saver = que.saver
		(*que.getThisQueue())[taskID] = task

		que.Watchdog().Aus()
		ret = taskID

		log.WithFields(
			log.Fields{
				"addr": task.crawlerData.addr,
			}).Info("Appended task")
	}
	que.mutex.Unlock()
	return ret
}

/*
RemoveFromQueue will remove a task from the queue by a given taskID.
*/
func (que *AnalyserQueue) RemoveFromQueue(taskID uint64) bool {
	que.mutex.Lock()
	ret := false
	if taskID > 0 && que.ContainsTaskID(taskID) {
		delete(*que.getThisQueue(), taskID)
		ret = true
	}
	que.mutex.Unlock()
	return ret
}

/*
getAllTaskIds will return the ids of all tasks in a slice.
*/
func (que *AnalyserQueue) getAllTaskIds() []uint64 {
	var ids []uint64
	for k := range *que.getThisQueue() {
		ids = append(ids, k)
	}
	return ids
}

/*
NewQueueStatus  will return a new queueStatus instance containing
information of all states.
*/
func (que *AnalyserQueue) NewQueueStatus() QueueStatus {
	undone, processing, crawlerError, saving, sendToCrawler, finished := que.getNumberOfStatus()
	return QueueStatus{
		undoneTasks:        undone,
		processingTasks:    processing,
		crawlerErrorTasks:  crawlerError,
		savingTasks:        saving,
		sendToCrawlerTasks: sendToCrawler,
		finishedTasks:      finished,
	}
}

/*
getAmountOfUndoneTasks will return the number of all tasks with the status undone.
*/
func (status *QueueStatus) getAmountOfUndoneTasks() uint64 {
	return status.undoneTasks
}

/*
getAmountOfProcessedTasks will return the number of all tasks with the status processing.
*/
func (status *QueueStatus) getAmountOfProcessedTasks() uint64 {
	return status.processingTasks
}

/*
getAmountOfCrawlerErrorTasks will return the number of all tasks with the status crawlerError.
*/
func (status *QueueStatus) getAmountOfCrawlerErrorTasks() uint64 {
	return status.crawlerErrorTasks
}

/*
getAmountOfSavingTasks will return the number of all tasks with the status saving.
*/
func (status *QueueStatus) getAmountOfSavingTasks() uint64 {
	return status.savingTasks
}

/*
getAmountOfSendToCrawlerTasks will return the number of all tasks with the status sendToCrawler.
*/
func (status *QueueStatus) getAmountOfSendToCrawlerTasks() uint64 {
	return status.sendToCrawlerTasks
}

/*
getAmountOfFinishedTasks will return the number of all tasks with the status finished.
*/
func (status *QueueStatus) getAmountOfFinishedTasks() uint64 {
	return status.finishedTasks
}

/*
getNumberOfUndoneTasks will return the number of Tasks of all states.
*/
func (que *AnalyserQueue) getNumberOfStatus() (uint64, uint64, uint64, uint64, uint64, uint64) {
	que.mutex.Lock()
	var undone, processing, crawlerError, saving, sendToCrawler, finished uint64 = 0, 0, 0, 0, 0, 0
	for _, v := range *que.getThisQueue() {
		switch v.State() {
		case UNDONE:
			undone++
		case PROCESSING:
			processing++
		case CRAWLERERROR:
			crawlerError++
		case SAVING:
			saving++
		case SENDTOCRAWLER:
			sendToCrawler++
		case FINISHED:
			finished++
		}
	}
	que.mutex.Unlock()
	return undone, processing, crawlerError, saving, sendToCrawler, finished
}

/*
GetStatusOfQueue will return the status of the queue.
*/
func (que *AnalyserQueue) GetStatusOfQueue() string {
	status := que.NewQueueStatus()
	return fmt.Sprintf("States of the Task currently in the Queue. \nUndone: %d, Processing: %d, CrawlerError: %d, Saving: %d, SendToCrawler: %d, Finished: %d",
		status.getAmountOfUndoneTasks(),
		status.getAmountOfProcessedTasks(),
		status.getAmountOfCrawlerErrorTasks(),
		status.getAmountOfSavingTasks(),
		status.getAmountOfSendToCrawlerTasks(),
		status.getAmountOfFinishedTasks())
}

/*
IsEmpty returns true if queue empty.
*/
func (que *AnalyserQueue) IsEmpty() bool {
	return len(que.Queue) == 0
}

/*
StopQueue will stop the queue.
*/
func (que *AnalyserQueue) StopQueue() {
	que.CleanQueue()
}

/*
CleanQueue will clean the queue.
*/
func (que *AnalyserQueue) CleanQueue() {
	que.deleteAllTasks()
}

/*
deleteAllTasks will empty the queue.
*/
func (que *AnalyserQueue) deleteAllTasks() {
	log.Info("Queue was emptied")
	que.Queue = make(map[uint64]*AnalyserTaskRequest)
}
