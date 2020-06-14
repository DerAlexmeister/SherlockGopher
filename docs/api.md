# REST-API of the Webserver

The webserver supplys a REST-API to fetch data from various points. For Example the state of the crawler and analysers queues can be retreived by the monitoring part of the API aswell as the data of the graph from the graph part (see structure below).

```apistructure
    - [Basicaddress]
        - /areyouthere 
        - /graph/v1   
            - /meta 
            - /all
            - /alloptimized
            - /performenceofsites
            - /detailsofnode 
            - /search
        - /monitor/v1
            - /meta
        - /controller/v1
            - /dropit
            - /status
            - /changestate
```

## /areyouthere

- [Basicaddress]/areyouthere
- Type: GET
- Will send you a "yes i am here" as json to respond to the "ping".

**Response:**
```Json
{
    "Message": "Yes i am here!"
}
```

## /graph/v1

### /meta 

- [Basicaddress]/graph/v1/meta
- Type: GET
- Will return some meta information of the Neo4J db such as the amount of nodes or number of relationships.

**Response:**
```Json
    [
        {
            "amountofimages": 18
        },
        {
            "amountofsheets": 10
        },
        {
            "amountofjs": 20
        } ...] 
```
### /all

- [Basicaddress]/graph/v1/all
- Type: GET
- Will return json which contains all nodes and their relationships to other nodes.

**Response:**
```Json
    ...
    {
        "Destination": "www.example.com/6",
        "DestinationType": "Javascript",
        "Relationship": "Requires",
        "Source": "www.example.com/23",
        "SourceType": "HTML"
    },
    {
        "Destination": "www.example.com/7",
        "DestinationType": "CSS",
        "Relationship": "Requires",
        "Source": "www.example.com/79",
        "SourceType": "HTML"
    }, ...
```

### /alloptimized

- [Basicaddress]/graph/v1/alloptimized
- Type: GET
- Will return json which contains all nodes and their relationships to other nodes but optimized for the react frontend.

**Response:**
```Json
    ...
    "links": [
        {
            "color": "#08206A",
            "label": "Links",
            "source": "www.example.com/43",
            "target": "www.example.com/0"
        },
        {
            "color": "#08206A",
            "label": "Links",
            "source": "www.example.com/41",
            "target": "www.example.com/0"
        }],
        "nodes": [
        {
            "color": "#7CA9EF",
            "id": "www.example.com/43"
        },
        {
            "color": "#7CA9EF",
            "id": "www.example.com/0"
        },
        {
            "color": "#7CA9EF",
            "id": "www.example.com/41"
        }] ...
```

### /performanceofsites

- [Basicaddress]/graph/v1/performenceofsites
- Type: GET
- Will return a list of performence indicators like statuscode and RTT.

**Response:**
```Json
    ...
        {
            "Address": "www.example.com/0",
            "ResponseTime": "522",
            "Status": "409"
        },
        {
            "Address": "www.example.com/1",
            "ResponseTime": "1034",
            "Status": "246"
        }, ...
```

### /detailsofnode

- [Basicaddress]/graph/v1/detailsofnode
- Type: POST
- Will return more details about a given node.

**Request:**
```Json
    ... {
        "url": "www.example.com/imprint"
    } ...
```

**Response:**
```Json
    ... "www.example.com/imprint": {
            "Accept-Ranges": "bytes",
            "Address": "1",
            "Age": "12",
            "Allow": "GET, HEAD",
            "Cache-Control": "max-age=3600",
            "Connection": "close",
            "Content-Encoding": "gzip",
            "Content-Language": "de",
            "Content-Location": "/imprint",
            "Content-MD5": "Q2hlY2sgSW50ZWdyaXR5IQ==",
            "Content-Type": "text/html; charset=utf-8",
        } ...
```

### /search
- [Basicaddress]/graph/v1/search
- Type: POST
- Will return a message containing information about the request like Fine or an error containing a message.

**Request:**
```Json
    ... {
        "url": "www.github.com/"
    } ...
```

**Response:**
```Json
    ...
    {
        "Status": "Fine"
    } 
    ...
```

## /monitor/v1

### /meta
- [Basicaddress]/monitor/v1/meta
- Type: GET
- Will return meta information about the services (Analyser and Crawler).

**Response:**
```Json
    "Crawler":{
        "Website":    45,
        "Undone":     9,
        "Processing": 25,
        "Finished":   10,
        "Failed":     1,
    },
    "Analyser":{
        "Website":       45,
        "Undone":        5,
        "Processing":    25,
        "CrawlerError":  0,
        "Saving":        5,
        "SendToCrawler": 5,
        "Finished":      5,
    }
```

## /controller/v1

### /changestate
- [Basicaddress]/controller/v1/changestate
- Type: POST
- Will change the state of various services.
- Operations: Clean, Stop, Pause, Resume
- Target: Crawler, Analyser, All

**Request:**
```Json
    ... {
        "operation": "clean",
        "target": "crawler"
    } ...
```

**Response:**
```Json
        {
		"Status": "Fine",
	}
```

### /status
- [Basicaddress]/controller/v1/status
- Type: GET
- Will return the status of the crawler and analyser.

**Response:**
```Json
   {
       "Analyser": "Running",
       "Crawler": "Paused"
   }
```

### /dropit
- [Basicaddress]/controller/v1/dropit
- Type: GET
- Will drop the database.

**Response:**
```Json
   {
       "Message": "Dropped the table."
   }
```