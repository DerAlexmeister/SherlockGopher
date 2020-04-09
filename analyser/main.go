package main

import (
	"github.com/micro/go-micro"
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
	/*
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

		newService := streamreceiver.NewServerGRPC()

		err1 := proto.RegisterReceiverHandler(service.Server(), newService)
		if err1 == nil {
			if err := service.Run(); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err1)
		}
	*/
}
