package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
CrawlWebServer will be the webserver being the man in the middle between the frontend and the backend.
*/
type CrawlWebServer struct {
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
New will return a new instance of the CrawlWebServer
*/
func New() *CrawlWebServer {
	return &CrawlWebServer{}
}

/*
Helloping Ping will return just for testing purposes a pong. Like PING PONG.
*/
func (server *CrawlWebServer) Helloping(context *gin.Context) {
	context.JSON(200, map[string]string{
		"message": "Pong",
	})
}

/*
RecieveURL will handle the requested url which should be crawled.
*/
func (server *CrawlWebServer) RecieveURL(context *gin.Context) {
	var url = NewRequestedURL()
	context.BindJSON(url)
	context.JSON(http.StatusOK, gin.H{
		"status": "Fine",
	})
	fmt.Println(url) //TODO check if url is empty or a well formed url.
}
