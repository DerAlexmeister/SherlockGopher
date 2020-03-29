package main

import (
	"github.com/micro/go-micro"
	proto "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/crawler/proto"
	impl "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/crawler/crawler"
)

const serviceName = "crawler-service"

func main(){

	service := micro.NewService(
		micro.Name(serviceName),
	)

	err := proto.RegisterCrawlerHandler(service.Server(), impl.NewCrawlerService())

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatalf("Failed to start service. Error:\n\t%s", err.Error())
	
}