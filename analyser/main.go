package main

import (
	"fmt"

	"github.com/micro/go-micro"
	proto "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/analyser/proto"
	"github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/analyser/sherlockanalyser"
)

const (
	serviceName      = "analyser-service"
	streamingService = "filestransfer-service"
)

func main() {

	// Analyserservice
	service := micro.NewService(
		micro.Name(serviceName),
	)
	service.Init()

	AnalyserService := sherlockanalyser.NewAnalyserServiceHandler()

	AnalyserService.InjectDependency(sherlockanalyser.NewAnalyserDependencies())

	err := proto.RegisterAnalyserHandler(service.Server(), AnalyserService)

	if err != nil {
		fmt.Println(err)
	} else if lerr := service.Run(); lerr != nil {
		fmt.Println(lerr)
	} else {
		fmt.Printf("Service %s started as intended... ", serviceName)
	}
	// FileTransferService.

	streamingservice := micro.NewService(
		micro.Name(streamingService),
	)

	streamingservice.Init()

}
