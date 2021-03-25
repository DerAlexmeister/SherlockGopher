package main

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	aproto "github.com/DerAlexx/SherlockGopher/analyser/proto"
	proto "github.com/DerAlexx/SherlockGopher/sherlockcrawler/proto"
	"github.com/micro/go-micro"

	sherlock "github.com/DerAlexx/SherlockGopher/sherlockcrawler/sherlockcrawler"
)

const (
	serviceName = "crawler-service"
)

func main() {
	SetupLogging()
	log.Info("Started analyser")

	service := micro.NewService(micro.Name(serviceName))
	service.Init()

	fmt.Printf("[+] Successfully initialized the serivce %s", serviceName)

	crawlerservice := sherlock.NewSherlockCrawlerService()

	deps := sherlock.SherlockDependencies{
		Analyser: func() aproto.AnalyserService {
			return aproto.NewAnalyserService("analyser-service", service.Client())
		},
	}
	fmt.Printf("[+] Injected dependencies in %s", serviceName)

	crawlerservice.InjectDependency(&deps)

	err := proto.RegisterCrawlerHandler(service.Server(), crawlerservice) //ändern

	if err != nil {
		log.Fatal("Crawler->main.go->RegisterCrawlerHandler failed!")
		log.Fatal(err)
	}

	ctx := context.Background()
	go crawlerservice.ManageTasks()
	crawlerservice.Consume(ctx)

	if err = service.Run(); err != nil {
		log.Fatal("Crawler->main.go->service.Run() failed!")
		log.Fatal(err)
	} else {
		log.Infof("Service %s started as intended... ", serviceName)
	}
}

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
