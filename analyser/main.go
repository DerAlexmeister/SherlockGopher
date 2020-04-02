package main

import (
	"github.com/micro/go-micro"
)

const serviceName = "analyser-service"

func main() {

	service := micro.NewService(
		micro.Name(serviceName),
	)

	//err := proto.AnalyserService()

	service.Init()

	/*
		if err := service.Run(); err != nil {
			log.Fatalf("Failed to start service. Error:\n\t%s", err.Error())
		}
	*/
}
