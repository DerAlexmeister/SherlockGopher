package test

import (
	"net/url"

	neo4j "github.com/neo4j/neo4j-go-driver/neo4j"
)

/*
TDriver.
*/
type TDriver struct {
}

/*
Session.
*/
func (driver TDriver) Session(accessMode neo4j.AccessMode, bookmarks ...string) (neo4j.Session, error) {
	return nil, nil
}

/*
Target.
*/
func (driver TDriver) Target() url.URL {
	return url.URL{}
}

/*
Close.
*/
func (driver TDriver) Close() error {
	return nil
}

/*
GetNeo4jDriverInstance.
*/
func GetNeo4jDriverInstance() neo4j.Driver {
	return TDriver{}
}
