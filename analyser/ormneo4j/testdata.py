#!/usr/bin/env python3

'''
Simple skript to produce testdata and put it in the neo4j db in order to 
support delevopment of the webserverservice and the frontendservice.
'''

# TODO 
# Merge to nodes 
# Add the Header


import random
import os
import sys

from neo4j import GraphDatabase

URL = "bolt://localhost:7687"
USER = "neo4j"
PASS = "test"

statements = {
    "create" :  "UNWIND {props} as prop CREATE (a:Website {address:prop.address, statuscode:prop.statuscode, responsetime:prop.responsetime, Header:prop.header, status:prop.status})",
    "merge" : ""
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

#"Dont not how the header will work at this point :("

def createNodes(tx, args):
    tx.run(findStatement("create"), props=args)

def add_constraints(tx):
    tx.run("CREATE CONSTRAINT ON (c:Website) ASSERT c.address IS UNIQUE")

def addConstrainsToDB(driver):
     with driver.session() as session:
        session.write_transaction(add_constraints)


def addNewNodes(driver):
    for i in range(0, 10000):
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

def runStatements(driver):
    '''  '''
    try:
        addConstrainsToDB(driver)
        addNewNodes(driver)
    except Exception as error:
        printError("runStatements", error)
            
def main():
    ''' Skript entrypoint. '''
    printmessage("Starting to fill in the Database")
    neo4jinst = Neo4JInstance(getURL(), getUSER(), getPASS())
    #print(neo4jinst._driver)
    runStatements(neo4jinst._driver)

if __name__ == "__main__":
    main()