#!/usr/bin/env python3

import os
from neo4j import GraphDatabase

def readFromENV(key, defaultVal):
    try:
        value = os.environ[key]
        if value is None or value == "":
            return defaultVal
        return value
    except:
        return defaultVal

def init():
    add = readFromENV("FLASKA_URL", "0.0.0.0")
    uri = "bolt://" + add + ":7687"
    return uri     

URL =  init()
USER = "neo4j"
PASS = "test"

statements = {
    "Images" : "Match (x:Image) Return x.Address as Address, id(x) as NodeID"
}

class Neo4JInstance():
    ''' Neo4J Class representing a neo4j instance.  '''

    def __init__(self, uri, user, password):
        ''' CTor for neo4jinstance. '''
        try:
            self._driver = GraphDatabase.driver(uri, auth=(user, password), encrypted=False)
        except Exception as error:
            print("__init__", error)

    def close(self):
        ''' Method to close the driver to the db. '''
        self._driver.close()

def findStatement(statement):
    ''' Function to find a querystatement in the dict. '''
    try:
        return statements[str(statement)]
    except Exception as error:
        print("findStatement", error)

def getURL(): 
    ''' Getter for the URL in Neo4j '''
    return URL

def getUSER():
    ''' Getter for the USER of Neo4j '''
    return USER

def getPASS():
    ''' Getter for the PASS of Neo4j '''
    return PASS

def Images(tx):
    a = tx.run(findStatement("Images"))
    a = [(record['NodeID'], record['Address']) for record in a]
    #print(a, type(a))
    return a

def runStatements(driver):
    ''' Will run all statements. '''
    try:
        with driver.session() as session:
            a = session.read_transaction(Images)
            return a
    except Exception as error:
        print("runStatements", error)
            
def GetImages():
    ''' Skript entrypoint. '''
    neo4jinst = Neo4JInstance(getURL(), getUSER(), getPASS())
    a = runStatements(neo4jinst._driver)
    return a