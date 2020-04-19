package sherlockanalyser

import (
	"fmt"
	"math/rand"
	"time"
)

/*
AnalyserQueue will be the queue of the current CrawlerTaskRequest.
*/
type AnalyserQueue struct {
	Queue map[uint64]*analyserTaskRequest
}

/*
QueueStatus struct contains a status of all tasks. E.g undonetaks 10, processingtask 2343, ...
*/
type QueueStatus struct {
	undoneTasks     uint64
	processingTasks uint64
	crawlerErrorTasks uint64
	savingTasks uint64
	sendToCrawlerTasks uint64
	finishedTasks   uint64
}

/*
NewAnalyserQueue will return a new Queue.
*/
func NewAnalyserQueue() AnalyserQueue {
	return AnalyserQueue{
		Queue: map[uint64]*analyserTaskRequest{},
	}
}

/*
getThisQueue will return a pointer to the current Queue.
*/
func (que *AnalyserQueue) getThisQueue() *map[uint64]*analyserTaskRequest {
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
Function to produce a random taskId.
*/
func (que *AnalyserQueue) getRandomUserID() uint64 {
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
Will return the taskId. In case an error occurred it will return 0.
The taskId is a uint64.
*/
func (que *AnalyserQueue) AppendQueue(task *analyserTaskRequest) uint64 {
	taskId := que.getRandomUserID()
	if !que.ContainsTaskID(taskId) && task != nil {
		task.SetId(taskId)
		(*que.getThisQueue())[taskId] = task
		return taskId
	}
	return 0
}

/*
RemoveFromQueue will remove a task from the queue by a given taskId.
*/
func (que *AnalyserQueue) RemoveFromQueue(taskId uint64) bool {
	if taskId > 0 && que.ContainsTaskID(taskId) {
		delete(*que.getThisQueue(), taskId)
		return true
	}
	return false
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
getNumberOfUndoneTasks will return the number of Tasks of all states
*/
func (que *AnalyserQueue) getNumberOfStatus() (uint64, uint64, uint64, uint64, uint64, uint64) {
	var undone, processing, crawlerError, saving, sendToCrawler, finished uint64 = 0, 0, 0, 0, 0, 0
	for _, v := range *que.getThisQueue() {
		switch v.State() {
		case UNDONE: undone++
		case PROCESSING: processing++
		case CRAWLERERROR: crawlerError++
		case SAVING: saving++
		case SENDTOCRAWLER: sendToCrawler++
		case FINISHED: finished++
		}
	}
	return undone, processing, crawlerError, saving, sendToCrawler, finished
}

/*
GetStatusOfQueue will return the
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
IsEmpty returns true if queue empty
*/
func (que *AnalyserQueue) IsEmpty() bool {
	return len(que.Queue) == 0
}
