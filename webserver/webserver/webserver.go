package webserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	analyserproto "github.com/DerAlexx/SherlockGopher/analyser/proto"
	crawlerproto "github.com/DerAlexx/SherlockGopher/sherlockcrawler/proto"
	sherlockneo "github.com/DerAlexx/SherlockGopher/sherlockneo"
)

/*
SherlockWebserver will be the webserver being the man in the middle between the frontend and the backend.
*/
type SherlockWebserver struct {
	Dependency *SherlockWebServerDependency
	Driver     neo4j.Driver
}

/*
SherlockWebServerDependency will be all dependencies needed for the WebServer to run.
*/
type SherlockWebServerDependency struct {
	Crawler  func() crawlerproto.CrawlerService
	Analyser func() analyserproto.AnalyserService
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
	context.JSON(http.StatusOK, map[string]string{
		"Message": "Yes i am here!",
	})
}

/*
RequestedURL will be the struct for the post request of the domain to crawl.
*/
type RequestedURL struct {
	URL string `json:"url" binding:"required"`
}

/*
RequestedStatus will be the struct for the post request of the status functions.
*/
type RequestedStatus struct {
	Operation string `json:"operation" binding:"required"`
	Target    string `json:"target" binding:"required"`
}

/*
MetaArray stores the JSON which contains the response of the crawler and the analyser after the webserver requests their status.
Help function of ReceiveMetadata.
*/
type MetaArray struct {
	metamap map[string]interface{}
}

/*
NewRequestedURL will be a new instance of RequestedURL.
*/
func NewRequestedURL() *RequestedURL {
	return &RequestedURL{}
}

/*
NewRequestedStatus will be a new instance of RequestedStatus.
*/
func NewRequestedStatus() *RequestedStatus {
	return &RequestedStatus{}
}

/*
New will return a new instance of the SherlockWebserver.
*/
func New() *SherlockWebserver {
	ldriver, err := sherlockneo.GetNewDatabaseConnection()
	if err == nil {
		return &SherlockWebserver{
			Driver: ldriver,
		}
	}
	return &SherlockWebserver{}
}

/*
ReceiveURL will handle the requested url which should be crawled.
*/
func (server *SherlockWebserver) ReceiveURL(ctx *gin.Context) {
	sherlockcrawlerService := server.Dependency.Crawler()
	var url = NewRequestedURL()
	err := ctx.BindJSON(url)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Error while reveiving Requested Url",
		})
	}
	if govalidator.IsURL(url.URL) {
		response, err := sherlockcrawlerService.ReceiveURL(context.TODO(), &crawlerproto.SubmitURLRequest{URL: url.URL})
		if err == nil && response.Recieved {
			ctx.JSON(http.StatusOK, gin.H{
				"Status": "Fine",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"Message": "The webserver cannot submit your URL to the Crawler, couldnt reach the crawler service",
			})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "Url was empty or malformed",
		})
	}
}

/*
ReceiveMetadata will get the status of the tasks inside the analyser and crawler queue.
*/
func (server *SherlockWebserver) ReceiveMetadata(context *gin.Context) {

	responseAnalyser, errAnalyser := server.receiveStatusTaskQAnalyser(context)
	responseCrawler, errCrawler := server.receiveStatusTaskQCrawler(context)

	switch {
	case errCrawler != nil && errAnalyser != nil:
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Couldnt get Metadata, analyser and crawler services are unavailable",
		})
	case errCrawler != nil:
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Couldnt get Metadata, crawler service is unavailable",
		})
	case errAnalyser != nil:
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Couldnt get Metadata, analyser service is unavailable",
		})
	default:
		metaArray := fillMetaArray(responseCrawler, responseAnalyser)
		context.JSON(http.StatusOK, metaArray.metamap)
	}
}

/*
receiveStatusTaskQAnalyser will get the status of the tasks in the analyser queue.
Help function of ReceiveMetadata.
*/
func (server *SherlockWebserver) receiveStatusTaskQAnalyser(context *gin.Context) (*analyserproto.WorkloadResponse, error) {
	sherlockanalyserService := server.Dependency.Analyser()
	if sherlockanalyserService == nil {
		return nil, errors.New("analyser is not alive")
	}
	inAnalyser := &analyserproto.WorkloadRequest{}
	responseAnalyser, errAnalyser := sherlockanalyserService.WorkloadRPC(context, inAnalyser)
	return responseAnalyser, errAnalyser
}

