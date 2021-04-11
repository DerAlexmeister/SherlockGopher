package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	sherlockneo "github.com/DerAlexx/SherlockGopher/sherlockneo"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type RequestedAddress struct {
	Address string `json:"address" binding:"required"`
}

/*
NewRequestedAddress will be a new instance of RequestedAddress.
*/
func NewRequestedAddress() *RequestedAddress {
	return &RequestedAddress{}
}

/*
GraphFetchWholeGraphHighPerformanceV1 will be a high performance endpoint to get optimized json for the Frontend.
*/
func (server *SherlockWebserver) GraphFetchWholeGraphHighPerformanceV1(context *gin.Context) {
	param := context.Param("query")
	paramtoint, err := strconv.Atoi(param)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"Status": "Path was malformed",
		})
	}
	var query string
	var add = NewRequestedAddress()
	switch paramtoint {
	case 0:
		query = sherlockneo.GetAllRels()
	case 1:
		query = sherlockneo.GetRoot()
	case 2:
		err := context.BindJSON(add)
		if err != nil || len(add.Address) == 0 {
			context.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Error while reveiving Requested Address",
			})
		}
		query = fmt.Sprintf(sherlockneo.GetChildren(), add.Address)
		fmt.Println(add.Address)

	default:
		query = sherlockneo.GetAllRels()
	}

	session, err := sherlockneo.GetSession(server.Driver)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database", //TODO improve message
		})
	} else {
		args := make(map[string]interface{})
		graph, _ := sherlockneo.GetAllNodesAndTheirRelationshipsOptimized(session, args, query)
		if paramtoint == 2 {
			for k, v := range graph["links"] {
				if v["target"] == add.Address {
					removeFromSlice(graph["links"], k)
				}
			}
			for k, v := range graph["nodes"] {
				if v["id"] == add.Address {
					removeFromSlice(graph["nodes"], k)
				}
			}
		}
		context.JSON(http.StatusOK, graph)
	}
}

func removeFromSlice(slice []map[string]string, s int) []map[string]string {
	return append(slice[:s], slice[s+1:]...)
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
	var meta []map[string]int64
	if images, _ := sherlockneo.GetAmountOfImages(session, args); len(images) != 0 {
		meta = append(meta, images[0])
	}
	if css, _ := sherlockneo.GetAmountOfStylesheets(session, args); len(css) != 0 {
		meta = append(meta, css[0])
	}
	if js, _ := sherlockneo.GetAmountOfJavascriptFiles(session, args); len(js) != 0 {
		meta = append(meta, js[0])
	}
	if html, _ := sherlockneo.GetAmountofHTMLNodes(session, args); len(html) != 0 {
		meta = append(meta, html[0])
	}
	if rels, _ := sherlockneo.GetAmountOfRels(session, args); len(rels) != 0 {
		meta = append(meta, rels[0])
	}
	if nodes, _ := sherlockneo.GetAmountOfNodes(session, args); len(nodes) != 0 {
		meta = append(meta, nodes[0])
	}

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
