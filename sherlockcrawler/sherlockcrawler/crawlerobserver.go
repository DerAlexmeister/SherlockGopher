package sherlockcrawler

import (
	log "github.com/sirupsen/logrus"
)

//TODO JW comments
//SystemState stores the current state of the crawler service (STOP, PAUSE, RUNNING, IDLE).
type SystemState int

const (
	//STOP System is shut down. Kill and remove all queue tasks. No further actions.
	STOP = 0

	//PAUSE System is paused. Pause task execution.
	PAUSE = 1

	//RUNNING System is running. Start all tasks.
	RUNNING = 2

	//IDLE System is idle. Waiting for tasks with low performance.
	IDLE = 4
)

/*
CrawlerObserverInterface is the interface of the CrawlerMediator.
*/
type CrawlerObserverInterface interface {
	UpdateState(state int)
	CleanCommand()
}

/*
CleanCommand cleans the crawler queue.
*/
func (obs CrawlerObserver) CleanCommand() {
	log.Info("Observer called CleanQueue")
	obs.queue.CleanQueue()
	obs.streamerQue.CleanQueue()
}

/*
UpdateState setter for the state of the service/queue.
*/
func (obs CrawlerObserver) UpdateState(state int) {
	log.Info("Observer updated state")
	if !(state > 4 || state < 0) {
		if obs.Crawler() != nil {
			obs.crawler.SetCrawlerState(state)
		}
		if obs.Queue() != nil {
			obs.queue.SetState(state)
		}
		if obs.StreamerQue() != nil {
			obs.streamerQue.SetState(state)
		}
	} else {
		log.Warn("Observer got an unknown state ", state)
	}
}

/*
CrawlerObserver mediates important values of the crawler system.
*/
type CrawlerObserver struct {
	crawler     *SherlockCrawler
	queue       CrawlerQueueInterface
	streamerQue CrawlerQueueInterface
}

/*
StreamerQue returns the streamerQue of the crawler.
*/
func (obs *CrawlerObserver) StreamerQue() CrawlerQueueInterface {
	return obs.streamerQue
}

/*
SetStreamerQue is a setter for the streamerQue of the crawler.
*/
func (obs *CrawlerObserver) SetStreamerQue(streamerQue CrawlerQueueInterface) {
	obs.streamerQue = streamerQue
	streamerQue.SetObserver(obs)
}

/*
Queue returns the queue of the crawler.
*/
func (obs *CrawlerObserver) Queue() CrawlerQueueInterface {
	return obs.queue
}

/*
SetQueue is a setter for the queue of the crawler.
*/
func (obs *CrawlerObserver) SetQueue(queue CrawlerQueueInterface) {
	obs.queue = queue
	queue.SetObserver(obs)
}

/*
Crawler returns the crawler stored in the observer.
*/
func (obs *CrawlerObserver) Crawler() *SherlockCrawler {
	return obs.crawler

}

/*
SetCrawler is a setter for the crawler stored in the observer.
*/
func (obs *CrawlerObserver) SetCrawler(crawler *SherlockCrawler) {
	obs.crawler = crawler
	crawler.SetObserver(obs)
}

/*
NewCrawlerObserver creates a new crawler oberver and returns it.
*/
func NewCrawlerObserver() CrawlerObserver {
	return CrawlerObserver{}
}
