package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertowebserver"
)

/*
SherlockWebserver will be the webserver being the man in the middle between the frontend and the backend.
*/
type SherlockWebserver struct {
	Dependency *SherlockWebServerDependency
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
	return &SherlockWebserver{}
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
RecieveURL will handle the requested url which should be crawled.
*/
func (server *SherlockWebserver) ReceiveURL(context *gin.Context) {
	sherlockcrawlerService := server.Dependency.Crawler()

	var url = NewRequestedURL()
	context.BindJSON(url)

	//check if url is empty or a well formed url.
	if govalidator.IsURL(url.URL) {
		context.JSON(http.StatusOK, gin.H{
			"Status": "Fine",
		}) //Send fine as response.
		fmt.Println(url)
		//send to crawler
		didSendCount := 0
		for didSend := false; !didSend; {
			if didSendCount < 5 {
				response, err := sherlockcrawlerService.ReceiveURL(context, &crawlerproto.SubmitURLRequest{URL: url.URL})
				if err == nil && response.Recieved {
					didSend = true
				} else {
					didSendCount++
					time.Sleep(100 * time.Millisecond)
				}
			} else {
				context.JSON(http.StatusInternalServerError, gin.H{
					"Status": "Cant send url to crawler",
				})
				didSend = true
			}
		}

	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"Status": "Url was empty or malformed",
		})
	}
}

/*
RecieveURL will handle the requested url which should be crawled.
*/
func (server *SherlockWebserver) ReceiveMetadata(context *gin.Context) {
	sherlockcrawlerService := server.Dependency.Crawler()
	in_crawler := &crawlerproto.TaskStatusRequest{}
	response_crawler, err_crawler := sherlockcrawlerService.StatusOfTaskQueue(context, in_crawler)

	//TODO Analyser Dependency
	/*
		sherlockanalyserService := server.Dependency.Analyser()
		in_analyser := &analyserproto.TaskStatusRequest{}
		response_analyser, err_analyser := sherlockanalyserService.StatusOfTaskQueue(context, in_analyser)
	*/

	if err_crawler == nil /*&& err_analyser == nil*/ {
		context.JSON(http.StatusOK, gin.H{
			"Website":    response_crawler.Website,
			"Undone":     response_crawler.Undone,
			"Processing": response_crawler.Processing,
			"Finished":   response_crawler.Finished,
			"Failed":     response_crawler.Failed,

			//TODO Analyser StatusOfTaskQueue out senden
		})
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Couldnt get Metadata",
		})
	}
}
