package sherlockanalyser

import (
	"errors"
	"net/http"
	"strings"
	"time"

	jw "github.com/jwalteri/GO/jwstring"
	model "github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockparser"
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
)

/*
AnalyserTaskRequest will be a request made by the analyser.
*/
type AnalyserTaskRequest struct {
	taskid       uint64 //taskid, send every time.
	addr         string //addr, once
	htmlCode     string
	taskstate    TASKSTATE
	linkLibrary  map[string]string
	rootAddr     string
	foundLinks   []string
	parserTime   int64
	analyserTime int64
	crawlerData  *CrawlerData
}

/*
CrawlerData contains the data send by the crawler.
*/
type CrawlerData struct {
	taskid            uint64
	addr              string
	taskerror         error
	responseHeader    *http.Header
	responseBodyBytes []byte
	statuscode        int
	responseTime      time.Duration
}

/*
NewTask will return an empty AnalyserTaskRequest.
*/
func NewTask(lcrawlerData CrawlerData) AnalyserTaskRequest {
	task := AnalyserTaskRequest{}
	if lcrawlerData.getCTaskError() == nil { //

		task.crawlerData = &lcrawlerData
		task.setAddr(lcrawlerData.addr)
		task.setHTMLCode(string(lcrawlerData.responseBodyBytes))

		task.initialze()
	} else {
		task.setAddr(lcrawlerData.addr)
	}

	return task
}

func (atask *AnalyserTaskRequest) initialze() {
	atask.linkLibrary = make(map[string]string)

	atask.linkLibrary["link"] = "href"
	atask.linkLibrary["script"] = "src"
	atask.linkLibrary["a"] = "href"
	atask.linkLibrary["img"] = "src"
	atask.linkLibrary["form"] = "action"
	atask.linkLibrary["input"] = "value"
	atask.linkLibrary["meta"] = "content"
}

/*
getCrawlerDada will return the id of a given task.
*/
func (atask *AnalyserTaskRequest) getCrawlerDada() *CrawlerData {
	return atask.crawlerData
}

/*
getTaskID will return the id of a given task.
*/
func (atask *AnalyserTaskRequest) getTaskID() uint64 {
	return atask.taskid
}

/*
getAddr getter for the address.
*/
func (atask *AnalyserTaskRequest) getAddr() string {
	return atask.addr
}

/*
getHTMLCode getter for the html code.
*/
func (atask *AnalyserTaskRequest) getHTMLCode() string {
	return atask.htmlCode
}

/*
getRootAddr getter for the root address.
*/
func (atask *AnalyserTaskRequest) getRootAddr() string {
	return atask.rootAddr
}

/*
getFoundLinks getter for the found links.
*/
func (atask *AnalyserTaskRequest) getFoundLinks() []string {
	return atask.foundLinks
}

/*
getTaskState will return the state of the task.
*/
func (atask *AnalyserTaskRequest) getTaskState() TASKSTATE {
	return atask.taskstate
}

/*
setTaskID will set the task id of a given task.
*/
func (atask *AnalyserTaskRequest) setTaskID(lid uint64) {
	atask.taskid = lid
}

/*
setAddr will set the addr to a given AnalyserTaskRequest.
*/
func (atask *AnalyserTaskRequest) setAddr(laddr string) {
	atask.addr = laddr
	rootEnd := jw.OrdinalIndexOf(laddr, "/", 3)

	if rootEnd > 0 {
		atask.rootAddr = laddr[:rootEnd]
	}
}

/*
setHTMLCode will set the html code to a given AnalyserTaskRequest.
*/
func (atask *AnalyserTaskRequest) setHTMLCode(htmlCode string) {
	atask.htmlCode = htmlCode
}

/*
traverse will traverse the tree.
*/
func (atask *AnalyserTaskRequest) traverse(node *model.Node) {
	for _, ele := range node.Children() {
		if len(ele.Children()) > 0 {
			atask.traverse(ele)
		} else {
			atask.classifyNode(ele)
		}
	}
}

/*
classifyNode inspects a node whether its a "link" or not
*/
func (atask *AnalyserTaskRequest) classifyNode(node *model.Node) {
	tag := node.Tag()

	if attributeType, ok := atask.linkLibrary[tag.TagType()]; ok {
		for _, attribute := range tag.Attributes() {
			if attributeType == attribute.AttributeType() {
				// TODO: REMOVE BUGHUNTER
				atask.handleLink(atask.Bughunter(attribute.Value()))
			}
		}
	}
}

