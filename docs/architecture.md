# Docs of the System architecture

<p align="center"><img alt="Sherlock Gopher" src="https://github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/blob/develop/assets/general/ServiceDest.png" width="550" height="550"></p>

## Neo4j

Neo4j is the persistens part of the system so like our brain. It stores all information collected in a huge graph. A API is provided via
the bolt protocol. To interact with the database, we build a package called sherlockneo which is a wrapper for neo4j and the needed function of the webserver and the analyser.

We have chosen neo4j for our project because it is kind of ideal for the purpose of storing a big Graph like we do by crawling and analysing a website.

## The Webserver

The Webserver is the mind behind SherlockGopher. Managing all services and providing data to the frontend. The webserver is connected via gRPC with crawler and the analyser, talks to neo4j via bolt and with the frontend via REST. So the main task is to controll the information-flow.

## The Frontend

As it is 2020 and SherlockGopher is a distributed system the Userinterface can not be a simple GUI. Especially the UI should by platform independent. So we built something more fitting the requirements. The solution is a modern webinterface using javascript or more specific with ReactJS. Also considering the fact that you wanna use your own Userinterface you can easily switch to a cli or something else by using the Skripting-API.

## The Crawler

Now we are looking at the eyes of our system. We can only see what the crawler saw. The crawler will visit a website and fetch all needed information, so that the analyser can work with it. After fetching the data the crawler will talk to the analyser via gRPC. But not in the normal boring way. No it uses something called Servicestreaming. This will resolve the Issue of too small packets in the gRPC protocol.

## The Analyser

Last but not least the analyser. This is the part of the system where the real magic happens. With a customized parser equiped it searches for new websites and the information needed to continue our intens investigation. This works by turning a website into a trie, traversing the trie and finding all needed tags. So the analyser uses optimized static text analysis-technics.