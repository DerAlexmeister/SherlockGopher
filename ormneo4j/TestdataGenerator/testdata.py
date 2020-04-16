#!/usr/bin/env python3

'''
Simple skript to produce testdata and put it in the neo4j db in order to 
support delevopment of the webserverservice and the frontendservice.
'''

# TODO 
# Add the Header

import random
import os
import sys

from neo4j import GraphDatabase

URL = "bolt://localhost:7687"
USER = "neo4j"
PASS = "test"

dropme = False

statements = {
    "constrains": "CREATE CONSTRAINT ON (c:Website) ASSERT c.address IS UNIQUE",
    "create" :  "UNWIND {props} as prop CREATE (a:Website {address:prop.address, statuscode:prop.statuscode, responsetime:prop.responsetime, Header:prop.header, status:prop.status});",
    "merge" : "MATCH (f:Website), (s:Website) WHERE f.address = {} AND s.address = {} MERGE (f)-[r:Links]->(s);",
    "drop": "MATCH (n) DETACH DELETE n",
}

class Neo4JInstance():
    ''' Neo4J Class representing a neo4j instance.  '''

    def __init__(self, uri, user, password):
        ''' CTor for neo4jinstance. '''
        try:
            self._driver = GraphDatabase.driver(uri, auth=(user, password), encrypted=False)
        except Exception as error:
            printError("__init__", error)

    def close(self):
        ''' Method to close the driver to the db. '''
        self._driver.close()

def findStatement(statement):
    ''' Function to find a querystatement in the dict. '''
    try:
        return statements[str(statement)]
    except Exception as error:
        printError("findStatement", error)

def getURL(): 
    ''' Getter for the URL in Neo4j '''
    return URL

def getUSER():
    ''' Getter for the USER of Neo4j '''
    return USER

def getPASS():
    ''' Getter for the PASS of Neo4j '''
    return PASS

def printmessage(message): 
    ''' Printmessage will print a message in a special format. '''
    print("[+] {}".format(message))

def printError(function, error): 
    ''' Printerror will print an error in a special format and will exit after the print'''
    print("[-] An error occurred in the function {} -  Error:  {}".format(function, error))
    sys.exit(0)

def getRandomNumber(start, stop):
    ''' 
    GetRandomNumber will return a random 
    number by a given start (inclusive) and stop (exclusive) and stepsize of one. 
    '''
    return random.randrange(start, stop)

def createNodes(tx, args):
    ''' CreateNodes function to execute the creation and transmit/commit it into the neo4jdb. '''
    tx.run(findStatement("create"), props=args)

def add_constraints(tx):
    ''' add_constraints function add constrains like the primarykey.'''
    tx.run(findStatement("constrains"))

def addConstrainsToDB(driver):
    ''' Function to run the neo4j function add_constraints '''
    with driver.session() as session:
        session.write_transaction(add_constraints)

def addNewNodes(driver):
    ''' addNewNodes will pseudorandom add some nodes '''
    for i in range(0, 100):
        with driver.session() as session:
            if i % 4 == 0: # unverified node 
                session.write_transaction(createNodes,{
                    "address": i,
                    "statuscode": getRandomNumber(200, 505),
                    "responsetime": getRandomNumber(1, 2000),
                    "status": "unverified"

                })
            elif i % 10 == 0: # error node
                session.write_transaction(createNodes, {
                    "address": i,
                    "statuscode": "An error occurred while trying to get this website. Error: ----",
                    "responsetime": getRandomNumber(1, 2000),
                    "status": "verified"

                })
            else: #normal node
                session.write_transaction(createNodes,{
                    "address": i,
                    "statuscode": getRandomNumber(200, 505),
                    "responsetime": getRandomNumber(1, 2000),
                    "status": "verified"

                })

def create_relationships(tx, args):
    ''' create_relationships function merge exisiting nodes in order to create relationships.'''
    tx.run(findStatement("merge").format(args['a'],args['b']))

def addRelationshipBetweenNodes(driver): 
    ''' Function to execute the create_relationships to create pseudo random relationships. '''
    try:
        for i in range(0, 2000):
            with driver.session() as session:
                session.write_transaction(create_relationships,
                {   
                    "a":str(getRandomNumber(0, 101)),
                    "b":str(getRandomNumber(0, 101)),
                })
                session.write_transaction(create_relationships,
                {
                    "a": str(i),  
                    "b":str(getRandomNumber(0, 101)),
                })
                print("Entry number {} has now a relationship".format(i))
    except Exception as error:
        printError("addRelationshipBetweenNodes", error)

def dropTable(tx):
    ''' Function to drop the entire DB. '''
    tx.run(findStatement("drop"))

def dropDB(driver):
    ''' Will drop the whole table. BE CAREFUL '''
    try:
        with driver.session() as session: session.write_transaction(dropTable)
    except Exception as error:
        printError("dropDB", error)

def runStatements(driver):
    global dropme
    try:
        if not dropme:
            printmessage("Starting to fill in the Database")
            addConstrainsToDB(driver)
            addNewNodes(driver)
            addRelationshipBetweenNodes(driver)
            printmessage("Finished - Database has now 100 nodes and 2000 relationships.")
        else: 
            dropDB(driver)
            printmessage("Dropped the entire Database.")
    except Exception as error:
        printError("runStatements", error)
            
def main():
    ''' Skript entrypoint. '''
    neo4jinst = Neo4JInstance(getURL(), getUSER(), getPASS())
    runStatements(neo4jinst._driver)
    

if __name__ == "__main__":
    main()