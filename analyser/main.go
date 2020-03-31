package main

import (
	"github.com/micro/go-micro"
	proto "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/analyser/proto"
	impl "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/analyser/analyser"
)

const serviceName = "analyser-service"

func main(){

	service := micro.NewService(
		micro.Name(serviceName),
	)

	err := proto.RegisterAnalyserHandler(service.Server(), impl.NewAnalyserService())

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatalf("Failed to start service. Error:\n\t%s", err.Error())
	
}