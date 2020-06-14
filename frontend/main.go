package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/web"
)

const (
	servicename string = "Reactservice"
	address     string = "0.0.0.0:8080"
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
	fmt.Println("[+] Initialized the webserver for the frontend.")
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Static("/static", "sherlockgopherfrontend/build/static")
	router.StaticFile("/", "sherlockgopherfrontend/build/index.html")

	if err := router.Run(getAddress()); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] Started the webserver. Listening on %s \n", getAddress())
	service.Handle("/", router)
	fmt.Println("[+] Started the service.")
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
