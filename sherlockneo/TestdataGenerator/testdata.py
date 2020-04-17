#!/usr/bin/env python3

'''
Simple skript to produce testdata and put it in the neo4j db in order to 
support delevopment of the webserverservice and the frontendservice.
'''

# TODO 
# Add the Header
# require als relationship
# 

import random
import os
import sys

from neo4j import GraphDatabase

URL = "bolt://localhost:7687"
USER = "neo4j"
PASS = "test"

dropme = False

#"UNWIND {props} as prop CREATE (c:Website {Address:prop.Address}) SET c += {props}",
#"create" :  "UNWIND {props} as prop CREATE (a:Address {Address:prop.Address, statuscode:prop.statuscode, responsetime:prop.responsetime, Header:prop.header, status:prop.status});",

statements = {
    "constrains": "CREATE CONSTRAINT ON (c:Website) ASSERT c.Address IS UNIQUE",
    "constrainsimg": "CREATE CONSTRAINT ON (c:Image) ASSERT c.Address IS UNIQUE",
    "constrainscss": "CREATE CONSTRAINT ON (c:StyleSheet) ASSERT c.Address IS UNIQUE",
    "constrainsjs": "CREATE CONSTRAINT ON (c:Javascript) ASSERT c.Address IS UNIQUE",
    "createadv" :  "UNWIND {props} as prop MERGE (c:Website {Address: prop.Address}) SET c += {props}",
    "createimg" :  "UNWIND {props} as prop MERGE (c:Image {Address: prop.Address}) SET c += {props}",
    "createcss" :  "UNWIND {props} as prop MERGE (c:StyleSheet {Address: prop.Address}) SET c += {props}",
    "createjs" :  "UNWIND {props} as prop MERGE (c:Javascript {Address: prop.Address}) SET c += {props}",
    "mergehtml" : "MATCH (f:Website), (s:Website) WHERE f.Address = {} AND s.Address = {} AND s.Filetype = \"HTML\" MERGE (f)-[r:Links]->(s);",
    "mergecss" : "MATCH (f:Website), (s:Javascript) WHERE f.Address = {} AND s.Address = {} AND s.Filetype = \"Javascript\" MERGE (f)-[r:Requires]->(s);",
    "mergeimg" : "MATCH (f:Website), (s:Image) WHERE f.Address = {} AND s.Address = {} AND s.Filetype = \"Image\" MERGE (f)-[r:Shows]->(s);",
    "mergejs" : "MATCH (f:Website), (s:StyleSheet) WHERE f.Address = {} AND s.Address = {} AND s.Filetype = \"CSS\" MERGE (f)-[r:Requires]->(s);",
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

def add_constraints(tx):
    ''' add_constraints function add constrains like the primarykey.'''
    tx.run(findStatement("constrains"))
    tx.run(findStatement("constrainscss"))
    tx.run(findStatement("constrainsimg"))
    tx.run(findStatement("constrainsjs"))

def addConstrainsToDB(driver):
    ''' Function to run the neo4j function add_constraints '''
    with driver.session() as session:
        session.write_transaction(add_constraints)

def createNodes(tx, args):
    ''' CreateNodes function to execute the creation and transmit/commit it into the neo4jdb. '''
    tx.run(findStatement("createadv"), props=args)

def createImg(tx, args):
    ''' CreateImg function to execute the creation and transmit/commit it into the neo4jdb just images. '''
    tx.run(findStatement("createimg"), props=args)

def createCss(tx, args):
    ''' CreateImg function to execute the creation and transmit/commit it into the neo4jdb just stylesheets. '''
    tx.run(findStatement("createcss"), props=args)

def createJs(tx, args):
    ''' CreateImg function to execute the creation and transmit/commit it into the neo4jdb just javascripts. '''
    tx.run(findStatement("createjs"), props=args)

def addNewNodes(driver):
    ''' addNewNodes will pseudorandom add some nodes '''
    headerfield = {
                    "Accept-Ranges": "bytes",
                    "Age": "12",
                    "Allow": "GET, HEAD",
                    "Cache-Control": "max-age=3600",
                    "Connection": "close",
                    "Content-Encoding": "gzip",
                    "Content-Language": "de",
                    "Content-Location": "/foo.html.de",
                    "Content-MD5": "Q2hlY2sgSW50ZWdyaXR5IQ==",
                    "Content-Type": "text/html; charset=utf-8",
                    "Expires": "Thu, 01 Dec 1994 16:00:00 GMT",
                    "Proxy-Authenticate": "Basic",
                }
    for i in range(0, 100):
        with driver.session() as session:
            if i % 4 == 0: # unverified node 
                props = {
                    "Address": i,
                    "Statuscode": getRandomNumber(200, 505),
                    "Responsetime": getRandomNumber(1, 2000),
                    "Status": "unverified",
                    "Filetype" : "HTML",
                }
                for elem in list(headerfield.items()): props[elem[0]] = elem[1]
                session.write_transaction(createNodes,props)
            elif i % 10 == 0: # error node
                session.write_transaction(createNodes, {
                    "Address": i,
                    "Statuscode": str(0),
                    "Responsetime": getRandomNumber(1, 2000),
                    "Status": "unverified",
                    "Filetype" : "HTML",
                    "Type" : "Error",
                })
            elif i % 7 == 0: # css
                props = {
                    "Address": i,
                    "Statuscode": getRandomNumber(200, 505),
                    "Responsetime": getRandomNumber(1, 2000),
                    "Status": "unverified",
                    "Filetype" : "CSS",
                }
                for elem in list(headerfield.items()): props[elem[0]] = elem[1]
                session.write_transaction(createCss,props)
            elif i % 3 == 0: # javascript
                props = {
                    "Address": i,
                    "Statuscode": getRandomNumber(200, 505),
                    "Responsetime": getRandomNumber(1, 2000),
                    "Status": "unverified",
                    "Filetype" : "Javascript",
                }
                for elem in list(headerfield.items()): props[elem[0]] = elem[1]
                session.write_transaction(createJs,props)
            if i % 8 == 0: # Img 
                props = {
                    "Address": i,
                    "Statuscode": getRandomNumber(200, 505),
                    "Responsetime": getRandomNumber(1, 2000),
                    "Status": "unverified",
                    "Filetype" : "Image",
                }
                for elem in list(headerfield.items()): props[elem[0]] = elem[1]
                session.write_transaction(createImg,props)
            else:                
                props = {
                    "Address": i,
                    "Statuscode": getRandomNumber(200, 505),
                    "Responsetime": getRandomNumber(1, 2000),
                    "Status": "unverified",
                    "Filetype" : "HTML",
                }
                for elem in list(headerfield.items()): props[elem[0]] = elem[1]
                session.write_transaction(createNodes, props)

def create_relationships(tx, args):
    ''' create_relationships function merge exisiting nodes in order to create relationships for html.'''
    tx.run(findStatement("mergehtml").format(args['a'],args['b']))

def create_relationships_css(tx, args):
    ''' create_relationships function merge exisiting nodes in order to create relationships for css.'''
    tx.run(findStatement("mergecss").format(args['a'],args['b']))

def create_relationships_img(tx, args):
    ''' create_relationships function merge exisiting nodes in order to create relationships for images.'''
    tx.run(findStatement("mergeimg").format(args['a'],args['b']))

def create_relationships_js(tx, args):
    ''' create_relationships function merge exisiting nodes in order to create relationships for javascripts.'''
    tx.run(findStatement("mergejs").format(args['a'],args['b']))

def getRandomWebsiteNumberAddress(start, stop):
    ''' 
    GetRandomNumber will return a random 
    number by a given start (inclusive) and stop (exclusive) and stepsize of one but file all things 
    which are not html so css/image/javascript. 
    '''
    x = random.randrange(start, stop)
    if x % 8 == 0 or x % 3 == 0 or x % 7 == 0: return getRandomWebsiteNumberAddress(start, stop)
    elif x is None: return getRandomWebsiteNumberAddress(start, stop)
    else: return x

def addRelationshipBetweenNodes(driver): 
    ''' Function to execute the create_relationships to create pseudo random relationships. '''
    try:
        printmessage("Creating HTML-Rels")
        for i in range(0, 2000):
            with driver.session() as session:
                session.write_transaction(create_relationships,
                {   
                    "a": str(getRandomNumber(0, 101)),
                    "b": str(getRandomNumber(0, 101)),
                })
                session.write_transaction(create_relationships,
                {
                    "a": str(i),  
                    "b": str(getRandomNumber(0, 101)),
                })
        printmessage("Creating CSS-Rels")
        for i in range(0, 500): #css
            with driver.session() as session:
                session.write_transaction(create_relationships_css,
                {   
                    "a": str(getRandomWebsiteNumberAddress(0, 101)),
                    "b": str(i), 
                })
        printmessage("Creating Image-Rels")
        for i in range(0, 250): #img
            with driver.session() as session:
                session.write_transaction(create_relationships_img,
                {
                    "a": str(getRandomWebsiteNumberAddress(0, 101)),
                    "b": str(i),
                })
        printmessage("Creating Javascript-Rels")
        for i in range(0, 1500): #js
            with driver.session() as session:
                session.write_transaction(create_relationships_js,
                {
                    "a": str(getRandomWebsiteNumberAddress(0, 101)),
                    "b": str(i),
                })
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
            printmessage("Finished - Database has now nodes and relationships.")
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