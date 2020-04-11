package sherlockcrawler

import (
	"context"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
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
CreateTask will append the current queue with a task.
*/
func (sherlock *Sherlockcrawler) CreateTask(ctx context.Context, in *proto.CrawlTaskCreateRequest, out *proto.CrawlTaskCreateResponse) error {
	return nil
}
