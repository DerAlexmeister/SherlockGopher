package sherlockanalyser

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"strings"
	"sync"
	"time"

	neo "github.com/DerAlexx/SherlockGopher/sherlockneo"
	log "github.com/sirupsen/logrus"

	crawlerproto "github.com/DerAlexx/SherlockGopher/sherlockcrawler/proto"

	model "github.com/DerAlexx/SherlockGopher/analyser/sherlockparser"
	jp "github.com/jpillora/go-tld"
)

/*
TaskState will be a type representing the current TaskState of the task.
*/
type TaskState int

const (
	//UNDONE A task has not been started yet.
	UNDONE TaskState = 0

	//PROCESSING A task has been started.
	PROCESSING TaskState = 1

	//CRAWLERERROR A task contains a crawler error.
	CRAWLERERROR TaskState = 2

	//SAVING A task is saving data to NEO4J.
	SAVING TaskState = 3

	//SENDTOCRAWLER A task is sending data to the crawler.
	SENDTOCRAWLER TaskState = 4

	//FINISHED A task is finished and ready to be remove from the queue.
	FINISHED TaskState = 5

	href string = "href"
)

/*
AnalyserTaskRequest will be a request made by the analyser.
*/
type AnalyserTaskRequest struct {
	id           uint64
	workAddr     string
	domainInfo   *jp.URL
	html         string
	state        TaskState
	linkTags     map[string]string
	fileType     neo.FileType
	foundLinks   []string
	parserTime   int64
	analyserTime int64
	crawlerData  *CrawlerData
	Dependencies *AnalyserDependency
	saver        neoSaverInterface
	cache        AnalyserCacheInterface
}

/*
LinkTags will get the link tags.
*/
func (analyserTask *AnalyserTaskRequest) LinkTags() map[string]string {
	return analyserTask.linkTags
}

/*
DomainInfo will get the domain info.
*/
func (analyserTask *AnalyserTaskRequest) DomainInfo() *jp.URL {
	return analyserTask.domainInfo
}

/*
SetDomainInfo will set the domain info.
*/
func (analyserTask *AnalyserTaskRequest) SetDomainInfo(domainInfo *jp.URL) {
	analyserTask.domainInfo = domainInfo
}

/*
FileType will return the file type.
*/
func (analyserTask *AnalyserTaskRequest) FileType() neo.FileType {
	return analyserTask.fileType
}

/*
SetFileType will set the file type.
*/
func (analyserTask *AnalyserTaskRequest) SetFileType() {
	analyserTask.fileType = analyserTask.extractFileType(analyserTask.DomainInfo())
}

/*
SetSaver sets the saver.
*/
func (analyserTask *AnalyserTaskRequest) SetSaver(saver neoSaverInterface) {
	analyserTask.saver = saver
}

/*
extractFileType extracts the file type of a domainInfo
*/
func (analyserTask *AnalyserTaskRequest) extractFileType(domainInfo *jp.URL) neo.FileType {
	link := domainInfo.Path
	endPoint := strings.LastIndex(link, ".")
	fileType := "html"
	if endPoint != -1 {
		fileType = link[endPoint+1:]
	}

	var ret neo.FileType
	switch fileType {
	case "html":
		ret = neo.HTML
	case "js":
		ret = neo.Javascript
	case "png":
		ret = neo.Image
	case "jpg":
		ret = neo.Image
	case "jepg":
		ret = neo.Image
	case "gif":
		ret = neo.Image
	case "css":
		ret = neo.CSS
	default:
		ret = neo.HTML
	}

	return ret
}

/*
InjectDependency will inject the dependencies for the a analyser instance
into the actual instance.
*/
func (analyserTask *AnalyserTaskRequest) InjectDependency(deps *AnalyserDependency) {
	analyserTask.Dependencies = deps
}

/*
CrawlerData returns CrawlerData.
*/
func (analyserTask *AnalyserTaskRequest) CrawlerData() *CrawlerData {
	return analyserTask.crawlerData
}

/*
SetCrawlerData sets CrawlerData.
*/
func (analyserTask *AnalyserTaskRequest) SetCrawlerData(crawlerData *CrawlerData) {
	analyserTask.crawlerData = crawlerData
}

/*
AnalyserTime returns AnalyserTime.
*/
func (analyserTask *AnalyserTaskRequest) AnalyserTime() int64 {
	return analyserTask.analyserTime
}

