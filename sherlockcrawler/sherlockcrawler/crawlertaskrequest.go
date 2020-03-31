package sherlockcrawler

import (
	"fmt"
	"net/http"
)

/*
CrawlerTaskRequest will be a request made by the analyser.
*/
type CrawlerTaskRequest struct {
	addr     string
	isdone   bool
	response *http.Response
}

/*
getAddr getter for the address.
*/
func (creq *CrawlerTaskRequest) getAddr() string {
	return creq.addr
}

/*
getStatus will return the status of the task.
*/
func (creq *CrawlerTaskRequest) isDone() bool {
	return creq.isdone
}

/*
getResponse will return the response of a crawlertask.
*/
func (creq *CrawlerTaskRequest) getResponse() http.Response {
	return *(creq.response)
}

/*
getResponseByReferenz will return the response of a crawlertask.
*/
func (creq *CrawlerTaskRequest) getResponseByReferenz() *http.Response {
	return creq.response
}

/*
setAddr will set the addr to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) setAddr(laddr string) {
	creq.addr = laddr
}

/*
setResponse will set the response of the Request to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) setDone(lisDone bool) {
	creq.isdone = lisDone
}

/*
setResponse will set the response of the Request to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) setResponse(lresponse *http.Response) {
	creq.response = lresponse
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
