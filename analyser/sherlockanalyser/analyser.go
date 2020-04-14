package sherlockanalyser

import (
	"context"
	"fmt"
	"sync"
	"time"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto/analyser"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
)

/*
AnalyserServiceHandler will be AnalyserService representation.
*/
type AnalyserServiceHandler struct {
	AnalyserQueue *AnalyserQueue
	Dependencies  *AnalyserDependency
}

/*
NewAnalyserServiceHandler will return an new AnalyserServiceHandler instance.
*/
func NewAnalyserServiceHandler() *AnalyserServiceHandler {
	return &AnalyserServiceHandler{}
}

/*
InjectDependency will inject the dependencies for the a sherlockcrawler instance
into the actual instance.
*/
func (analyser AnalyserServiceHandler) InjectDependency(deps *AnalyserDependency) {
	analyser.Dependencies = deps
}

/*
SendResult will send the result to the crawler.
*/
func (analyser AnalyserServiceHandler) SendResult() {

}

/*
manageUndoneTasks manage all tasks which are undone.
*/
func (analyser *AnalyserServiceHandler) manageUndoneTasks(waitgroup *sync.WaitGroup) {
	var localwaitgroup sync.WaitGroup

	for _, v := range *analyser.AnalyserQueue.getCurrentQueue() {
		if v.getTaskState() == FINISHED {
			go v.Execute()
			localwaitgroup.Add(1)
		}
	}

	waitgroup.Done()
}

//Time to wait in milliseconds after checking for the tasks again.
const delaytime = 10

/*
runManager will run all tasks related functions on the sherlockcrawlerservice and put them into a go routine.
*/
func (analyser *AnalyserServiceHandler) runManager() {
	var waitgroup sync.WaitGroup
	go analyser.manageUndoneTasks(&waitgroup)
	go analyser.manageFinishedTasks(&waitgroup)
	waitgroup.Wait()
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
manageFinishedTasks will managed all finished tasks.
*/
func (analyser *AnalyserServiceHandler) manageFinishedTasks(waitgroup *sync.WaitGroup) {
	for _, v := range *analyser.AnalyserQueue.getCurrentQueue() {
		if v.getTaskState() == FINISHED {
			analyser.SendResultToCrawler(v)
		}
	}
	waitgroup.Done()
}

/*
SendResultToCrawler sends the found links to the crawler.
*/
func (analyser *AnalyserServiceHandler) SendResultToCrawler(task *AnalyserTaskRequest) {
	serv := analyser.Dependencies.Crawler()

	for _, link := range task.foundLinks {
		message := &crawlerproto.CrawlTaskCreateRequest{
			Url: link,
		}

		res, err := serv.CreateTask(context.TODO(), message)
		if err != nil {
			fmt.Println(err)
		}

		if res.Statuscode == crawlerproto.URL_STATUS_ok {
			analyser.AnalyserQueue.RemoveFromQueue(task.getTaskID())
		}

	}
}

/*
StatusOfTaskQueue will send the status of the queue.
*/
func (analyser *AnalyserServiceHandler) StatusOfTaskQueue(ctx context.Context, _ *proto.TaskStatusRequest, out *proto.TaskStatusResponse) error {
	undone, processing, finished := analyser.AnalyserQueue.getNumberOfStatus()
	out.Undone = undone
	out.Processing = processing
	out.Finished = finished
	return nil
}