/*
receiveStatusTaskQCrawler will get the status of the tasks in the crawler queue.
Help function of ReceiveMetadata.
*/
func (server *SherlockWebserver) receiveStatusTaskQCrawler(context *gin.Context) (*crawlerproto.TaskStatusResponse, error) {
	sherlockcrawlerService := server.Dependency.Crawler()
	if sherlockcrawlerService == nil {
		return nil, errors.New("crawler is not alive")
	}
	inCrawler := &crawlerproto.TaskStatusRequest{}
	responseCrawler, errCrawler := sherlockcrawlerService.StatusOfTaskQueue(context, inCrawler)
	return responseCrawler, errCrawler
}

/*
Will return a new Meta Array.
*/
func newMetaArray() *MetaArray {
	return &MetaArray{metamap: make(map[string]interface{})}
}

/*
fillMetaArray will put the response of the functions: receiveStatusTaskQAnalyser and receiveStatusTaskQCrawler in a MetaArra.
Help function of ReceiveMetadata.
*/
func fillMetaArray(responseCrawler *crawlerproto.TaskStatusResponse, responseAnalyser *analyserproto.WorkloadResponse) MetaArray {
	metaarray := newMetaArray()
	metaarray.metamap["Crawler"] = gin.H{
		"Website":    responseCrawler.Website,
		"Undone":     responseCrawler.Undone,
		"Processing": responseCrawler.Processing,
		"Finished":   responseCrawler.Finished,
		"Failed":     responseCrawler.Failed,
	}

	metaarray.metamap["Analyser"] = gin.H{
		"Website":       responseAnalyser.CrawledWebsite,
		"Undone":        responseAnalyser.Undone,
		"Processing":    responseAnalyser.Processing,
		"CrawlerError":  responseAnalyser.CrawlerError,
		"Saving":        responseAnalyser.Saving,
		"SendToCrawler": responseAnalyser.SendToCrawler,
		"Finished":      responseAnalyser.Finished,
	}
	return *metaarray
}

/*
ChangeState is used to change the state of the analyser/crawler service.
Can send the new status to one of them or both at once.
Targeted Service and the new status is transmitted via the post request.
*/
func (server *SherlockWebserver) ChangeState(context *gin.Context) {
	var status = NewRequestedStatus()
	err := context.BindJSON(status)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Error while reveiving Requested Status",
		})
	}

	statusToAnalyser, statusToCrawler, statuserr := setRequestedStatus(status)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"Status": statuserr,
		})
	}

	switch status.Target {
	case "Crawler":
		server.sendStateToCrawler(statusToCrawler, context)
	case "Analyser":
		server.sendStateToAnalyser(statusToAnalyser, context)
	case "All":
		server.sendStateToCrawlerAndAnalyser(statusToCrawler, statusToAnalyser, context)
	default:
		context.JSON(http.StatusBadRequest, gin.H{
			"Status": "Unknown Target, expected Crawler, Analyser or All",
		})
	}
}

/*
sendStateToAnalyser sends the new state to the analyser.
Help function of ChangeState.
*/
func (server *SherlockWebserver) sendStateToAnalyser(status analyserproto.AnalyserStateEnum, context *gin.Context) {
	stateanalyserService := server.Dependency.Analyser()
	inAnalyser := &analyserproto.ChangeStateRequest{State: status}
	_, errAnalyser := stateanalyserService.ChangeStateRPC(context, inAnalyser)

	if errAnalyser == nil {
		context.JSON(http.StatusOK, gin.H{
			"Status": "Fine",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"Analyser": "Unknown",
		})
	}
}

/*
sendStateToCrawler sends the new state to the crawler.
Help function of ChangeState.
*/
func (server *SherlockWebserver) sendStateToCrawler(status crawlerproto.CurrentState, context *gin.Context) {
	sherlockcrawlerService := server.Dependency.Crawler()
	inCrawler := &crawlerproto.StateRequest{State: status}
	_, errCrawler := sherlockcrawlerService.SetState(context, inCrawler)

	if errCrawler == nil {
		context.JSON(http.StatusOK, gin.H{
			"Status": "Fine",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"Crawler": "Unknown",
		})
	}
}

