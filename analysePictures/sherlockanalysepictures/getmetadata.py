#!/usr/bin/env python3

import os,sys,requests,psycopg2
from sherlockneo import GetImages
from exif import Image

def databaseConnection():
    con =  psycopg2.connect(database="metadata", user="gopher", password="gopher", host="127.0.0.1", port="5432")
    return con

def databaseCreateTable():
    con = databaseConnection()
    cur = con.cursor()
    cur.execute('''CREATE TABLE METADATA
        (ID INT PRIMARY KEY     NOT NULL,
        CONDITION           BOOLEAN    NOT NULL,
        DATETIMEORIGINAL        CHAR(20),
        MODEL        CHAR(20),
        MAKE        CHAR(20),
        MAKERNOTE   CHAR(50),
        SOFTWARE    CHAR(20),
        GPSLATITUDE CHAR(20),
        GPSLONGITUDE    CHAR(20));''')
    con.commit()
    con.close()

def databaseDeleteTable():
    con = databaseConnection()
    cur = con.cursor()
    cur.execute('''DROP TABLE METADATA;''')
    con.commit()
    con.close()

def databaseInsertData(id, cond, listWithExif):
    con = databaseConnection()
    cur = con.cursor()
    cur.execute('''INSERT INTO METADATA (ID, CONDITION, DATETIMEORIGINAL, MODEL, MAKE, MAKERNOTE, SOFTWARE, GPSLATITUDE, GPSLONGITUDE) VALUES ({}, {}, {}, {}, {}, {}, {});'''.format(id, cond, listWithExif[0], listWithExif[1], listWithExif[2], listWithExif[3], listWithExif[4], listWithExif[5], listWithExif[6]))
    con.commit()
    con.close()

def DatabaseRetreiveData():
    relevantExifTags = ["id", "condition", "datetime_original", "model", "make", "maker_note", "software", "gps_latitude", "gps_longitude"]
    listExifTags = []
    con = databaseConnection()
    cur = con.cursor()
    cur.execute('''SELECT id, condition, datetimefromoriginal, model, make, makernote, software, gpslatitude, gpslongtitude from METADATA;''')
    rows = cur.fetchall()
    for i, row in enumerate(rows):
        listExifTags.append((relevantExifTags[i], row[i]))
    con.commit()
    con.close()
    return listExifTags

def DownloadImage():

    filePathToImage = "/tmp/tmp"
    relevantExifTags = ["datetime_original", "model", "make", "maker_note", "software", "gps_latitude", "gps_longitude"]
    listExifTags = []
    

    #list = neo.GetImages() 
    #liste = [(1, "https://walterj.de/hukafallsCR.jpg")]
    listWithIdAndUrl = [(1, "https://www.aboutbenita.com/wp-content/uploads/Juli-2020-800x1000.jpg")]

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

databaseDeleteTable()
databaseCreateTable()
DownloadImage()
a = DatabaseRetreiveData()
print (a)