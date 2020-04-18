package sherlockanalyser

import (
	"errors"
	"strings"
	"time"

	jw "github.com/jwalteri/GO/jwstring"
	model "github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockparser/model"
)

/*
TASKSTATE will be a type representing the current TASKSTATE of the task.
*/
type TASKSTATE int

const (
	//UNDONE: A task has not been started yet.
	UNDONE TASKSTATE = 0

	//PROCESSING: A task has been started.
	PROCESSING TASKSTATE = 1

	//SAVING: A task is saving data to NEO4J.
	SAVING TASKSTATE = 2

	//SENDTOCRAWLER: A task is sending data to the crawler.
	SENDTOCRAWLER TASKSTATE = 3

	//FINISHED: A task is finished and ready to be remove from the queue.
	FINISHED TASKSTATE = 4
)

/*
analyserTaskRequest will be a request made by the analyser.
*/
type analyserTaskRequest struct {
	id    		 uint64
	workAddr     string
	html    	 string
	state	     TASKSTATE
	linkTags  	 map[string]string
	rootAddr     string
	foundLinks   []string
	parserTime   int64
	analyserTime int64
	crawlerData  *CrawlerData
}

/*
CrawlerData returns CrawlerData
 */
func (analyserTask *analyserTaskRequest) CrawlerData() *CrawlerData {
	return analyserTask.crawlerData
}

/*
SetCrawlerData sets CrawlerData
*/
func (analyserTask *analyserTaskRequest) SetCrawlerData(crawlerData *CrawlerData) {
	analyserTask.crawlerData = crawlerData
}

/*
AnalyserTime returns AnalyserTime
*/
func (analyserTask *analyserTaskRequest) AnalyserTime() int64 {
	return analyserTask.analyserTime
}

/*
SetAnalyserTime sets AnalyserTime
*/
func (analyserTask *analyserTaskRequest) SetAnalyserTime(analyserTime int64) {
	analyserTask.analyserTime = analyserTime
}

/*
ParserTime returns ParserTime
*/
func (analyserTask *analyserTaskRequest) ParserTime() int64 {
	return analyserTask.parserTime
}

/*
SetParserTime sets ParserTime
*/
func (analyserTask *analyserTaskRequest) SetParserTime(parserTime int64) {
	analyserTask.parserTime = parserTime
}

/*
FoundLinks returns FoundLinks
*/
func (analyserTask *analyserTaskRequest) FoundLinks() []string {
	return analyserTask.foundLinks
}

/*
Links(foundL sets FoundLinks
*/
func (analyserTask *analyserTaskRequest) SetFoundLinks(foundLinks []string) {
	analyserTask.foundLinks = foundLinks
}

/*
RootAddr returns RootAddr
*/
func (analyserTask *analyserTaskRequest) RootAddr() string {
	return analyserTask.rootAddr
}

/*
SetRootAddr sets RootAddr
*/
func (analyserTask *analyserTaskRequest) SetRootAddr(rootAddr string) {
	analyserTask.rootAddr = rootAddr
}

/*
LinkTags returns LinkTags
*/
func (analyserTask *analyserTaskRequest) LinkTags() map[string]string {
	return analyserTask.linkTags
}

/*
SetLinkTags sets LinkTags
*/
func (analyserTask *analyserTaskRequest) SetLinkTags(linkTags map[string]string) {
	analyserTask.linkTags = linkTags
}

/*
State returns State
*/
func (analyserTask *analyserTaskRequest) State() TASKSTATE {
	return analyserTask.state
}

/*
SetState sets State
*/
func (analyserTask *analyserTaskRequest) SetState(state TASKSTATE) {
	analyserTask.state = state
}

/*
Html returns Html
*/
func (analyserTask *analyserTaskRequest) Html() string {
	return analyserTask.html
}

/*
SetHtml sets Html
*/
func (analyserTask *analyserTaskRequest) SetHtml(html string) {
	analyserTask.html = html
}

/*
SetWorkAddr sets the work address.
Identifies and sets the root address.
*/
func (analyserTask *analyserTaskRequest) SetWorkAddr(workAddr string) {
	analyserTask.SetWorkAddr(workAddr)
	rootEnd := jw.OrdinalIndexOf(workAddr, "/", 3)

	if rootEnd > 0 {
		analyserTask.SetRootAddr(workAddr[:rootEnd])
	}
}

