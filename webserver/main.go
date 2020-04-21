package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/web"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertowebserver"
	webserver "github.com/ob-algdatii-20ss/SherlockGopher/webserver/webserver"
)

const (
	servicename        string = "sherlockwebserver" // Name of the Service
	serviceNameCrawler string = "crawler-service"   //Name of the cralwer
	address            string = "0.0.0.0:8081"      // Address of the Webserver
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

	grpcservice := micro.NewService(micro.Name(servicename))
	grpcerr := service.Init()
	if grpcerr != nil {
		log.Fatal(err)
	}
	webServerService := webserver.New()
	webServerService.SetCrawlerServiceDependency(&webserver.SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawlerproto.NewCrawlerService(serviceNameCrawler, grpcservice.Client())
		},
	})
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}

	router.Use(cors.New(config))

	//Get Requests.
	router.GET("/areyouthere", webServerService.Helloping)

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
	graphsapi.GET("/meta", webServerService.GraphMetaV1)                             // Get all meta information about neo4j.
	graphsapi.GET("/all", webServerService.GraphFetchWholeGraphV1)                   // Will return the entire graph, maybe build a stream.
	graphsapi.GET("/performenceofsites", webServerService.GraphPerformenceOfSitesV1) //Will return address with statuscode and reponsetime.

	graphsapi.POST("/detailsofnode", webServerService.GraphNodeDetailsV1) // get all information of a node.
	graphsapi.POST("/search", webServerService.ReceiveURL)

	router.Run(getAddress())

	service.Handle("/", router)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
