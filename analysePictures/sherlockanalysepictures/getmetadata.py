#!/usr/bin/env python3

import os,sys,requests,psycopg2
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
    img_id = Column(Integer, primary_key = True)
    img_url = Column(String(50))
    condition = Column(Boolean)
    datetime_original = Column(String(50))  
    model = Column(String(50))
    make = Column(String(50))
    maker_note = Column(String(50))
    software = Column(String(50))
    gps_latitude = Column(String(50))
    gps_longitude = Column(String(50))

    def __init__(self, img_id, img_url, condition, datetime_original, model, make, maker_note, software, gps_latitude, gps_longitude):
        self.img_id = img_id
        self.img_url = img_url
        self.condition = condition
        self.datetime_original = datetime_original
        self.model = model
        self.make = make
        self.maker_note = maker_note
        self.software = software
        self.gps_latitude = gps_latitude
        self.gps_longitude = gps_longitude

    # for print debugging
    def __str__(self):
        return "{},{},{},{},".format(self.img_id, self.software, self.gps_latitude, self.gps_longitude)

# create database
def databaseCreateTable():
    Base.metadata.create_all(engine)

# drops database
def databaseDeleteTable():
    Base.metadata.drop_all(engine)
 
# inserts exif data in database
def databaseInsertData(img_id, cond, listWithExif, imgurl):
    s = Session()
    latitude = "{}".format(listWithExif[5])
    longtitude = "{}".format(listWithExif[6])
    object1 = Mdata(img_id, imgurl, cond, listWithExif[0], listWithExif[1], listWithExif[2], listWithExif[3], listWithExif[4], latitude, longtitude)

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

    listWithIdAndUrl = neo.GetImages() 
    #for test purpose
    #listWithIdAndUrl = [(6, "https://www.aboutbenita.com/wp-content/uploads/Foto-07.01.21-21-07-58-1.jpg")]

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
            if myImage.has_exif:
                availableExifTags = myImage.list_all()
                for currentRelevantExifTag in relevantExifTags:
                    if currentRelevantExifTag in availableExifTags:
                        listExifTags.append((myImage.get(currentRelevantExifTag)))
                    else:
                        listExifTags.append(None)
                databaseInsertData(pair[0], 1, listExifTags, pair[1])
            else:
                databaseInsertData(pair[0], 0, None, pair[1])