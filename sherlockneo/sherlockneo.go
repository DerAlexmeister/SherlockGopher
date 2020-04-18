package sherlockneo

import (
	"encoding/json"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// TODO
// - Querys GETTER
// - Function to turn a result into json and return it
// - Implementing the stuff defined in Graph
// - Add missing Querys from Python Skript
// - merge this branch and the webserver branch in order to build the rest api
// - Write the first tests
// - Mocking for neo4j with neoism possible?
// - Add querys for info like amout of nodes of type image, js, css etc.

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

/*
DropEntireGraph will drop the entrie graph. Be careful incase you use this one.
*/
func DropEntireGraph(session *neo4j.Session) (neo4j.Result, error) {
	args := make(map[string]interface{})
	result, err := RunStatement(session, fmt.Sprintf(getStarterKidOfNode()), args)
	if err != nil {
		return nil, fmt.Errorf("An error occurred while trying to execute the droptable statement. Error: %s", err.Error())
	}
	return result, nil
}

/*
JsonfiyNeo will try to turn a given neo4j result into json.
It will return a byte array containing the formated json-output.
If an err occurred the byte array is nil
*/
func JsonfiyNeo(res neo4j.Result) ([]byte, error) {
	sliceofrecords := []map[string]interface{}{}
	for res.Next() {
		recor := make(map[string]interface{})
		for _, element := range res.Record().Keys() {
			if value, contains := res.Record().Get(element); contains {
				recor[element] = value
			}
		}
		sliceofrecords = append(sliceofrecords, recor)
	}
	jsonData, err := json.Marshal(sliceofrecords)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
