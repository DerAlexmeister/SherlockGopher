package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sherlockneo "github.com/ob-algdatii-20ss/SherlockGopher/sherlockneo"
)

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
	css, _ := sherlockneo.GetAmountOfImages(&session, args)
	js, _ := sherlockneo.GetAmountOfImages(&session, args)
	html, _ := sherlockneo.GetAmountOfImages(&session, args)
	rels, _ := sherlockneo.GetAmountOfImages(&session, args)
	nodes, _ := sherlockneo.GetAmountOfImages(&session, args)

	var meta [][]map[string]int64
	meta = append(meta, images, css, js, html, rels, nodes)
	context.JSON(http.StatusOK, meta)
}

//func (server *SherlockWebserver) GraphMetaV1(context *gin.Context)
