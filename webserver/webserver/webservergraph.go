package webserver

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	sherlockneo "github.com/ob-algdatii-20ss/SherlockGopher/sherlockneo"
)

//TODO close session after finishing the response.

/*
GraphFetchWholeGraphV1 will fetch the entire graph.
*/
func (server *SherlockWebserver) GraphFetchWholeGraphV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	}
	args := make(map[string]interface{})
	graph, _ := sherlockneo.GetAllNodesAndTheirRelationships(&session, args)
	context.JSON(http.StatusOK, graph)
}

/*
GraphMetaV1 will return all metainformation in json format.
*/
func (server *SherlockWebserver) GraphMetaV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	}
	args := make(map[string]interface{})
	images, _ := sherlockneo.GetAmountOfImages(&session, args)
	css, _ := sherlockneo.GetAmountOfStylesheets(&session, args)
	js, _ := sherlockneo.GetAmountOfJavascriptFiles(&session, args)
	html, _ := sherlockneo.GetAmountofHTMLNodes(&session, args)
	rels, _ := sherlockneo.GetAmountOfRels(&session, args)
	nodes, _ := sherlockneo.GetAmountOfNodes(&session, args)

	var meta [][]map[string]int64
	meta = append(meta, images, css, js, html, rels, nodes)
	context.JSON(http.StatusOK, meta)
}

/*
GraphPerformenceOfSitesV1 will return the performence of all sites like statuscode and RTT.
*/
func (server *SherlockWebserver) GraphPerformenceOfSitesV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	}
	args := make(map[string]interface{})
	performence, _ := sherlockneo.GetPerformenceOfSite(&session, args)

	var meta [][]map[string]string
	meta = append(meta, performence)
	context.JSON(http.StatusOK, meta)
}

/*
GraphNodeDetailsV1 will
*/
func (server *SherlockWebserver) GraphNodeDetailsV1(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)

	var url = NewRequestedURL()
	context.BindJSON(url)

	//check if url is empty or a well formed url.
	if govalidator.IsURL(url.URL) {
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
			})
		}
	}
	args := make(map[string]interface{})
	args["address"] = string(url.URL)
	performence, _ := sherlockneo.GetPerformenceOfSite(&session, args) //TODO
	context.JSON(http.StatusOK, performence)
}
