package main

import (
	"log"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/web"
)

const (
	servicename string = "Reactservice"
	address     string = "localhost:8080"
)

/*
getServiceName will return the name of the service.
*/
func getServiceName() string {
	return servicename
}

/*
getAddress will return the address of the webserver.
*/
func getAddress() string {
	return address
}

func main() {
	service := web.NewService(web.Name(getServiceName()))

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	router.Run(getAddress())

	service.Handle("/", router)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
