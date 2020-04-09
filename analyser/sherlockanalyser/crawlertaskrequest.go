package sherlockanalyser

import (
	"net/http"
	"time"

	model "github.com/ob-algdatii-20ss/SherlockGopher/analyser/html2treeparser"
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
AnalyserTaskRequest will be a request made by the analyser.
*/
type AnalyserTaskRequest struct {
	taskid            int64  //taskid, send every time.
	addr              string //addr, once
	taskstate         TASKSTATE
	taskerror         error //error, send as string incase there is an error then dont send a body
	response          *http.Response
	responseHeader    *http.Header //header, once (typ map)
	responseBody      []string
	responseBodyBytes []byte        //body, split
	statuscode        int           //statuscode, once
	responseTime      time.Duration //response time, once
}

/*
NewTask will return an empty AnalyserTaskRequest.
*/
func NewTask() AnalyserTaskRequest {
	return AnalyserTaskRequest{}
}

/*
getTaskID will return the id of a given task.
*/
func (creq *AnalyserTaskRequest) getTaskID() int64 {
	return creq.taskid
}

/*
getAddr getter for the address.
*/
func (creq *AnalyserTaskRequest) getAddr() string {
	return creq.addr
}

/*
getTASKSTATE will return the TASKSTATE of the task.
*/
func (creq *AnalyserTaskRequest) getTaskState() TASKSTATE {
	return creq.taskstate
}

/*lTASKSTATE
getResponse will return the response of a Analysertask.
*/
func (creq *AnalyserTaskRequest) getResponse() http.Response {
	return *(creq.response)
}

/*
getResponseHeader will return the Header of the Response.
*/
func (creq *AnalyserTaskRequest) getResponseHeader() http.Header {
	return *(creq.responseHeader)
}

/*
getResponseBody will return the Header of the Response.
*/
func (creq *AnalyserTaskRequest) getResponseBody() []string {
	return creq.responseBody
}

/*
getResponseByReferenz will return the response of a Analysertask.
*/
func (creq *AnalyserTaskRequest) getResponseByReferenz() *http.Response {
	return creq.response
}

/*
getResponseBodyInBytes will return the responsebody as a bytearray.
*/
func (creq *AnalyserTaskRequest) getResponseBodyInBytes() []byte {
	return creq.responseBodyBytes
}

/*
setTaskID will set the task id of a given task.
*/
func (creq *AnalyserTaskRequest) setTaskID(lid int64) {
	creq.taskid = lid
}

/*
setAddr will set the addr to a given AnalyserTaskRequest.
*/
func (creq *AnalyserTaskRequest) setAddr(laddr string) {
	creq.addr = laddr
}

/*
setResponse will set the response of the Request to a given AnalyserTaskRequest.
*/
func (creq *AnalyserTaskRequest) setDone(ltaskstate TASKSTATE) {
	creq.taskstate = ltaskstate
}

/*
setResponse will set the response of the Request to a given AnalyserTaskRequest.
*/
func (creq *AnalyserTaskRequest) setResponse(lresponse *http.Response) {
	creq.response = lresponse
}

/*
setResponseBody will set the responsebody of the Response to a given AnalyserTaskRequest.
*/
func (creq *AnalyserTaskRequest) setResponseBody(lbody []string) {
	creq.responseBody = lbody
}

/*
setResponseHeader will set the responseheader of the Response to a given AnalyserTaskRequest.
*/
func (creq *AnalyserTaskRequest) setResponseHeader(lheader *http.Header) {
	creq.responseHeader = lheader
}

/*
setResponseBodyInBytes will set the responsebody in bytes of the Response to a given AnalyserTaskRequest.
*/
func (creq *AnalyserTaskRequest) setResponseBodyInBytes(lbody []byte) {
	creq.responseBodyBytes = lbody
}

/*
SearchForLinks will search the html for links and returns them.
*/
func (creq *AnalyserTaskRequest) SearchForLinks() ([]string, error) {

	return nil, nil
}

/*
Traverse will traverse the tree.
*/
func (creq *AnalyserTaskRequest) Traverse(node *model.Node) {

}

/*
MakeRequestAndStoreResponse will search the tree for links and stores the result in the field response of the task
*/
func (creq *AnalyserTaskRequest) MakeRequestAndStoreResponse() bool {
	return true
}
