package sherlockcrawler

import (
	"context"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
	"github.com/pkg/errors"
)

/*
Sherlockcrawler will be the Crawlerservice.
*/
type Sherlockcrawler struct {
	Queue        CrawlerQueue //Queue with all tasks
	Dependencies *Sherlockdependencies
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
func (sherlock *Sherlockcrawler) CreateTask(ctx context.Context, in *proto.CrawlTaskCreateRequest, out *proto.CrawlTaskCreateResponse) error {
	task := NewTask()
	task.setAddr(in.GetUrl())
	if id := (*sherlock.getQueue()).AppendQueue(&task); id > 0 {
		return nil
	}
	out.Statuscode = proto.URL_STATUS_failure //TODO
	return errors.New("cannot create a task for this queue")
}
