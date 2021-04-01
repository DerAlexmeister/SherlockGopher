package main

import (
	"fmt"
	"log"

	analyserproto "github.com/DerAlexx/SherlockGopher/analyser/proto"
	crawlerproto "github.com/DerAlexx/SherlockGopher/sherlockcrawler/proto"
	webserver "github.com/DerAlexx/SherlockGopher/webserver/webserver"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/web"
)

const (
	servicename         string = "sherlockwebserver" // Name of the Service
	serviceNameCrawler  string = "crawler-service"   //Name of the cralwer
	serviceNameAnalyser string = "analyser-service"
	address             string = "0.0.0.0:8081" // Address of the WebServer
	acceptcors          bool   = true
)

//CORS are all accepted Cross-Origin-Resource-Sharing-Addresses.
var CORS []string = []string{"http://localhost:8080", "http://localhost:3000"}

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

/*
getCors will return all cors
*/
func getCors() []string {
	return CORS
}

/*
acceptCors will return true incase cors should be accepted.
*/
func acceptCors() bool {
	return acceptcors
}

func main() {
	service := web.NewService(web.Name(getServiceName()))
	err := service.Init()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] Initialized the webserver.")
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
		Analyser: func() analyserproto.AnalyserService {
			return analyserproto.NewAnalyserService(serviceNameAnalyser, grpcservice.Client())
		},
	})
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if acceptCors() {
		config := cors.DefaultConfig()
		config.AllowOrigins = getCors()
		router.Use(cors.New(config))
		fmt.Println("[+] Accepting the following CORS", getCors())
	}

	//Get Requests.
	router.GET("/areyouthere", webServerService.Helloping)

	//Controller Group.
	controller := router.Group("/controller/v1")
	controller.POST("/changestate", webServerService.ChangeState) //will change the state of the cralwer and the analyser.
	controller.GET("/status", webServerService.GetServiceStatus)  // will return the status of the analyser/crawler service.
	controller.GET("/dropit", webServerService.DropGraphTable)    // will drop the neo4j table.

	//new services ba
	//controller.GET("/getmetadata/:page/", webServerService.GetMetaData)
	controller.GET("/getscreenshots/:page/", webServerService.GetScreenshots)

	//Monitor Group.
	monitorapi := router.Group("/monitor/v1")                 //missing handler.
	monitorapi.GET("/meta", webServerService.ReceiveMetadata) // Get all meta information about the crawler and the analyser.

	//Graph Group.
	graphsapi := router.Group("/graph/v1")
	graphsapi.GET("/meta", webServerService.GraphMetaV1) // Get all meta information about neo4j.
	graphsapi.GET("/all", webServerService.GraphFetchWholeGraphV1)
	graphsapi.GET("/alloptimized", webServerService.GraphFetchWholeGraphHighPerformanceV1) // Will return the entire graph, maybe build a stream.
	graphsapi.GET("/performenceofsites", webServerService.GraphPerformanceOfSitesV1)       //Will return address with statuscode and reponsetime.

	graphsapi.POST("/detailsofnode", webServerService.GraphNodeDetailsV1) // get all information of a node.
	graphsapi.POST("/search", webServerService.ReceiveURL)                // will handle the requested url which should be crawled.

	err = router.Run(getAddress())
	if err != nil {
		log.Fatal(err)
	}

	service.Handle("/", router)
	fmt.Printf("[+] Started the webserver. Listening on %s \n", getAddress())
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
