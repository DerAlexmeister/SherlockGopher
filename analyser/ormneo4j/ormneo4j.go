package ormneo4j

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
*/
func CloseDatabaseConnection(driver neo4j.Driver) {
	defer driver.Close()
}

/*
GetSession will return a session to operate on inorder to send and recieve data.
*/
func GetSession(driver neo4j.Driver) (neo4j.Session, error) {
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	return session, nil
}

/*
CloseSession will close a session to the DB.
*/
func CloseSession(session neo4j.Session) {
	defer session.Close()
}

/*
RunStatement will execute a given statement on a given session and return the result.
*/
func RunStatement(session neo4j.Session, statments string, args map[string]interface{}) {

}
