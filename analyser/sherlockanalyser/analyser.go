package sherlockanalyser

import (
	"context"
	"fmt"
	"sync"
	"time"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto"
)

//Time to wait in milliseconds after checking for the tasks again.
const delaytime = 20

/*
AnalyserServiceHandler will be AnalyserService representation.
*/
type AnalyserServiceHandler struct {
	AnalyserQueue *AnalyserQueue
	Dependencies  *AnalyserDependency
}

/*
InjectDependency will inject the dependencies for the a analyser instance
into the actual instance.
*/
func (analyser AnalyserServiceHandler) InjectDependency(deps *AnalyserDependency) {
	analyser.Dependencies = deps
}

/*
NewAnalyserServiceHandler will return an new AnalyserServiceHandler instance.
*/
func NewAnalyserServiceHandler() *AnalyserServiceHandler {
	que := NewAnalyserQueue()
	return &AnalyserServiceHandler{
		AnalyserQueue: &que,
		Dependencies:  nil,
	}
}

/*
getDependency will return a pointer to the dependencies instance of this service.
*/
func (analyser *AnalyserServiceHandler) getDependency() *AnalyserDependency {
	return analyser.Dependencies
}

/*
getQueue will return a pointer to the queue.
*/
func (analyser *AnalyserServiceHandler) getQueue() *AnalyserQueue {
	return analyser.AnalyserQueue
}

/*
manageUndoneTasks will manage all undone tasks and start them.
*/
func (analyser *AnalyserServiceHandler) manageUndoneTasks(waitGroup *sync.WaitGroup) {
	var localWaitGroup sync.WaitGroup
	for _, v := range *(*analyser.getQueue()).getThisQueue() {
		if v.State() == UNDONE {
			localWaitGroup.Add(1)
			go v.Execute(&localWaitGroup)
		}
	}
	localWaitGroup.Wait()
	defer waitGroup.Done()
}

/*
manageFinishedTasks will managed all finished tasks.
*/
func (analyser *AnalyserServiceHandler) manageFinishedTasks(waitGroup *sync.WaitGroup) {
	for k, v := range *(*analyser.getQueue()).getThisQueue() {
		if v.State() == FINISHED {
			(*analyser.getQueue()).RemoveFromQueue(k)
		}
	}
	defer waitGroup.Done()
}

/*
runManager will run all tasks related functions on the analyserServiceHandler and put them into a go routine.
*/
func (analyser *AnalyserServiceHandler) runManager() {
	var localWaitGroup sync.WaitGroup
	localWaitGroup.Add(1)
	go analyser.manageUndoneTasks(&localWaitGroup)
	localWaitGroup.Wait()
	localWaitGroup.Add(1)
	go analyser.manageFinishedTasks(&localWaitGroup)
	localWaitGroup.Wait()
	localWaitGroup.Add(1)
}

/*
ManageTasks will be a function to check all tasks in a period of time.
Should run as a go routine. Will work like a cronJob.
*/
func (analyser *AnalyserServiceHandler) ManageTasks() { //TODO kill loop on signal
	for {
		analyser.runManager()
		fmt.Println(analyser.AnalyserQueue.GetStatusOfQueue())
		time.Tick(delaytime * time.Millisecond)
	}
}

/*
StatusOfTaskQueue will send the status of the queue.
*/
func (analyser *AnalyserServiceHandler) StatusOfTaskQueue(ctx context.Context, _ *proto.TaskStatusRequest, out *proto.TaskStatusResponse) error {
	undone, processing, crawlerError, saving, sendToCrawler, finished := analyser.getQueue().getNumberOfStatus()
	out.Undone = undone
	out.Processing = processing
	out.CrawlerError = crawlerError
	out.Saving = saving
	out.SendToCrawler = sendToCrawler
	out.Finished = finished
	return nil
}