/*
WorkAddr returns WorkAddr
*/
func (analyserTask *analyserTaskRequest) WorkAddr() string {
	return analyserTask.workAddr
}


/*
Id returns Id
*/
func (analyserTask *analyserTaskRequest) Id() uint64 {
	return analyserTask.id
}

/*
SetId sets Id
*/
func (analyserTask *analyserTaskRequest) SetId(id uint64) {
	analyserTask.id = id
}

/*
NewTask returns a new task initialized by crawler data
*/
func NewTask(cdata CrawlerData) analyserTaskRequest {
	task := analyserTaskRequest{}

	if cdata.getTaskError() == nil {
		task.SetCrawlerData(&cdata)
		task.SetWorkAddr(cdata.getAddr())
		task.SetHtml(string(cdata.getResponseBody()))

		task.initializeLinkTags()

		task.SetState(UNDONE)
	} else {
		task.SetWorkAddr(cdata.getAddr())
		task.SetState(CRAWLERROR)
	}

	return task
}

func (analyserTask *analyserTaskRequest) initializeLinkTags() {
	analyserTask.linkTags = make(map[string]string)

	analyserTask.linkTags["link"] = "href"
	analyserTask.linkTags["script"] = "src"
	analyserTask.linkTags["a"] = "href"
	analyserTask.linkTags["img"] = "src"
	analyserTask.linkTags["form"] = "action"
	analyserTask.linkTags["input"] = "value"
	analyserTask.linkTags["meta"] = "content"
}

/*
analyze will traverse the tree and classify every node.
*/
func (analyserTask *analyserTaskRequest) analyze(node *model.Node) {
	for _, ele := range node.Children() {
		if len(ele.Children()) > 0 {
			analyserTask.analyze(ele)
		} else {
			analyserTask.classifyNode(ele)
		}
	}
}

/*
classifyNode inspects a node whether its a "link" or not
*/
func (analyserTask *analyserTaskRequest) classifyNode(node *model.Node) {
	tag := node.Tag()

	if attributeType, ok := analyserTask.linkTags[tag.TagType()]; ok {
		for _, attribute := range tag.Attributes() {
			if attributeType == attribute.AttributeType() {
				// TODO: REMOVE bugHunter
				val := analyserTask.bugHunter(attribute.Value())
				analyserTask.SetFoundLinks(append(analyserTask.FoundLinks(), val))
				//analyserTask.handleLink(analyserTask.bugHunter(attribute.Value()))
			}
		}
	}
}

/*
TODO: DELETE AFTER FIX
bugHunter kills bugs
*/
func (analyserTask *analyserTaskRequest) bugHunter(link string) string {
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
func (analyserTask *analyserTaskRequest) prettyPrintLink(link string) (string, error) {
	if link == "" {
		return "", errors.New("it's not a link")
	}

	// TODO: REMOVE bugHunter
	link = analyserTask.bugHunter(link)

	if link[0] == '/' {
		return analyserTask.RootAddr() + link, nil
	}

	if !strings.HasPrefix(link, "www") && !strings.HasPrefix(link, "http") {
		return "", errors.New("it's not a link")
	}

	return link, nil
}

/*
handleLink Verifies whether link is valid (add to NEO4J and send to crawler) or not
*/
func (analyserTask *analyserTaskRequest) handleLink(link string) {
	if link, err := analyserTask.prettyPrintLink(link); err == nil {
		if !analyserTask.containedInNEO4J(link) {
			analyserTask.foundLinks = append(analyserTask.foundLinks, link)
			analyserTask.addToNEO4J(link)
		}
	}
}

/*
containedInNEO4J verifies whether link is already contained in NEO4J or not
*/
func (analyserTask *analyserTaskRequest) containedInNEO4J(link string) bool {
	return false
}

/*
addToNEO4J adds a link to NEO4J
*/
func (analyserTask *analyserTaskRequest) addToNEO4J(link string) bool {
	return false
}

/*
Execute will search the tree for links and stores the result in the field response of the task
*/
func (analyserTask *analyserTaskRequest) Execute() bool {
	analyserTask.taskstate = PROCESSING
	tree := model.NewHTMLTree(analyserTask.gethtml())

	start := time.Now()
	tree.Parse()
	analyserTask.parserTime = time.Since(start).Nanoseconds()

	root := tree.RootNode()

	start = time.Now()
	analyserTask.analyze(root)
	analyserTask.analyserTime = time.Since(start).Nanoseconds()

	analyserTask.taskstate = FINISHED

	return true
}
