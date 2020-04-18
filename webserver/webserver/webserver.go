package webserver

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertowebserver"
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

type TaskStatusResponse struct {
    website string
    undone int
    processing int
    finished int
    failed int
}

/*
RecieveURL will handle the requested url which should be crawled.
*/
func (server *CrawlWebServer) ReceiveURL(context *gin.Context) {
	var url = NewRequestedURL()
	context.BindJSON(url)
	context.JSON(http.StatusOK, gin.H{
		"Status": "Fine",
	}) //Send fine as response.
	fmt.Println(url)

	if url.URL != "" && govalidator.IsURL(url.URL) {
		sherlockcrawlerService := server.dependencies.
		in := &crawlerproto.TaskStatusRequest{}
		response, err := sherlockcrawlerService.StatusOfTaskQueue(context, in)

		if (err == nil && response != nil) {
			context.JSON(http.StatusOK, gin.H{
				"Website": response.Website,
				"Undone": response.Undone,
				"Processing": response.Processing,
				"Finished": response.Finished,
				"Failed": response.Failed,
			}) 
		}

	
	}
	//TODO check if url is empty or a well formed url.
	//TODO send to crawler.
}

