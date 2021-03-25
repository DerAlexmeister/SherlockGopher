package sherlockanalyser

import (
	"net/http"
	"time"
)

/*
CrawlerData contains the data send by the crawler.
*/
type CrawlerData struct {
	taskID            uint64
	addr              string
	taskError         error
	responseHeader    *http.Header
	responseBodyBytes []byte
	statusCode        int
	responseTime      time.Duration
}

func NewCrawlerData() *CrawlerData {
	data := CrawlerData{}
	return &data
}

/*
getTaskID will return the id of a given task.
*/
func (cdata *CrawlerData) getTaskID() uint64 {
	return cdata.taskID
}

/*
getAddr getter for the address.
*/
func (cdata *CrawlerData) getAddr() string {
	return cdata.addr
}

/*
getTaskError will return an error which was caused by the http package.
*/
func (cdata *CrawlerData) getTaskError() error {
	return cdata.taskError
}

/*
getResponseHeader will return the Header of the Response.
*/
func (cdata *CrawlerData) getResponseHeader() http.Header {
	return *(cdata.responseHeader)
}

/*
getResponseBody will return the Header of the Response.
*/
func (cdata *CrawlerData) getResponseBody() []byte {
	return cdata.responseBodyBytes
}

/*
getStatusCode will return the status code.
*/
func (cdata *CrawlerData) getStatusCode() int {
	return cdata.statusCode
}

/*
getResponseTime will return the time it took to make the response and get an answer.
*/
func (cdata *CrawlerData) getResponseTime() time.Duration {
	return cdata.responseTime
}

/*
setTaskID will set the id of a given task.
*/
func (cdata *CrawlerData) setTaskID(lid uint64) {
	cdata.taskID = lid
}

/*
setAddr setter for the address.
*/
func (cdata *CrawlerData) setAddr(addr string) {
	cdata.addr = addr
}

/*
setTaskError will set an error which was caused by the http package.
*/
func (cdata *CrawlerData) setTaskError(err error) {
	cdata.taskError = err
}

/*
setResponseHeader will set the Header of the Response.
*/
func (cdata *CrawlerData) setResponseHeader(header *http.Header) {
	cdata.responseHeader = header
}

/*
setResponseBody will set the Header of the Response.
*/
func (cdata *CrawlerData) setResponseBody(body []byte) {
	cdata.responseBodyBytes = body
}

/*
setStatusCode will set the statusCode.
*/
func (cdata *CrawlerData) setStatusCode(statusCode int) {
	cdata.statusCode = statusCode
}

/*
setResponseTime will set the time.
*/
func (cdata *CrawlerData) setResponseTime(time time.Duration) {
	cdata.responseTime = time
}
