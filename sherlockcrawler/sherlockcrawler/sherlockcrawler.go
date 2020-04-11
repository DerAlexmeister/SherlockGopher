package sherlockcrawler

import (
	"context"
	"fmt"
	"sync"
	"time"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
	"github.com/pkg/errors"
)

//Time to wait in milliseconds after checking for the tasks again.
const delaytime = 10

/*
Sherlockcrawler will be the Crawlerservice.
*/
type Sherlockcrawler struct {
	Queue        CrawlerQueue //Queue with all tasks
	Dependencies *Sherlockdependencies
}

/*
NewSherlockCrawlerService will return a new sherlockcrawler instance.
*/
func NewSherlockCrawlerService() Sherlockcrawler {
	return Sherlockcrawler{
		Queue: NewCrawlerQueue(),
	}
}

/*
InjectDependency will inject the dependencies for the a sherlockcrawler instance
into the actual instance.
*/
func (sherlock *Sherlockcrawler) InjectDependency(deps *Sherlockdependencies) {
	sherlock.Dependencies = deps
}

/*
getQueue will return a pointer to the queue.
*/
func (sherlock *Sherlockcrawler) getQueue() *CrawlerQueue {
	return &sherlock.Queue
}

/*
CreateTask will append the current queue with a task.
*/
func (sherlock Sherlockcrawler) CreateTask(ctx context.Context, in *proto.CrawlTaskCreateRequest, out *proto.CrawlTaskCreateResponse) error {
	task := NewTask()
	task.setAddr(in.GetUrl())
	if id := (*sherlock.getQueue()).AppendQueue(&task); id > 0 {
		return nil
	}
	out.Statuscode = proto.URL_STATUS_failure //TODO improve struct
	return errors.New("cannot create a task for this queue")
}

/*
manageUndoneTasks will manage all undone tasks and start them.
*/
func (sherlock *Sherlockcrawler) manageUndoneTasks(waitgroup *sync.WaitGroup) {
	var localwaitgroup sync.WaitGroup
	for _, v := range *sherlock.getQueue().getCurrentQueue() {
		if v.getTaskState() == UNDONE {
			go v.MakeRequestAndStoreResponse(&localwaitgroup)
			localwaitgroup.Add(1)
		}
	}
	waitgroup.Done()
}

/*
manageFinishedTasks will managed all finished tasks.
*/
func (sherlock *Sherlockcrawler) manageFinishedTasks(waitgroup *sync.WaitGroup) {
	for _, v := range *sherlock.getQueue().getCurrentQueue() {
		if v.getTaskState() == FINISHED {
			//TODO
		}
	}
	waitgroup.Done()
}

/*
manageFailedTasks manage all tasks which failed.
*/
func (sherlock *Sherlockcrawler) manageFailedTasks(waitgroup *sync.WaitGroup) {
	var localwaitgroup sync.WaitGroup
	for _, v := range *sherlock.getQueue().getCurrentQueue() {
		if v.getTaskState() == FINISHED {
			if v.getTrysError() < 3 {
				go v.MakeRequestAndStoreResponse(&localwaitgroup)
				localwaitgroup.Add(1)
			} else {
				//TODO
			}
		}
	}
	waitgroup.Done()
}

/*
ManageTasks will be a function to check all tasks in a period of time.
Should run as a go routine. Will work like a cronJob.
*/
func (sherlock *Sherlockcrawler) ManageTasks() {
	for {
		var waitgroup sync.WaitGroup
		go sherlock.manageUndoneTasks(&waitgroup)
		go sherlock.manageFinishedTasks(&waitgroup)

		//TODO manage Tasks the right way.
		waitgroup.Wait()
		fmt.Println(sherlock.getQueue().GetStatusOfQueue())
		time.Tick(delaytime * time.Millisecond)
	}
}
