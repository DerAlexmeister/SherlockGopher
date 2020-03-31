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
	isDone   bool
	response *http.Response
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
SetResponse will set the response of the Request to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) SetResponse(lresponse *http.Response) {

}