/*
SetAnalyserTime sets AnalyserTime.
*/
func (analyserTask *AnalyserTaskRequest) SetAnalyserTime(analyserTime int64) {
	analyserTask.analyserTime = analyserTime
}

/*
ParserTime returns ParserTime.
*/
func (analyserTask *AnalyserTaskRequest) ParserTime() int64 {
	return analyserTask.parserTime
}

/*
SetParserTime sets ParserTime.
*/
func (analyserTask *AnalyserTaskRequest) SetParserTime(parserTime int64) {
	analyserTask.parserTime = parserTime
}

/*
FoundLinks returns FoundLinks.
*/
func (analyserTask *AnalyserTaskRequest) FoundLinks() []string {
	return analyserTask.foundLinks
}

/*
SetFoundLinks sets FoundLinks.
*/
func (analyserTask *AnalyserTaskRequest) SetFoundLinks(foundLinks []string) {
	analyserTask.foundLinks = foundLinks
}

/*
State returns State.
*/
func (analyserTask *AnalyserTaskRequest) State() TaskState {
	return analyserTask.state
}

/*
SetState sets State.
*/
func (analyserTask *AnalyserTaskRequest) SetState(state TaskState) {
	analyserTask.state = state
}

/*
HTML returns Html.
*/
func (analyserTask *AnalyserTaskRequest) HTML() string {
	return analyserTask.html
}

/*
SetHTML sets Html.
*/
func (analyserTask *AnalyserTaskRequest) SetHTML(html string) {
	analyserTask.html = html
}

/*
SetWorkAddr sets the work address.
Identifies and sets the root address.
*/
func (analyserTask *AnalyserTaskRequest) SetWorkAddr(workAddr string) {
	analyserTask.workAddr = workAddr
	domainInfo, _ := jp.Parse(workAddr)
	analyserTask.SetDomainInfo(domainInfo)
}

/*
WorkAddr returns WorkAddr.
*/
func (analyserTask *AnalyserTaskRequest) WorkAddr() string {
	return analyserTask.workAddr
}

/*
getID returns Id.
*/
func (analyserTask *AnalyserTaskRequest) getID() uint64 {
	return analyserTask.id
}

/*
SetID sets Id.
*/
func (analyserTask *AnalyserTaskRequest) SetID(id uint64) {
	analyserTask.id = id
}

/*
NewTask returns a new task initialized by crawler data.
*/
func NewTask(cdata *CrawlerData /*, que *AnalyserQueue*/) *AnalyserTaskRequest {
	task := AnalyserTaskRequest{}
	task.SetCrawlerData(cdata)
	task.SetState(UNDONE)
	return &task
}

/*
SetLinkTags will set the tags which contain links.
*/
func (analyserTask *AnalyserTaskRequest) SetLinkTags() {
	analyserTask.linkTags = make(map[string]string)

	analyserTask.linkTags["meta"] = "content"
	analyserTask.linkTags["base"] = href
	analyserTask.linkTags["link"] = href
	analyserTask.linkTags["script"] = "src"
	analyserTask.linkTags["a"] = href
	analyserTask.linkTags["img"] = "src"
	analyserTask.linkTags["form"] = "action"
	analyserTask.linkTags["input"] = "value"
}

/*
Execute will execute the task.
*/
func (analyserTask *AnalyserTaskRequest) Execute(waitGroup *sync.WaitGroup, analyser *AnalyserServiceHandler) bool {
	if analyserTask.CrawlerData().getTaskError().Error() == "" {
		analyserTask.Initialize()
		analyserTask.Process(analyser)
	} else {
		analyserTask.ProcessError()
	}

	analyserTask.SetState(FINISHED)

	if waitGroup != nil {
		defer waitGroup.Done()
	}

	return true
}

/*
VerifyCorrectness will verify the correctness of the workAddr.
*/
func (analyserTask *AnalyserTaskRequest) VerifyCorrectness() bool {
	_, err := jp.Parse(strings.TrimSpace(analyserTask.workAddr))
	return err == nil
}

/*
ClearMemory will clear the memory.
*/
func (analyserTask *AnalyserTaskRequest) ClearMemory() {
	analyserTask.SetHTML("")
	analyserTask.linkTags = nil
}

