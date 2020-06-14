# Docs about the Webserver

The Webserver serves as Rest Endpoint. 
It provides a way to communicate between the frontend and the backend. (neo4j database, crawlerservice and the analyserservice)

To receive JSON from the frontend we bind it in structs like this:

```go
type RequestedURL struct {
	URL string `json:"url" binding:"required"`
}
```

We bind it like this:

```go
	var url = NewRequestedURL() // creates a RequestedURL struct
	err := ctx.BindJSON(url)    // binds the json in the struct
```

Sending JSON Back to the frontend would look like this:

```go
context.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Error while sending Metadata",
})
```

Implementation:
* Http web framework: gin-gonic
* Communication between the frontend and the webserver: json
* Communication between the webserver and the analyser/crawler service: GRPC

All Functions implemented in the webserver
* Helloping: used to check if the webserver is alive
* ChangeState: will change the state of the cralwer and the analyser
* GetServiceStatus: will return the status of the analyser/crawler service
* DropGraphTable: will drop the neo4j table
* ReceiveMetadata: get all meta information about the crawler and the analyser
* GraphMetaV1: get all meta information about neo4j
* GraphFetchWholeGraphV1
* GraphFetchWholeGraphHighPerformanceV1: will return the entire graph, maybe build a stream
* GraphPerformanceOfSitesV1: will return address with statuscode and reponsetime
* GraphNodeDetailsV1: get all information of a node
* ReceiveURL: will handle the requested url which should be crawled

For more information about the API, see the api.md in the docs or the Scripting-API.
