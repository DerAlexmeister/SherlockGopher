package sherlockscreenshot

import (
	"sync"
)

type Screenshot struct {
	Picture []byte
	URL     string
}

/*
ScreenshotServiceHandler will be ScreenshotService representation.
*/
type ScreenshotServiceHandler struct {
	ScreenshotQueue *ScreenshotQueue
}

/*
getQueue will return a pointer to the queue.
*/
func (screenshot *ScreenshotServiceHandler) getQueue() *ScreenshotQueue {
	return screenshot.ScreenshotQueue
}

/*
manageUndoneTasks will manage all undone tasks and start them.
*/
func (screenshot *ScreenshotServiceHandler) manageUndoneTasks() {
	var localWaitGroup sync.WaitGroup

	for _, v := range *(*screenshot.getQueue()).getThisQueue() {
		if v.getState() == UNDONE {
			v.SetState(PROCESSING)
			localWaitGroup.Add(1)
			go v.Execute(&localWaitGroup)
		}
	}
	localWaitGroup.Wait()
}

/*
manageFinishedTasks will managed all finished tasks.
*/
func (screenshot *ScreenshotServiceHandler) manageFinishedTasks() {
	for k, v := range *(*screenshot.getQueue()).getThisQueue() {
		if v.getState() == FINISHED {
			(*screenshot.getQueue()).RemoveFromQueue(k)
		}
	}
}

/*
runManager will run all tasks related functions on the analyserServiceHandler and put them into a go routine.
*/
func (screenshot *ScreenshotServiceHandler) runManager() {
	screenshot.manageUndoneTasks()
	screenshot.manageFinishedTasks()
}