/*
Process will process a valid task.
*/
func (analyserTask *AnalyserTaskRequest) Process(analyser *AnalyserServiceHandler) {
	analyserTask.SetState(PROCESSING)

	if analyserTask.VerifyCorrectness() {
		analyserTask.Analyse(analyserTask.Parse())
		analyserTask.ClearMemory()
		analyserTask.CheckLinks()

		links2crawl := analyserTask.VerifyNeo4J()
		analyserTask.Save()

		analyserTask.SetFoundLinks(links2crawl)

		// NEW
		newLinks := make([]string, 0)
		for _, link := range analyserTask.FoundLinks() {
			if !analyserTask.cache.Request(link) {
				newLinks = append(newLinks, link)
			}
		}
		analyserTask.SetFoundLinks(newLinks)

		analyserTask.CollectLinks()
		if analyserTask.saver.GetSession() != nil {
			defer neo.CloseSession(analyserTask.saver.GetSession())
		}
		analyserTask.NextSend(analyser)
	}
}

/*
Initialize will initialize the task with the given crawler data.
*/
func (analyserTask *AnalyserTaskRequest) Initialize() {
	adr := analyserTask.CrawlerData().getAddr()
	adr = strings.TrimSpace(adr)

	analyserTask.SetWorkAddr(adr)
	analyserTask.SetHTML(string(analyserTask.CrawlerData().responseBodyBytes))
	analyserTask.SetLinkTags()
	analyserTask.SetFileType()
}

/*
Parse will parse the html to create a tree with a root node.
*/
func (analyserTask *AnalyserTaskRequest) Parse() *model.Node {
	htmlTree := model.NewHTMLTree(analyserTask.HTML())
	start := time.Now()
	htmlTree.Parse(false)
	rootNode := htmlTree.RootNode()

	analyserTask.parserTime = time.Since(start).Nanoseconds()

	return rootNode
}

/*
Analyse will analyse the root node to get the links of the site.
*/
func (analyserTask *AnalyserTaskRequest) Analyse(rootNode *model.Node) {
	start := time.Now()
	analyserTask.Traverse(rootNode)
	analyserTask.analyserTime = time.Since(start).Nanoseconds()
}

/*
CollectLinks will collect all links which will be sent to the crawler
*/
func (analyserTask *AnalyserTaskRequest) CollectLinks() {
	analyserTask.RemoveExternals()
	analyserTask.RemoveFileTypes()
	analyserTask.CorrectLinks()

}

/*
VerifyNeo4J will verify if the address is already contained in neo4j
*/
func (analyserTask *AnalyserTaskRequest) VerifyNeo4J() []string {
	return analyserTask.saver.Contains(analyserTask)
}

/*
Save will save the valid links to Neo4J
*/
func (analyserTask *AnalyserTaskRequest) Save() {
	analyserTask.SetState(SAVING)
	analyserTask.saver.Save(analyserTask)
}

/*
ProcessError will process an error task.
*/
func (analyserTask *AnalyserTaskRequest) ProcessError() {
	analyserTask.SetState(CRAWLERERROR)
	analyserTask.saver.Save(analyserTask)
}

/*
Traverse will traverse the tree and classify every node.
*/
func (analyserTask *AnalyserTaskRequest) Traverse(node *model.Node) {
	for _, ele := range node.Children() {
		if len(ele.Children()) > 0 {
			analyserTask.Traverse(ele)
		} else {
			analyserTask.Classify(ele)
		}
	}
}

/*
Classify inspects a node whether its a "link" or not.
*/
func (analyserTask *AnalyserTaskRequest) Classify(node *model.Node) {
	tag := node.Tag()

	if attributeType, ok := analyserTask.LinkTags()[tag.TagType()]; ok {
		for _, attribute := range tag.Attributes() {
			if attributeType == attribute.AttributeType() {
				analyserTask.HandleValue(attribute.Value())
			}
		}
	}
}

/*
Send will send all valid links to the crawler as new task.
*/
func (analyserTask *AnalyserTaskRequest) Send() {
	numberVar := 300
	analyserTask.SetState(SENDTOCRAWLER)

	log.WithFields(log.Fields{
		"analyserID": analyserTask.getID(),
		"links":      analyserTask.FoundLinks(),
	}).Info("SEND LINKS TO CRAWLER")

	serv := analyserTask.Dependencies.Crawler()

	for _, link := range analyserTask.FoundLinks() {
		r := rand.Intn(numberVar)
		time.Sleep(time.Duration(r) * time.Millisecond)

		message := &crawlerproto.CrawlTaskCreateRequest{
			Url: link,
		}
		_, err := serv.CreateTask(context.TODO(), message)
		if err != nil {
			log.Error("Error while sending to crawler")
			log.Error(err)
		}
	}
}



