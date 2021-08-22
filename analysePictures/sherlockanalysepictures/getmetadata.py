#!/usr/bin/env python3

import os,sys,requests,psycopg2
from .sherlockneo import GetImages
from exif import Image
from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
from sqlalchemy import Column, Integer, String, Date, Boolean

Base = declarative_base()

# sets host ip, neccessary for docker and local env
def readFromENV(key, defaultVal):
    try:
        value = os.environ[key]
        if value is None or value == "":
            return defaultVal
        return value
    except:
        return defaultVal

# prepares url
def init():
    add = readFromENV("FLASKA_URL", "0.0.0.0")
    uri = "postgresql://gopher:gopher@" + add + ":5432/metadata"
    return uri     

DATABASE_URI = init()
engine = create_engine(DATABASE_URI)
Session = sessionmaker(bind=engine)

# DB class
class Mdata(Base):
    __tablename__ = 'metadata'
    neo4j_node_id = Column(Integer, primary_key = True)
    img_url = Column(String(500))
    datetime_original = Column(String(50))  
    model = Column(String(50))
    make = Column(String(50))
    maker_note = Column(String(50))
    software = Column(String(50))
    gps_latitude = Column(String(50))
    gps_longitude = Column(String(50))

    def __init__(self, neo4j_node_id, img_url, datetime_original, model, make, maker_note, software, gps_latitude, gps_longitude):
        self.neo4j_node_id = neo4j_node_id
        self.img_url = img_url
        self.datetime_original = datetime_original
        self.model = model
        self.make = make
        self.maker_note = maker_note
        self.software = software
        self.gps_latitude = gps_latitude
        self.gps_longitude = gps_longitude


# create database
def databaseCreateTable():
    Base.metadata.create_all(engine)

# drops database
def databaseDeleteTable():
    Base.metadata.drop_all(engine)
 
# inserts exif data in database
def databaseInsertData(node_id, listWithExif, imgurl):
    s = Session()
    exists = s.query(Mdata.neo4j_node_id).filter_by(neo4j_node_id=node_id).first() is not None
    if exists:
        return
    else:   
        latitude = "{}".format(listWithExif[5])
        longtitude = "{}".format(listWithExif[6])
        object1 = Mdata(node_id, imgurl, listWithExif[0], listWithExif[1], listWithExif[2], listWithExif[3], listWithExif[4], latitude, longtitude)
        s.add(object1)
        s.commit()
        s.close()

# return all entries from dp
def DatabaseRetreiveData():
    s = Session()
    res = s.query(Mdata).all()
    s.close()
    return res

# downloads all images that are currently in the neo4j database, analyses their exif data, then saves them in the postgres db
def DownloadImage():

    filePathToImage = "/tmp/tmp"
    relevantExifTags = ["datetime_original", "model", "make", "maker_note", "software", "gps_latitude", "gps_longitude"]
    listExifTags = []

    listWithIdAndUrl = GetImages() 
    #for test purpose
    #listWithIdAndUrl = [(6, "https://www.aboutbenita.com/wp-content/uploads/benita-thenhaus-body.jpg")]

    # check list
    if listWithIdAndUrl is not None and isinstance(listWithIdAndUrl, list) and len(listWithIdAndUrl) != 0:
        for pair in listWithIdAndUrl:

            response = requests.get(pair[1])
            _, fileExtension = os.path.splitext(pair[1])
            pathplusext = filePathToImage+fileExtension
            file = open(pathplusext, "wb")
            file.write(response.content)
            file.close()
            
            #get metadata and save the result in a postgresql database
            with open(pathplusext, 'rb') as imageFile:
                myImage = Image(imageFile)
                availableExifTags = myImage.list_all()
                for currentRelevantExifTag in relevantExifTags:
                    if currentRelevantExifTag in availableExifTags:
                        listExifTags.append((myImage.get(currentRelevantExifTag)))
                    else:
                        listExifTags.append("-")
                databaseInsertData(pair[0], listExifTags, pair[1])