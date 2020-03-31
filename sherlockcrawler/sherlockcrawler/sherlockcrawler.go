package sherlockcrawler

import (
	"fmt"
	"net/http"
)

/*
CrawlerTaskRequest will be a request made by the analyser.
*/
type CrawlerTaskRequest struct {
	addr  string
	state bool
}

/*
CrawlerQueue will be the queue of the current CrawlerTaskRequest.
*/
type CrawlerQueue struct {
	Queue map[string]*CrawlerTaskRequest
}

/*
getCurrentQueue will return a pointer to the current Queue.
*/
func (que *CrawlerQueue) getCurrentQueue() *(map[string]*CrawlerTaskRequest) {
	return &que.Queue
}

/*
MakeRequestForHTML will make a request to a given Website and return its HTML-Code.
*/
func (creq *CrawlerTaskRequest) MakeRequestForHTML() (*http.Response, error) {
	response, err := http.Get(creq.addr)
	if err != nil {
		return nil, fmt.Errorf("An error occured while trying to get the Website: %s", creq.addr)
	}
	return response, nil
}

/*
ContainsAddress will check whether or not a addr is allready in use or not.
*/
func (que *CrawlerQueue) ContainsAddress(addr string) bool {
	_, contains := (*que.getCurrentQueue())[addr]
	if !contains {
		return false
	}
	return true
}

/*
AppendQueue will append the current queue with a new CrawlerTaskRequest.
*/
func (que *CrawlerQueue) AppendQueue(target string, task *CrawlerTaskRequest) error {
	if !que.ContainsAddress(target) {
		(*que.getCurrentQueue())[target] = task
		return nil //TODO Returntype noch mal überarbeiten weil Error vlt nicht das beste.
	}
	return nil //TODO Returntype noch mal überarbeiten weil Error vlt nicht das beste.
}
