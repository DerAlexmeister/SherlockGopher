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
        - /controller/v1
            - 
```

## /areyouthere

- [Basicaddress]/areyouthere
- Type: GET
- Will send you as json a Pong as message as response to the "ping".

Response: 
```
{
    "message": "Pong"
}
```

## /graph/v1

### /meta 

- [Basicaddress]/graph/v1/meta
- Type: GET
- Will return some meta information of the Neo4J db such as the amount of nodes or number of relationships.

Response: 
```
...
    [
    [
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
    ], ... ]
]
```
### /all

- [Basicaddress]/graph/v1/all
- Type: GET
- Will return json which contains all nodes and their relationships to other nodes.

Response: 
```
[
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