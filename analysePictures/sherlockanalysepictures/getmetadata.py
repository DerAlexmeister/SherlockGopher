#!/usr/bin/env python3

import os,sys,requests,psycopg2
#from .sherlockneo import GetImages
from exif import Image
from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
from sqlalchemy import Column, Integer, String, Date, Boolean

Base = declarative_base()

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
    uri = "postgresql://gopher:gopher@" + add + ":5432/metadata"
    return uri     

DATABASE_URI = init()
engine = create_engine(DATABASE_URI)
Session = sessionmaker(bind=engine)

class Mdata(Base):
    __tablename__ = 'metadata'
    img_id = Column(Integer, primary_key = True)
    condition = Column(Boolean)
    datetime_original = Column(String(50))  
    model = Column(String(50))
    make = Column(String(50))
    maker_note = Column(String(50))
    software = Column(String(50))
    gps_latitude = Column(String(50))
    gps_longitude = Column(String(50))

    def __init__(self, img_id, condition, datetime_original, model, make, maker_note, software, gps_latitude, gps_longitude):
        self.img_id = img_id
        self.condition = condition
        self.datetime_original = datetime_original
        self.model = model
        self.make = make
        self.maker_note = maker_note
        self.software = software
        self.gps_latitude = gps_latitude
        self.gps_longitude = gps_longitude

    def __str__(self):
        return "{},{},{},{},".format(self.img_id, self.software, self.gps_latitude, self.gps_longitude)

def databaseCreateTable():
    Base.metadata.create_all(engine)

def databaseDeleteTable():
    Base.metadata.drop_all(engine)
 

def databaseInsertData(img_id, cond, listWithExif):
    s = Session()
    latitude = "{}".format(listWithExif[5])
    longtitude = "{}".format(listWithExif[6])
    object1 = Mdata(img_id, cond, listWithExif[0], listWithExif[1], listWithExif[2], listWithExif[3], listWithExif[4], latitude, longtitude)

    s.add(object1)
    s.commit()
    s.close()

def DatabaseRetreiveData():
    s = Session()
    res = s.query(Mdata).all()
    s.close()
    return res


def DownloadImage():

    filePathToImage = "/tmp/tmp"
    relevantExifTags = ["datetime_original", "model", "make", "maker_note", "software", "gps_latitude", "gps_longitude"]
    listExifTags = []
    

    listWithIdAndUrl = neo.GetImages() 
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
                databaseInsertData(pair[0], 1, listExifTags)
            else:
                databaseInsertData(pair[0], 0, None)