func (analyserTask *AnalyserTaskRequest) NextSend(analyser *AnalyserServiceHandler) {
	analyserTask.SetState(SENDTOCRAWLER)

	log.WithFields(log.Fields{
		"analyserID": analyserTask.getID(),
		"links":      analyserTask.FoundLinks(),
	}).Info("SEND LINKS TO CRAWLER")

	for _, link := range analyserTask.FoundLinks() {
		err := analyser.SendUrlToCrawler(context.TODO(), link)
		if err != nil {
			log.Error("Error while sending to crawler")
			log.Error(err)
		}
	}
}

/*
HandleValue will handle a possible link.
*/
func (analyserTask *AnalyserTaskRequest) HandleValue(link string) {
	if link, err := analyserTask.VerifyLink(link); err == nil {
		analyserTask.foundLinks = append(analyserTask.foundLinks, link)
	}
}

/*
VerifyLink will verify whether link is valid or not.
*/
func (analyserTask *AnalyserTaskRequest) VerifyLink(link string) (string, error) {
	switch {
	// Empty link
	case link == "":
		return "", fmt.Errorf("link is empty: %v", link)
	// Relative link
	case strings.Contains(link, "width=device"):
		return "", fmt.Errorf("link is empty: %v", link)
	case len(link) > 2 && link[0] == '/' && link[1] != '/':
		return analyserTask.DomainInfo().Scheme + "://" + analyserTask.DomainInfo().Hostname() + link, nil
	case !strings.Contains(link, "/") && strings.Contains(link, "."):
		return analyserTask.DomainInfo().Scheme + "://" + analyserTask.DomainInfo().Hostname() + "/" + link, nil
	// Onion links
	case !strings.HasPrefix(link, "www") && !strings.HasPrefix(link, "http"):
		return "", fmt.Errorf("link is invalid: %v", link)
	default:
		return link, nil

	}
}

/*
CheckLinks will check a link before saved.
*/
func (analyserTask *AnalyserTaskRequest) CheckLinks() {
	analyserTask.RemoveDuplicates()
}

/*
RemoveFileTypes will remove all to big files.
*/
func (analyserTask *AnalyserTaskRequest) RemoveFileTypes() {
	list := []string{
		"pdf", "zip", "png", "jpg", "jpeg", "vdi", "iso", "zst",
	}

	links := make([]string, 0)

	for _, link := range analyserTask.FoundLinks() {
		no := false
		for _, fileType := range list {
			if strings.Contains(link, fileType) {
				no = true
			}
		}

		if !no {
			links = append(links, link)
		}
	}

	analyserTask.SetFoundLinks(links)
}

/*
RemoveDuplicates removes duplicated links.
*/
func (analyserTask *AnalyserTaskRequest) RemoveDuplicates() {
	keys := make(map[string]bool)
	cleanLinks := make([]string, 0)
	for _, entry := range analyserTask.FoundLinks() {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			cleanLinks = append(cleanLinks, entry)
		}
	}

	analyserTask.SetFoundLinks(cleanLinks)
}

/*
RemoveExternals will remove external links.
*/
func (analyserTask *AnalyserTaskRequest) RemoveExternals() {
	var cleanLinks []string
	for _, link := range analyserTask.FoundLinks() {
		domainInfo, err := jp.Parse(strings.TrimSpace(link))
		if err != nil {
			log.WithFields(log.Fields{
				"link": link,
			}).Info("DOMAIN INFO")
		} else {

			rootDomain := analyserTask.DomainInfo().Domain + "." + analyserTask.DomainInfo().TLD
			linkDomain := domainInfo.Domain + "." + domainInfo.TLD

			if rootDomain == linkDomain {
				cleanLinks = append(cleanLinks, link)
			}
		}
	}

	analyserTask.SetFoundLinks(cleanLinks)
}

/*
CorrectLinks will correct wrong links.
*/
func (analyserTask *AnalyserTaskRequest) CorrectLinks() {
	var cleanLinks []string
	for _, link := range analyserTask.FoundLinks() {
		domainInfo, err := jp.Parse(strings.TrimSpace(link))
		if err != nil {
			log.WithFields(log.Fields{
				"link": link,
			}).Info("DOMAIN INFO")
		} else {
			newLink := link
			if len(domainInfo.Path) != 0 {
				path := template.URLQueryEscaper(domainInfo.Path[1:])
				newLink = domainInfo.Scheme + "://" + domainInfo.Host + "/" + path
			}

			cleanLinks = append(cleanLinks, newLink)
		}
	}

	analyserTask.SetFoundLinks(cleanLinks)
}
