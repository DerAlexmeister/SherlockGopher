# REST-API of the Webserver

The webserver supplys a REST-API to fetch data from various points. For Example the state of the crawler and analysers queues can be retreived be the monitoring part of the API aswell as the data of the graph from the graph part (see structure below).

```apistructure
    - [Basicaddress]
        - /areyouthere 
        - /graph/v1   
            - /meta 
            - /all
            - /performenceofsites
            - /detailsofnode 
            - /search
        - /monitor/v1
            - /meta
```

## /areyouthere

- [Basicaddress]/areyouthere
- Type: GET
- Will send you a Pong as json to respond to the "ping".

**Response:**
```Json
{
    "message": "Pong"
}
```

## /graph/v1

### /meta 

- [Basicaddress]/graph/v1/meta
- Type: GET
- Will return some meta information of the Neo4J db such as the amount of nodes or number of relationships.

**Response:**
```Json
    ...
        {
            "amountofimages": 18
        }
    ],
    [
        {
            "amountofsheets": 10
        }
    ],
    [
        {
            "amountofjs": 20
        }
    ], ...
]
```
### /all

- [Basicaddress]/graph/v1/all
- Type: GET
- Will return json which contains all nodes and their relationships to other nodes.

**Response:**
```Json
    ...
    {
        "Type(r)": "Links",
        "k.Address": "0",
        "k.Filetype": "HTML",
        "n.Address": "12",
        "n.Filetype": "HTML"
    },
    {
        "Type(r)": "Links",
        "k.Address": "0",
        "k.Filetype": "HTML",
        "n.Address": "61",
        "n.Filetype": "HTML"
    },
    {
        "Type(r)": "Links",
        "k.Address": "0",
        "k.Filetype": "HTML",
        "n.Address": "48",
        "n.Filetype": "HTML"
    }, ...
```

### /performenceofsites

- [Basicaddress]/graph/v1/performenceofsites
- Type: GET
- Will return a list of performence indicators like statuscode and rtt.

**Response:**
```Json
    ...
        {
            "n.Address": "0",
            "n.Responsetime": "333",
            "n.Statuscode": "293"
        },
        {
            "n.Address": "1",
            "n.Responsetime": "1521",
            "n.Statuscode": "446"
        }, ...
```

### /detailsofnode

- [Basicaddress]/graph/v1/detailsofnode
- Type: POST
- Will return more details about a given node.

**Request:**
```Json
    ... {
        "url": "www.github.com/imprint"
    } ...
```

**Response:**
```Json
    ... "www.github.com/imprint": {
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
    ...
```