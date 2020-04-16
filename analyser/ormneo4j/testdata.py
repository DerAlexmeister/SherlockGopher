#!/bin/python3

'''
Simple skript to produce testdata and put it in the neo4j db in order to 
support delevopment of the webserverservice and the frontendservice.
'''

from neo4j import GraphDatabase

def printmessage(message): 
    ''' Printmessage will print a message in a special format. '''
    print("[+] {}".format(message))


def main():
    ''' Skript entrypoint. '''
    printmessage("Starting to fill in the Database")