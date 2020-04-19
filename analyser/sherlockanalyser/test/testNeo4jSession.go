package test

import (
	neo4j "github.com/neo4j/neo4j-go-driver/neo4j"
)

type testSession struct {

}

func (session testSession) LastBookmark() string {
	return ""
}

func (session testSession) BeginTransaction(configurers ...func(*neo4j.TransactionConfig)) (neo4j.Transaction, error)  {
	return nil, nil
}

func (session testSession) ReadTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	return nil, nil
}

func (session testSession) WriteTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	return nil, nil
}

func (session testSession) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	return nil, nil
}

func (session testSession) Close() error {
	return nil
}

func GetNeo4jSessionInstance() neo4j.Session {
	return testSession{}
}
