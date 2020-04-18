package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/web"
	webserver "github.com/ob-algdatii-20ss/SherlockGopher/webserver/webserver"
)

const (
	servicename string = "CrawlWebServer" // Name of the Service
	address     string = "0.0.0.0:8081"   // Address of the Webserver
)

/*
getServiceName will return a the service name of this service.
*/
func getServiceName() string {
	return servicename
}

/*
getAddress will return the address of the service.
*/
func getAddress() string {
	return address
}

func main() {
	service := web.NewService(web.Name(getServiceName()))

	err := service.Init()

	if err != nil {
		log.Fatal(err)
	}

	webServerService := webserver.New()

	router := gin.Default()

	//Get Requests.
	router.GET("/areyouthere", webServerService.Helloping)

	//POST Requests.
	router.POST("/search", webServerService.RecieveURL)

	//DONT implement this yet.
	controller := router.Group("/controller/v1")
	controller.GET("/stop")   //will stop the cralwer and the analyser
	controller.GET("/pause")  //will pause the cralwer and the analyser.
	controller.GET("/resume") //will resume the cralwer and the analyser.
	controller.GET("/clean")  //will clean the the cralwer and the analyser works like stop but will clean the queues.

	//Monitor Group.
	monitorapi := router.Group("/monitor/v1") //missing handler
	monitorapi.GET("/meta")                   // Get all meta information about the crawler and the analyser.

	//Graph Group.
	graphsapi := router.Group("/graph/v1")
	graphsapi.GET("/meta")     // Get all meta information about neo4j.
	graphsapi.GET("/all")      // Will return the entire graph, maybe build a stream.
	graphsapi.POST("/snipped") // Snipped search, submit a target to get a snipped.

	router.Run(getAddress())

	service.Handle("/", router)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
