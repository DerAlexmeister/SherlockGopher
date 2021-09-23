package webserver

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	analyserproto "github.com/DerAlexx/SherlockGopher/analyser/proto"
	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
	crawlerproto "github.com/DerAlexx/SherlockGopher/sherlockcrawler/proto"
	sherlockneo "github.com/DerAlexx/SherlockGopher/sherlockneo"
)

var postgresuri string

/*
Init prepares urls for postgres
*/
func Init() {
	tmp := readFromENV("POSTG_URL", "0.0.0.0")
	postgresuri = "host=" + tmp + " user=gopher password=gopher dbname=metadata port=5432"
}

/*
readFromENV allows docker usage
*/
func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

/*
SherlockWebserver will be the webserver being the man in the middle between the frontend and the backend.
*/
type SherlockWebserver struct {
	Dependency *SherlockWebServerDependency
	Driver     neo4j.Driver
	PGdriver   *gorm.DB
	MGdriver   *screenshot.DB
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
RequestedPagination  will be the struct for the post request of the metadata and screenshot service.
*/
type RequestedPagination struct {
	Map         []interface{} `json:"map" binding:"required"`
	Maxpage     int           `json:"maxpage" binding:"required"`
	CurrentPage int           `json:"currentpage" binding:"required"`
	Pagerange   int           `json:"pagerange" binding:"required"`
}

/*
MetaArray stores the JSON which contains the response of the crawler and the analyser after the webserver requests their status.
Help function of ReceiveMetadata.
*/
type MetaArray struct {
	metamap map[string]interface{}
}

type Metadata struct {
	neo4j_node_id     int
	img_url           string
	datetime_original string
	model             string
	make              string
	maker_note        string
	software          string
	gps_latitude      string
	gps_longitude     string
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
NewRequestedStatus will be a new instance of RequestedPagination.
*/
func NewRequestedPagination() *RequestedPagination {
	return &RequestedPagination{}
}

/*
New will return a new instance of the SherlockWebserver.
*/
func New() *SherlockWebserver {
	ldriver, _ := sherlockneo.GetNewDatabaseConnection()
	Init()
	pgdriver, _ := connectToPostgresDb()
	mgdriver := screenshot.Connect()
	return &SherlockWebserver{
		Driver:   ldriver,
		PGdriver: pgdriver,
		MGdriver: mgdriver,
	}
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

/*
DropMongoTable should drop the mongo table.
*/
func (server *SherlockWebserver) DropMongoTable(context *gin.Context) {
	err := server.MGdriver.DropDB()
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

/*
DropPostgresTable should drop the postgres table.
*/
func (server *SherlockWebserver) DropPostgresTable(context *gin.Context) {

	err := server.PGdriver.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Metadata{})

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

/*
isNil checks wheter a interface is nil or not
*/
func isNil(i interface{}) string {
	switch i.(type) {
	case nil:
		return "-"
	default:
		return i.(string)
	}
}

/*
getStartStopMaxPage is a help function for the frontend pagination
*/
func getStartStopMaxPage(showpersite int, page int, size int) (start int, stop int, maxpage int) {
	start = page * showpersite
	stop = start + showpersite

	if size%showpersite != 0 {
		maxpage = (size / showpersite) + 1
	} else {
		maxpage = (size / showpersite)
	}
	if page+1 == maxpage {
		start = page * showpersite
		stop = size
	}
	if page > maxpage || page < 0 {
		start = 0
		stop = showpersite
	}
	if size < showpersite {
		start = 0
		stop = size
	}
	return start, stop, maxpage
}

/*
buildRequestedPagination builds the response for the frontend
*/
func buildRequestedPagination(mapparam []interface{}, maxpage int, currentpage int) RequestedPagination {
	tmpstruct := NewRequestedPagination()
	tmppagerange := 0
	tmpstruct.Map = mapparam
	tmpstruct.Maxpage = maxpage
	tmpstruct.CurrentPage = currentpage
	if maxpage > 5 {
		tmppagerange = 5
	}
	tmpstruct.Pagerange = tmppagerange
	return *tmpstruct
}

/*
GetScreenshots gets al screenshot data from the mongo db, picks 25 entries depending on the current pagination page.
*/
func (server *SherlockWebserver) GetScreenshots(ctx *gin.Context) {
	imagespersite := 25
	allscreenshots, err := server.MGdriver.ReturnAllScreenshots()
	if err != nil || len(allscreenshots) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Error while receiving data from database",
		})
	}
	param := ctx.Param("page")
	paramtoint, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "Path was malformed",
		})
	}
	start, stop, maxpage := getStartStopMaxPage(imagespersite, paramtoint, len(allscreenshots))
	partofallscreenshots := allscreenshots[start:stop]
	var tmpmap []interface{}
	for k, v := range partofallscreenshots {
		imagename := strconv.Itoa(k) + ".png"
		path := "../images/" + imagename
		if err := ioutil.WriteFile(path, *(v.GetPicture()), 0644); err != nil {
			panic(err)
		}
		tmpmap = append(tmpmap, gin.H{
			"imagepath": imagename,
			"imageurl":  v.GetUrl(),
		})
	}
	res := buildRequestedPagination(tmpmap, maxpage, paramtoint)
	ctx.JSON(http.StatusOK, res)
}

func connectToPostgresDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(postgresuri), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

/*
GetMetaData gets all metadata entries from the postgres db, picks 10 entries depending on the current pagination page.
*/
func (server *SherlockWebserver) GetMetaData(ctx *gin.Context) {
	metadatapersite := 10
	param := ctx.Param("page")
	paramtoint, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "Path was malformed",
		})
	}

	// get all entries
	resultmap := []map[string]interface{}{}
	result := server.PGdriver.Table("metadata").Find(&resultmap)
	if result.Error != nil || result.RowsAffected == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Error while reveiving data from database",
		})
	}

	var tmpmeta []Metadata

	for i := range resultmap {
		tmp := Metadata{}
		tmp.neo4j_node_id = int(resultmap[i]["neo4j_node_id"].(int32))
		tmp.img_url = isNil(resultmap[i]["img_url"])
		tmp.datetime_original = isNil(resultmap[i]["datetime_original"])
		tmp.model = isNil(resultmap[i]["model"])
		tmp.make = isNil(resultmap[i]["make"])
		tmp.maker_note = isNil(resultmap[i]["maker_note"])
		tmp.software = isNil(resultmap[i]["software"])
		tmp.gps_longitude = isNil(resultmap[i]["gps_longitude"])
		tmp.gps_latitude = isNil(resultmap[i]["gps_latitude"])
		tmpmeta = append(tmpmeta, tmp)
	}

	start, stop, maxpage := getStartStopMaxPage(metadatapersite, paramtoint, int(result.RowsAffected))
	partofallmeta := tmpmeta[start:stop]
	var tmpmap []interface{}
	for _, v := range partofallmeta {
		tmpmap = append(tmpmap, gin.H{
			"neo4j_node_id":     v.neo4j_node_id,
			"img_url":           v.img_url,
			"datetime_original": v.datetime_original,
			"model":             v.model,
			"make":              v.make,
			"maker_note":        v.maker_note,
			"software":          v.software,
			"gps_latitude":      v.gps_latitude,
			"gps_longitude":     v.gps_longitude,
		})
	}
	res := buildRequestedPagination(tmpmap, maxpage, paramtoint)
	ctx.JSON(http.StatusOK, res)
}
