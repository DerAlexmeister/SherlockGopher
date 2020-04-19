package sherlockanalyser

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"

	jw "github.com/jwalteri/GO/jwstring"
	model "github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockparser"
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

	//CRAWLERERROR: A task contains a crawler error.
	CRAWLERERROR TASKSTATE = 2

	//SAVING: A task is saving data to NEO4J.
	SAVING TASKSTATE = 3

	//SENDTOCRAWLER: A task is sending data to the crawler.
	SENDTOCRAWLER TASKSTATE = 4

	//FINISHED: A task is finished and ready to be remove from the queue.
	FINISHED TASKSTATE = 5
)

/*
analyserTaskRequest will be a request made by the analyser.
*/
type analyserTaskRequest struct {
	id           uint64
	workAddr     string
	html         string
	state        TASKSTATE
	linkTags     map[string]string
	rootAddr     string
	foundLinks   []string
	parserTime   int64
	analyserTime int64
	crawlerData  *CrawlerData
	Dependencies *AnalyserDependency
}

/*
InjectDependency will inject the dependencies for the a analyser instance
into the actual instance.
*/
func (analyserTask *analyserTaskRequest) InjectDependency(deps *AnalyserDependency) {
	analyserTask.Dependencies = deps
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
	analyserTask.workAddr = workAddr
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
	task.SetCrawlerData(&cdata)

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
				analyserTask.handleLink(val)
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
verifyLink pretty prints a link.
E.g.:
	Input: 	"/stuff/blub"
	Output:	"https://randamonium.bay/stuff/blub"
*/
func (analyserTask *analyserTaskRequest) verifyLink(link string) (string, error) {
	if link == "" {
		return "", errors.New("it's not a link")
	}

	if link[0] == '/' {
		return analyserTask.RootAddr() + link, nil
	} else if !strings.HasPrefix(link, "www") && !strings.HasPrefix(link, "http") {
		return "", errors.New("it's not a link")
	} else {
		return link, nil
	}
}

/*
handleLink Verifies whether link is valid (add to NEO4J and send to crawler) or not
*/
func (analyserTask *analyserTaskRequest) handleLink(link string) {
	if link, err := analyserTask.verifyLink(link); err == nil {
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
func (analyserTask *analyserTaskRequest) Execute(waitGroup *sync.WaitGroup) bool {
	analyserTask.preprocess()

	analyserTask.SetState(FINISHED)

	if waitGroup != nil {
		defer waitGroup.Done()
	}
	return true
}

func (analyserTask *analyserTaskRequest) preprocess() {
	if analyserTask.CrawlerData().getTaskError() == nil {
		analyserTask.process()
	} else {
		analyserTask.handleCrawlerError()
	}
}

func (analyserTask *analyserTaskRequest) handleCrawlerError() {
	analyserTask.SetState(CRAWLERERROR)
	//TODO: ADD NEO4J
	analyserTask.SetState(FINISHED)
}

func (analyserTask *analyserTaskRequest) saveToNeo4J() {
	analyserTask.SetState(SAVING)
	//TODO: ADD NEO4J
}

func (analyserTask *analyserTaskRequest) sendToCrawler() {
	analyserTask.SetState(SENDTOCRAWLER)

	serv := analyserTask.Dependencies.Crawler()

	for _, link := range analyserTask.FoundLinks() {
		message := &crawlerproto.CrawlTaskCreateRequest{
			Url: link,
		}

		_, err := serv.CreateTask(context.TODO(), message)
		if err != nil {
			fmt.Println("Error while sending link to crawler")
			fmt.Println(err)
		}
	}
}

func (analyserTask *analyserTaskRequest) process() {
	analyserTask.SetState(PROCESSING)

	analyserTask.SetWorkAddr(analyserTask.CrawlerData().getAddr())
	analyserTask.SetHtml(string(analyserTask.CrawlerData().responseBodyBytes))
	analyserTask.initializeLinkTags()

	htmlTree := model.NewHTMLTree(analyserTask.Html())
	start := time.Now()
	htmlTree.Parse(false)
	rootNode := htmlTree.RootNode()
	analyserTask.parserTime = time.Since(start).Nanoseconds()

	start = time.Now()
	analyserTask.analyze(rootNode)
	analyserTask.analyserTime = time.Since(start).Nanoseconds()

	analyserTask.saveToNeo4J()
	analyserTask.sendToCrawler()
	analyserTask.SetState(FINISHED)
}
