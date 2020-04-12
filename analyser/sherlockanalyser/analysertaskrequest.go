package sherlockanalyser

import (
	"errors"
	"strings"

	model "github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockparser/model"
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
	taskid      int64  //taskid, send every time.
	addr        string //addr, once
	htmlCode    string
	taskstate   TASKSTATE
	linkLibrary map[string]string
	rootAddr    string
	foundLinks  []string
}

/*
NewTask will return an empty AnalyserTaskRequest.
*/
func NewTask() AnalyserTaskRequest {
	task := AnalyserTaskRequest{}
	task.taskstate = UNDONE
	task.initialze()

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
getTaskID will return the id of a given task.
*/
func (atask *AnalyserTaskRequest) getTaskID() int64 {
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
func (atask *AnalyserTaskRequest) setTaskID(lid int64) {
	atask.taskid = lid
}

/*
setAddr will set the addr to a given AnalyserTaskRequest.
*/
func (atask *AnalyserTaskRequest) setAddr(laddr string) {
	atask.addr = laddr
	rootEnd := atask.ordinalIndexOf(laddr, "/", 3)
	atask.rootAddr = laddr[:rootEnd]
}

/*
setHTMLCode will set the html code to a given AnalyserTaskRequest.
*/
func (atask *AnalyserTaskRequest) setHTMLCode(htmlCode string) {
	atask.htmlCode = htmlCode
}

/*
See jwalteri go
*/
func (atask *AnalyserTaskRequest) ordinalIndexOf(str string, searchStr string, ordinal int) int {
	if str == "" {
		return -2
	}

	if searchStr == "" {
		return -3
	}

	if ordinal <= 0 {
		return -4
	}

	ret := -1
	pos := strings.Index(str, searchStr)

	if pos == -1 {
		return -1
	}

	for ordinal > 0 && pos != -1 {
		str = str[pos+1:]
		ret += pos + 1
		ordinal--

		pos = strings.Index(str, searchStr)
	}

	if ordinal != 0 {
		return -1
	}

	return ret
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
	tree.Parse()
	root := tree.RootNode()
	atask.traverse(root)
	atask.taskstate = FINISHED

	return true
}
