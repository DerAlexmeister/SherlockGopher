package webserver

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	sherlockneo "github.com/ob-algdatii-20ss/SherlockGopher/sherlockneo"
)

/*
GraphFetchWholeGraphHighPerformanceV1 will be a high performance endpoint to get optimized json for the Frontend.
*/
func (server *SherlockWebserver) GraphFetchWholeGraphHighPerformanceV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	} else {
		args := make(map[string]interface{})
		graph, _ := sherlockneo.GetAllNodesAndTheirRelationshipsOptimized(session, args)
		context.JSON(http.StatusOK, graph)
	}
}

/*
GraphFetchWholeGraphV1 will fetch the entire graph.
*/
func (server *SherlockWebserver) GraphFetchWholeGraphV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	} else {
		args := make(map[string]interface{})
		graph, _ := sherlockneo.GetAllNodesAndTheirRelationships(session, args)
		context.JSON(http.StatusOK, graph)
		defer sherlockneo.CloseSession(&session)
	}
}

/*
GraphMetaV1 will return all metainformation in json format.
*/
func (server *SherlockWebserver) GraphMetaV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	}
	args := make(map[string]interface{})
	images, _ := sherlockneo.GetAmountOfImages(session, args)
	css, _ := sherlockneo.GetAmountOfStylesheets(session, args)
	js, _ := sherlockneo.GetAmountOfJavascriptFiles(session, args)
	html, _ := sherlockneo.GetAmountofHTMLNodes(session, args)
	rels, _ := sherlockneo.GetAmountOfRels(session, args)
	nodes, _ := sherlockneo.GetAmountOfNodes(session, args)

	var meta []map[string]int64
	meta = append(meta, nodes[0], rels[0], html[0], css[0], js[0], images[0])
	context.JSON(http.StatusOK, meta)
	defer sherlockneo.CloseSession(&session)
}

/*
GraphPerformanceOfSitesV1 will return the performance of all sites like statuscode and RTT.
*/
func (server *SherlockWebserver) GraphPerformanceOfSitesV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	} else {
		args := make(map[string]interface{})
		performance, _ := sherlockneo.GetPerformanceOfSite(session, args)
		context.JSON(http.StatusOK, performance)
		defer sherlockneo.CloseSession(&session)
	}
}

/*
GraphNodeDetailsV1 will receive a URL and return the properties of the node incase
the node exists.
*/
func (server *SherlockWebserver) GraphNodeDetailsV1(context *gin.Context) {
	session, sessionErr := sherlockneo.GetSession(server.Driver)

	if sessionErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	}

	var url = NewRequestedURL()
	bindErr := context.BindJSON(url)

	if bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"Message": "Malformed JSON", //TODO improve message
		})
	}

	validURL := govalidator.IsURL(url.URL)
	if !validURL {
		context.JSON(http.StatusBadRequest, gin.H{
			"Message": "Malformed URL", //TODO improve message
		})
	}

	if sessionErr == nil && bindErr == nil && validURL {
		details, err := sherlockneo.GetDetailsOfNode(session, url.URL) //TODO
		fmt.Println(details)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"Message": "Sherlockneo Error", //TODO improve message
			})
		} else {
			context.JSON(http.StatusOK, details)
		}
		defer sherlockneo.CloseSession(&session)
	}
}
