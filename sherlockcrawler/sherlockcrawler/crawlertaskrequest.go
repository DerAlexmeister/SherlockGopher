package sherlockcrawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*
TASKSTATE will be a type representing the current TASKSTATE of the task.
*/
type TASKSTATE int

const (
	//UNDONE will be a task untouch.
	UNDONE TASKSTATE = 0
	//PROCESSING will be a task currently working on.
	PROCESSING TASKSTATE = 1
	//FINISHED will be a task which is successfully completed.
	FINISHED TASKSTATE = 2
	//FAILED is a task which was in the state of PROCESSING but failed to complet.
	FAILED TASKSTATE = 3
)

/*
CrawlerTaskRequest will be a request made by the analyser.
*/
type CrawlerTaskRequest struct {
	taskid            int64
	addr              string
	taskstate         TASKSTATE
	taskerror         error
	response          *http.Response
	responseHeader    *http.Header
	responseBody      string
	responseBodyBytes []byte
	statuscode        int
}

/*
NewTask will return an empty CrawlerTaskRequest.
*/
func NewTask() CrawlerTaskRequest {
	return CrawlerTaskRequest{}
}

/*
getTaskID will return the id of a given task.
*/
func (creq *CrawlerTaskRequest) getTaskID() int64 {
	return creq.taskid
}

/*
getAddr getter for the address.
*/
func (creq *CrawlerTaskRequest) getAddr() string {
	return creq.addr
}

/*
getTASKSTATE will return the TASKSTATE of the task.
*/
func (creq *CrawlerTaskRequest) getTaskState() TASKSTATE {
	return creq.taskstate
}

/*lTASKSTATE
getResponse will return the response of a crawlertask.
*/
func (creq *CrawlerTaskRequest) getResponse() http.Response {
	return *(creq.response)
}

/*
getResponseHeader will return the Header of the Response.
*/
func (creq *CrawlerTaskRequest) getResponseHeader() http.Header {
	return *(creq.responseHeader)
}

/*
getResponseBody will return the Header of the Response.
*/
func (creq *CrawlerTaskRequest) getResponseBody() string {
	return creq.responseBody
}

/*
getResponseByReferenz will return the response of a crawlertask.
*/
func (creq *CrawlerTaskRequest) getResponseByReferenz() *http.Response {
	return creq.response
}

/*
getResponseBodyInBytes will return the responsebody as a bytearray.
*/
func (creq *CrawlerTaskRequest) getResponseBodyInBytes() []byte {
	return creq.responseBodyBytes
}

/*
getTask will return an error which was caused by the http package.
*/
func (creq *CrawlerTaskRequest) getTaskError() error {
	return creq.taskerror
}

/*
getStatusCode will return the statuscode.
*/
func (creq *CrawlerTaskRequest) getStatusCode() int {
	return creq.statuscode
}

/*
setTaskID will set the task id of a given task.
*/
func (creq *CrawlerTaskRequest) setTaskID(lid int64) {
	creq.taskid = lid
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
func (creq *CrawlerTaskRequest) setDone(ltaskstate TASKSTATE) {
	creq.taskstate = ltaskstate
}

/*
setResponse will set the response of the Request to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) setResponse(lresponse *http.Response) {
	creq.response = lresponse
}

/*
setResponseBody will set the responsebody of the Response to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) setResponseBody(lbody string) {
	creq.responseBody = lbody
}

/*
setResponseHeader will set the responseheader of the Response to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) setResponseHeader(lheader *http.Header) {
	creq.responseHeader = lheader
}

/*
setResponseBodyInBytes will set the responsebody in bytes of the Response to a given CrawlerTaskRequest.
*/
func (creq *CrawlerTaskRequest) setResponseBodyInBytes(lbody []byte) {
	creq.responseBodyBytes = lbody
}

/*
setTaskError will set the error raised by the http package.
*/
func (creq *CrawlerTaskRequest) setTaskError(lerror error) {
	creq.taskerror = lerror
}

/*
setStatusCode set the statuscode of the task.
*/
func (creq *CrawlerTaskRequest) setStatusCode(lstatuscode int) {
	creq.statuscode = lstatuscode
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
MakeRequestAndStoreResponse will make a request and store the result in the field response of the task.
*/
func (creq *CrawlerTaskRequest) MakeRequestAndStoreResponse() bool {
	response, err := creq.MakeRequestForHTML()
	if err != nil {
		log.Fatal(err) //TODO formated error
		creq.setTaskError(err)
		return false
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err) //TODO formated error
		creq.setTaskError(err)
		return false
	}
	//TODO Check wo alles via setter gesetzt wird.
	creq.setResponse(response)
	creq.setResponseBody(string(bodyBytes))
	creq.setResponseBodyInBytes(bodyBytes)
	return true
}
