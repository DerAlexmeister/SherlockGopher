package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertowebserver"
	sherlockneo "github.com/ob-algdatii-20ss/SherlockGopher/sherlockneo"
)

/*
SherlockWebserver will be the webserver being the man in the middle between the frontend and the backend.
*/
type SherlockWebserver struct {
	Dependency *SherlockWebServerDependency
	Driver     *neo4j.Driver
}

/*
SherlockWebServerDependency will be all dependencies needed for the Webserver to run.
*/
type SherlockWebServerDependency struct {
	Crawler func() crawlerproto.CrawlerService
}

/*
SetCrawlerServiceDependency will set the dependency for the sherlockwebserver package.
*/
func (server *SherlockWebserver) SetCrawlerServiceDependency(deps *SherlockWebServerDependency) {
	server.Dependency = deps
}

/*
Helloping Ping will return just for testing purposes a pong. Like PING PONG.
*/
func (server *SherlockWebserver) Helloping(context *gin.Context) {
	context.JSON(200, map[string]string{
		"message": "Pong",
	})
}

/*
RequestedURL will be the struct for the post request of the domain to crawl.
*/
type RequestedURL struct {
	URL string `json:"url" binding:"required"`
}

/*
NewRequestedURL will be a new instance of RequestedURL.
*/
func NewRequestedURL() *RequestedURL {
	return &RequestedURL{}
}

/*
New will return a new instance of the SherlockWebserver
*/
func New() *SherlockWebserver {
	ldriver, err := sherlockneo.GetNewDatabaseConnection()
	if err == nil {
		return &SherlockWebserver{
			Driver: &ldriver,
		}
	}
	return &SherlockWebserver{}
}

/*
RecieveURL will handle the requested url which should be crawled.
*/
func (server *SherlockWebserver) RecieveURL(context *gin.Context) {
	var url = NewRequestedURL()
	context.BindJSON(url)
	context.JSON(http.StatusOK, gin.H{
		"Status": "Fine",
	}) //Send fine as response.
	fmt.Println(url)

	//TODO check if url is empty or a well formed url.
	//TODO send to crawler.
}
