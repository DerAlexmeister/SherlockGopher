package sherlockneo

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// TODO
// - Write the first tests
// - Mocking for neo4j with neoism possible?
// - Drop constrains
// - RunStatement for only write query from analyser to neo4j.

const (
	//ADDRESS of the Neo4j Dockercontainer.
	ADDRESS string = "bolt://localhost:7687"

	//USER will be the user of the db.
	USER string = "neo4j" //Standard username change this in production.

	//PASSWORD will be the password of the neo4j db.
	PASSWORD string = "test" //Standard password change this in production.
)

const (
	//Magicconstant.

	//EMPTYSTRING will return an emptystring.
	EMPTYSTRING string = ""
)

/*
Will return a emptystring.
*/
func getEmptyString() string {
	return EMPTYSTRING
}

/*
GetGraphDBAddress will return the Address of neo4j.
*/
func GetGraphDBAddress() string {
	return ADDRESS
}

/*
getUserName will return the username of the current user.
*/
func getUserName() string {
	return USER
}

/*
getPasswort will return the password of the current user.
*/
func getPasswort() string {
	return PASSWORD
}

/*
GetNewDatabaseConnection will return a new connection/driver for the current neo4j Instance.
*/
func GetNewDatabaseConnection() (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(GetGraphDBAddress(), neo4j.BasicAuth(getUserName(), getPasswort(), getEmptyString()))
	if err != nil {
		return nil, fmt.Errorf("An error occrued: %s", err)
	}
	return driver, nil
}

/*
CloseDatabaseConnection will close the drivers to the DB.
Eg. defer localneo.Close()
*/
func CloseDatabaseConnection(driver *neo4j.Driver) {
	(*driver).Close()
}

/*
GetSession will return a session to operate on inorder to send and recieve data.
*/
func GetSession(driver *neo4j.Driver) (neo4j.Session, error) {
	session, err := (*driver).Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	return session, nil
}

/*
CloseSession will close a session to the DB.
Eg. defer localneo.Close()
*/
func CloseSession(session *neo4j.Session) {
	(*session).Close()
}

/*
GetAllNodesAndTheirRelationships will return all nodes with address and the rels to other nodes.
*/
func GetAllNodesAndTheirRelationships(session *neo4j.Session, args map[string]interface{}) ([]map[string]string, error) {
	results, _ := (*session).Run(getAllRels(), args)
	var tojson []map[string]string
	for results.Next() {
		elem := make(map[string]string)
		for _, element := range results.Record().Keys() {
			if value, contains := results.Record().Get(element); contains {
				elem[element] = value.(string)
			}

		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetAmountOfNodes will return the amount of nodes
*/
func GetAmountOfNodes(session *neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	var tojson []map[string]int64
	amountofnodes, _ := (*session).Run(getCountNumberOfNodes(), args)

	for amountofnodes.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofnodes.Record().Keys() {
			if value, contains := amountofnodes.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetAmountOfRels will return the amout of relationship.
*/
func GetAmountOfRels(session *neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	var tojson []map[string]int64
	amountofrels, _ := (*session).Run(getCountRelsToNodes(), args)
	for amountofrels.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofrels.Record().Keys() {
			if value, contains := amountofrels.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetAmountofHTMLNodes will return the amount of html nodes.
*/
func GetAmountofHTMLNodes(session *neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	var tojson []map[string]int64
	amountofhtmls, _ := (*session).Run(getCountHtmlsNodes(), args)
	for amountofhtmls.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofhtmls.Record().Keys() {
			if value, contains := amountofhtmls.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetAmountOfStylesheets will return the amount of Stylesheets.
*/
func GetAmountOfStylesheets(session *neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	var tojson []map[string]int64
	amountofcss, _ := (*session).Run(getCountCSSNodes(), args)
	for amountofcss.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofcss.Record().Keys() {
			if value, contains := amountofcss.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetAmountOfJavascriptFiles will return the amount of Javascript files.
*/
func GetAmountOfJavascriptFiles(session *neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	var tojson []map[string]int64
	amountofjs, _ := (*session).Run(getCountJavascriptNodes(), args)
	for amountofjs.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofjs.Record().Keys() {
			if value, contains := amountofjs.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetAmountOfImages will return amount of images.
*/
func GetAmountOfImages(session *neo4j.Session, args map[string]interface{}) ([]map[string]int64, error) {
	var tojson []map[string]int64
	amountofimages, _ := (*session).Run(getCountImageNodes(), args)
	for amountofimages.Next() {
		elem := make(map[string]int64)
		for _, element := range amountofimages.Record().Keys() {
			if value, contains := amountofimages.Record().Get(element); contains {
				elem[element] = value.(int64)
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetPerformenceOfSite will return the performence index of each site like address, rTT and statuscode.
*/
func GetPerformenceOfSite(session *neo4j.Session, args map[string]interface{}) ([]map[string]string, error) {
	var tojson []map[string]string
	performence, _ := (*session).Run(getResponseTimeInTableAndStatusCode(), args)
	for performence.Next() {
		elem := make(map[string]string)
		for _, element := range performence.Record().Keys() {
			if value, contains := performence.Record().Get(element); contains {
				switch value.(type) {
				case int64:
					elem[element] = fmt.Sprintf("%d", value)
				case string:
					elem[element] = value.(string)
				}
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}

/*
GetDetailsOfNode will return information of a node given as parameter.
*/
func GetDetailsOfNode(session *neo4j.Session, args map[string]interface{}) ([]map[string]string, error) {
	var tojson []map[string]string
	details, _ := (*session).Run(getReturnNode(), args)
	for details.Next() {
		elem := make(map[string]string)
		for _, element := range details.Record().Keys() {
			if value, contains := details.Record().Get(element); contains {
				switch value.(type) {
				case int64:
					elem[element] = fmt.Sprintf("%d", value)
				case string:
					elem[element] = value.(string)
				}
			}
		}
		tojson = append(tojson, elem)
	}
	return tojson, nil
}
