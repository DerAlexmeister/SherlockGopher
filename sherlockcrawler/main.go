package main

import (
	"log"

	"github.com/micro/go-micro"
)

const serviceName = "crawler-service"

func main() {

	service := micro.NewService(micro.Name(serviceName))
	service.Init()

	if err := service.Run(); err != nil {
		log.Fatalf("Failed to start service. Error:\n\t%s", err.Error())
	}
}
