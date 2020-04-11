package main

import (
	"fmt"

	"github.com/micro/go-micro"
	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
	sherlock "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/sherlockcrawler"
)

const serviceName = "crawler-service"

func main() {
	// CrawlerService.
	service := micro.NewService(micro.Name(serviceName))
	service.Init()

	crawlerservice := sherlock.NewSherlockCrawlerService()

	deps := sherlock.NewSherlockDependencies()

	//TODO missing setters for the Dependencies.

	crawlerservice.InjectDependency(deps)

	err := proto.RegisterAnalyserInterfaceHandler(service.Server(), crawlerservice)

	go crawlerservice.ManageTasks()

	if err != nil {
		fmt.Println(err)
	} else if lerr := service.Run(); lerr != nil {
		fmt.Println(lerr)
	} else {
		fmt.Printf("Service %s started as intended... ", serviceName)
	}

	//Filetransferservice.
}
