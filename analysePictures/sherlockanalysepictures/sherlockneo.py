#!/usr/bin/env python3

from neo4j import GraphDatabase

URL = "bolt://10.0.2.15:7687" #"bolt://0.0.0.0:7687"
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
    print(a, type(a))
    return a

def runStatements(driver):
    ''' Will run all statements. '''
    try:
        with driver.session() as session:
            session.read_transaction(Images)
    except Exception as error:
        print("runStatements", error)
            
def GetImages():
    ''' Skript entrypoint. '''
    neo4jinst = Neo4JInstance(getURL(), getUSER(), getPASS())
    runStatements(neo4jinst._driver)