/*
sendStateToCrawlerAndAnalyser sends the new state to the crawler and analyser.
Help function of ChangeState.
*/
func (server *SherlockWebserver) sendStateToCrawlerAndAnalyser(statusCrawler crawlerproto.CurrentState, statusAnalyser analyserproto.AnalyserStateEnum, context *gin.Context) {
	sherlockcrawlerService := server.Dependency.Crawler()
	inCrawler := &crawlerproto.StateRequest{State: statusCrawler}
	_, errCrawler := sherlockcrawlerService.SetState(context, inCrawler)

	stateanalyserService := server.Dependency.Analyser()
	inAnalyser := &analyserproto.ChangeStateRequest{State: statusAnalyser}
	_, errAnalyser := stateanalyserService.ChangeStateRPC(context, inAnalyser)

	switch {
	case errAnalyser == nil && errCrawler == nil:
		context.JSON(http.StatusOK, gin.H{
			"Status": "Fine",
		})
	case errAnalyser != nil && errCrawler != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": "Unknown",
			"Crawler":  "Unknown",
		})
	case errAnalyser != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": "Unknown",
		})
	default:
		context.JSON(http.StatusOK, gin.H{
			"Crawler": "Unknown",
		})
	}
}

/*
setRequestedStatus will determine the requested Status and prepare the content of proto message.
Help function of ChangeState.
*/
func setRequestedStatus(requestedStatus *RequestedStatus) (statusAnalyser analyserproto.AnalyserStateEnum, statusCrawler crawlerproto.CurrentState, err error) {

	switch requestedStatus.Operation {
	case "Stop":
		statusAnalyser = analyserproto.AnalyserStateEnum_Stop
		statusCrawler = crawlerproto.CurrentState_Stop
	case "Pause":
		statusAnalyser = analyserproto.AnalyserStateEnum_Pause
		statusCrawler = crawlerproto.CurrentState_Pause
	case "Resume":
		statusAnalyser = analyserproto.AnalyserStateEnum_Running
		statusCrawler = crawlerproto.CurrentState_Running
	case "Clear": // JW
		statusAnalyser = analyserproto.AnalyserStateEnum_Clean
		statusCrawler = crawlerproto.CurrentState_Clean
	default:
		err = errors.New("unknown Operation, cant change state of service/services")
	}

	return statusAnalyser, statusCrawler, err
}

/*
GetServiceStatus returns the status of the crawler/analyser service. (checks if they are running).
*/
func (server *SherlockWebserver) GetServiceStatus(context *gin.Context) {

	responseAnalyser, errAnalyser := server.getStateFromAnalyser(context)
	responseCrawler, errCrawler := server.getStateFromCrawler(context)

	switch {
	case errAnalyser != nil && errCrawler != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": "Unknown",
			"Crawler":  "Unknown",
		})
	case errAnalyser != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": "Unknown",
			"Crawler":  responseCrawler.GetState().String(),
		})
	case errCrawler != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": responseAnalyser.GetState().String(),
			"Crawler":  "Unknown",
		})
	default:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": responseAnalyser.GetState().String(),
			"Crawler":  responseCrawler.GetState().String(),
		})
	}
}

/*
getStateFromAnalyser gets the current state from the analyser service.
Help function of GetServiceStatus.
*/
func (server *SherlockWebserver) getStateFromAnalyser(context *gin.Context) (*analyserproto.StateResponse, error) {
	stateanalyserService := server.Dependency.Analyser()
	inAnalyser := &analyserproto.StateRequest{}
	responseAnalyser, errAnalyser := stateanalyserService.StateRPC(context, inAnalyser)
	return responseAnalyser, errAnalyser
}

/*
getStateFromCrawler gets the current state from the crawler service.
Help function of GetServiceStatus.
*/
func (server *SherlockWebserver) getStateFromCrawler(context *gin.Context) (*crawlerproto.StateGetResponse, error) {
	sherlockcrawlerService := server.Dependency.Crawler()
	inCrawler := &crawlerproto.StateGetRequest{}
	responseCrawler, errCrawler := sherlockcrawlerService.GetState(context, inCrawler)
	return responseCrawler, errCrawler
}

/*
DropGraphTable should drop the neo4j table.
*/
func (server *SherlockWebserver) DropGraphTable(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database",
		})
	} else {
		_, err := sherlockneo.DropTable(session)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"Message": "A Problem occurred while trying to drop the Database",
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"Message": "Dropped the table.",
			})
		}
	}
}
