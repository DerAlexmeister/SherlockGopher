package sherlockanalyser

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	neo "github.com/DerAlexx/SherlockGopher/sherlockneo"
	jp "github.com/jpillora/go-tld"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type neoSaverInterface interface {
	Contains(analyserTask *AnalyserTaskRequest) []string
	Save(analyserTask *AnalyserTaskRequest)
	GetSession() *neo4j.Session
	SetDriver(driver neo4j.Driver)
}

type neoSaver struct {
	driver  neo4j.Driver
	session neo4j.Session
	mutex   *sync.Mutex
}

func newNeoSaver() neoSaver {
	driver, _ := neo.GetNewDatabaseConnection()
	session, _ := neo.GetSession(driver)
	mutex := sync.Mutex{}
	return neoSaver{
		driver:  driver,
		session: session,
		mutex:   &mutex,
	}
}

func (h *neoSaver) SetDriver(driver neo4j.Driver) {
	h.driver = driver
}

func (h neoSaver) GetSession() *neo4j.Session {
	return &h.session
}

func (h neoSaver) Contains(analyserTask *AnalyserTaskRequest) []string {
	h.mutex.Lock()
	newLinks := make([]string, 0)
	for _, link := range analyserTask.FoundLinks() {
		if !neo.ContainsNode(h.session, link) {
			newLinks = append(newLinks, link)
		}
	}

	h.mutex.Unlock()
	return newLinks
}

func (h neoSaver) Save(analyserTask *AnalyserTaskRequest) {
	h.mutex.Lock()
	crawledAddress := neo.NewNeoLink(analyserTask.WorkAddr(), analyserTask.FileType())

	linksToSave := make([]*neo.NeoLink, 0)

	for _, link := range analyserTask.FoundLinks() {
		domainInfo, _ := jp.Parse(strings.TrimSpace(link))
		neoLink := neo.NewNeoLink(link, analyserTask.extractFileType(domainInfo))
		linksToSave = append(linksToSave, neoLink)
	}

	responseHeader := analyserTask.CrawlerData().getResponseHeader()
	neoData := neo.NewNeoData(crawledAddress, analyserTask.CrawlerData().getStatusCode(), analyserTask.CrawlerData().getResponseTime(),
		&responseHeader, analyserTask.CrawlerData().getTaskError().Error(), linksToSave)

	err := neoData.Save(h.driver)
	if err != nil {
		log.Info(err)
	}
	h.mutex.Unlock()
}
