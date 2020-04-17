package sherlockneo

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	//ADDRESS of the Neo4j Dockercontainer.
	ADDRESS string = "bolt://localhost:7687"

	//USER will be the user of the db.
	USER string = "neo4j" //Standard username change this in production.

	//PASSWORD will be the password of the neo4j db.
	PASSWORD string = "test" //Standard password change this in production.
)

const (
	//Cypherstatements.

	//DROPGRAPH the entire graph.
	DROPGRAPH string = "MATCH (n) DETACH DELETE n"

	//CONTAINS will check whether or not a node is already in use.
	CONTAINS string = "MATCH (x) WHERE x.name = \"%s\" RETURN x"

	//ADDNODE will add a node into the neo4j DB.
	ADDNODE string = "UNWIND {props} as prop CREATE (a:Website {address:prop.address, statuscode:prop.statuscode, responsetime:prop.responsetime, Header:prop.header, status:prop.status});"

	//CONSTRAINS will be the constrains for a DB.
	CONSTRAINS string = "CREATE CONSTRAINT ON (c:Website) ASSERT c.address IS UNIQUE"

	//RETURNALL will return all nodes and their relationsships in the db.
	RETURNALL string = "MATCH (n) RETURN n"

	//CONNECTBYLINK will connect two given nodes by a link relationship.
	CONNECTBYLINK string = "MATCH (f:Website), (s:Website) WHERE f.address = \"%s\" AND s.address = \"%s\" MERGE (f)-[:]->(s);"

	//STARTERKIDOFNODE Will return a subset of nodes connected directly to a given node.
	STARTERKIDOFNODE string = "MATCH (a)-[:]->(b) WHERE a.address = \"%s\" RETURN a, b"
	// Vlt. mit Limit.
)

/*
getStarterKidOfNode will return the query for a small subset of nodes.
*/
func getStarterKidOfNode() string {
	return STARTERKIDOFNODE
}

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
RunStatement will execute a given statement on a given session and return the result.
*/
func RunStatement(session *neo4j.Session, statment string, args map[string]interface{}) (neo4j.Result, error) {
	result, err := (*session).Run(statment, args)
	if err != nil {
		return nil, fmt.Errorf("An error occured while trying to run the cypherstatement. Error: %s", err)
	}
	return result, nil
}

/*
GetNeoSnipped will return a subset of nodes and relationships.
Submit a target as parameter and the result will be the subset.
*/
func GetNeoSnipped(session *neo4j.Session, target string) (neo4j.Result, error) {
	args := make(map[string]interface{})
	result, err := RunStatement(session, fmt.Sprintf(getStarterKidOfNode(), target), args)
	if err != nil {
		return nil, fmt.Errorf("An error occurred while trying to execute the StarterKidOfNode statement. Error: %s", err.Error())
	}
	return result, nil
}

//func DropEntireGFraph() (neo4j.Result, error)