/*
TODO: DELETE AFTER FIX
Bughunter kills bugs
*/
func (atask *AnalyserTaskRequest) Bughunter(link string) string {
	link = strings.ReplaceAll(link, "'", "")
	link = strings.ReplaceAll(link, "\"", "")

	return link
}

/*
prettyPrintLink pretty prints a link.
E.g.:
	Input: 	"/stuff/blub"
	Output:	"https://randamonium.bay/stuff/blub"
*/
func (atask *AnalyserTaskRequest) prettyPrintLink(link string) (string, error) {
	if link == "" {
		return "", errors.New("it's not a link")
	}

	// TODO: REMOVE BUGHUNTER
	link = atask.Bughunter(link)

	if link[0] == '/' {
		return atask.getRootAddr() + link, nil
	}

	if !strings.HasPrefix(link, "www") && !strings.HasPrefix(link, "http") {
		return "", errors.New("it's not a link")
	}

	return link, nil
}

/*
handleLink Verifies wheater link is valid (add to NEO4J and send to crawler) or not
*/
func (atask *AnalyserTaskRequest) handleLink(link string) {
	if link, err := atask.prettyPrintLink(link); err == nil {
		if !atask.containedInNEO4J(link) {
			atask.foundLinks = append(atask.foundLinks, link)
			atask.addToNEO4J(link)
		}
	}
}

/*
containedInNEO4J verifies wheater link is already contained in NEO4J or not
*/
func (atask *AnalyserTaskRequest) containedInNEO4J(link string) bool {
	return false
}

/*
addToNEO4J adds a link to NEO4J
*/
func (atask *AnalyserTaskRequest) addToNEO4J(link string) bool {
	return false
}

/*
Execute will search the tree for links and stores the result in the field response of the task
*/
func (atask *AnalyserTaskRequest) Execute() bool {
	atask.taskstate = PROCESSING
	tree := model.NewHTMLTree(atask.getHTMLCode())

	start := time.Now()
	tree.Parse()
	atask.parserTime = time.Since(start).Nanoseconds()

	root := tree.RootNode()

	start = time.Now()
	atask.traverse(root)
	atask.analyserTime = time.Since(start).Nanoseconds()

	atask.taskstate = FINISHED

	return true
}

/*
getCTaskID will return the id of a given task.
*/
func (ctask *CrawlerData) getCTaskID() uint64 {
	return ctask.taskid
}

/*
getCAddr getter for the address.
*/
func (ctask *CrawlerData) getCAddr() string {
	return ctask.addr
}

/*
getCTask will return an error which was caused by the http package.
*/
func (ctask *CrawlerData) getCTaskError() error {
	return ctask.taskerror
}

/*
getCResponseHeader will return the Header of the Response.
*/
func (ctask *CrawlerData) getCResponseHeader() http.Header {
	return *(ctask.responseHeader)
}

/*
getCResponseBody will return the Header of the Response.
*/
func (ctask *CrawlerData) getCResponseBody() []byte {
	return ctask.responseBodyBytes
}

/*
getCStatusCode will return the statuscode.
*/
func (ctask *CrawlerData) getCStatusCode() int {
	return ctask.statuscode
}

/*
getCResponseTime will return the time it took to make the response and get an answer.
*/
func (ctask *CrawlerData) getCResponseTime() time.Duration {
	return ctask.responseTime
}

/*
setCTaskID will set the id of a given task.
*/
func (ctask *CrawlerData) setCTaskID(lid uint64) {
	ctask.taskid = lid
}

/*
setCAddr setter for the address.
*/
func (ctask *CrawlerData) setCAddr(laddr string) {
	ctask.addr = laddr
}

/*
setCTask will set an error which was caused by the http package.
*/
func (ctask *CrawlerData) setCTaskError(lerror error) {
	ctask.taskerror = lerror
}

/*
setCResponseHeader will set the Header of the Response.
*/
func (ctask *CrawlerData) setCResponseHeader(lheader *http.Header) {
	*(ctask.responseHeader) = *lheader
}

/*
setCResponseBody will set the Header of the Response.
*/
func (ctask *CrawlerData) setCResponseBody(lbody []byte) {
	ctask.responseBodyBytes = lbody
}

/*
setCStatusCode will set the statuscode.
*/
func (ctask *CrawlerData) setCStatusCode(lstatuscode int) {
	ctask.statuscode = lstatuscode
}

/*
setCResponseTime will set the time.
*/
func (ctask *CrawlerData) setCResponseTime(ltime time.Duration) {
	ctask.responseTime = ltime
}
