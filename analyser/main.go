package main

import (
	"os"

	"github.com/micro/go-micro"
	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto"
	"github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockanalyser"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto"
	"github.com/ob-algdatii-20ss/SherlockGopher/sherlockneo"
	log "github.com/sirupsen/logrus"
)

const (
	serviceName = "analyser-service"
)

func main() {
	SetupLogging()
	log.Info("Started analyser")

	service := micro.NewService(
		micro.Name(serviceName),
	)
	service.Init()

	AnalyserService := sherlockanalyser.NewAnalyserServiceHandler()

	dep := sherlockanalyser.AnalyserDependency{}

	if driver, err := sherlockneo.GetNewDatabaseConnection(); driver == nil || err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Neo4J failed!")
	} else {
		dep.Neo4J = &driver
	}

	dep.Crawler = func() crawlerproto.CrawlerService {
		return crawlerproto.NewCrawlerService("crawler-service", service.Client())
	}

	if session, err := sherlockneo.GetSession(*dep.Neo4J); session == nil || err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"session": session,
		}).Error("getting session failed")
	} else {
		sherlockneo.RunConstrains(session)
	}

	AnalyserService.InjectDependency(&dep)

	err := proto.RegisterAnalyserHandler(service.Server(), AnalyserService)

	go AnalyserService.ManageTasks()

	if err != nil {
		log.Fatal("Analyser->main.go->RegisterAnalyserHandler failed!")
		log.Fatal(err)
	} else if err = service.Run(); err != nil {
		log.Fatal("Analyser->main.go->service.Run() failed!")
		log.Fatal(err)
	} else {
		log.Infof("Service %s started as intended... ", serviceName)
	}
}

/*
SetupLogging will init all needed things for logging.
*/
func SetupLogging() {
	_ = os.Remove("info.log")
	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0644)

	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		ForceQuote:                true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
	})

	log.SetOutput(file)
}
