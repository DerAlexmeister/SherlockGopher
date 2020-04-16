package main

import (
	"fmt"

	"github.com/micro/go-micro"
	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
	sherlock "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/sherlockcrawler"
)

const (
	serviceName     = "crawler-service"
	fileservicename = "filestransfer-service-crawler"
)

func main() {
	// CrawlerService.
	service := micro.NewService(micro.Name(serviceName))
	service.Init()

	crawlerservice := sherlock.NewSherlockCrawlerService()
	deps := sherlock.NewSherlockDependencies()
	streamingserver := sherlock.NewStreamingServer()

	//TODO missing setters for the Dependencies.
	crawlerservice.InjectDependency(deps)
	crawlerservice.SetSherlockStreamer(&streamingserver) // Add the current streaminserver to the current sherlock crawler.

	err := proto.RegisterAnalyserInterfaceHandler(service.Server(), crawlerservice)

	if err != nil {
		fmt.Println(err)
	} else if lerr := service.Run(); lerr != nil {
		fmt.Println(lerr)
	} else {
		fmt.Printf("Service %s started as intended... ", serviceName)
	}

	//Filestreaming Service.

	//streamingService := micro.NewService(
	//	micro.Name(fileservicename),
	//)

	//service.Init()

	//fileproto.RegisterSenderHandler(streamingService.Server(), streamingserver)

	go crawlerservice.ManageTasks() //Maybe a waitgroup needed.
	//go streamingserver.Upload(context.TODO())

}
