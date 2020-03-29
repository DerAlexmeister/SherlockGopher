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

/*
New ill return a new ServeReact service.
*/
func New() *ServeReact {
	return &ServeReact{}
}

/*
ServeReact this service will serve the react files.
*/
type ServeReact struct{}

func (react *ServeReact) serveIndex(context *gin.Context) {
	context.JSON(200, map[string]string{
		"message": "Welcome json!",
	})
}

func main() {
	service := web.NewService(web.Name(getServiceName()))

	service.Init()

	reactService := New()

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))
	router.GET("/testing", reactService.serveIndex)

	router.Run(getAddress())

	service.Handle("/", router)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
