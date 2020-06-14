# Neo4J-Testdatagenerator

## Explanation in one sentence

With this Pythonskript you can produce partly random data and put it in the database in order to get a feeling how the data of sherlock gopher might look like. 

## Installation

0) Requirements:
    - Docker
    - Python3
    - pip3 
    - Python3 in the PATH

1) Run the following command in order to create and run the Neo4j-Dockercontainer.

Incase you wanna change the username and password, just change the parameter NEO4J_AUTH= (yourusername)/(password). Also make sure you change the username/password in the skript.

```bash

sudo docker run --name neo4j3.5 -p7474:7474 -p7687:7687 -d -v $HOME/neo4j/data:/data -v $HOME/neo4j/logs:/logs -v $HOME/neo4j/import:/var/lib/neo4j/import -v $HOME/neo4j/plugins:/plugins --env NEO4J_AUTH=neo4j/test neo4j:3.5

```

2) Access the Neo4J Webinterface at:   

    ``` 127.0.0.1:7474 ```

3) Login with 
    
    ``` USER: neo4j, PASS: test ```

4) Install all Python dependencies

    ``` pip3 install -r ./requirements.txt ```

5) Run the Skript

    ``` ./testdata ```


6) Your output should look like this

```
[+] Starting to fill in the Database
[+] Creating HTML-Rels
[+] Creating CSS-Rels
[+] Creating Image-Rels
[+] Creating Javascript-Rels
[+] Finished - Database has now nodes and relationships.
```