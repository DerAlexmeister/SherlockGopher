# Manual to install or run SherlockGropher


Vorraussetzung GO 1.12

go get go.mod to install all dependencys

Most of the scripts and the installation guide are only working on debian based systems

## How to install all dependencies to start the services via main file

execute the seabolt.sh script:
- ./scripts/seabolt.sh
- this script installs all necessary dependencies for neo4j

start the neo4j docker container with the following command:
* sudo docker run --name neo4j3.5 -p7474:7474 -p7687:7687 -d -v $HOME/  neo4j/data:/data -v $HOME/neo4j/logs:/logs -v $HOME/neo4j/import:/var/lib/neo4j/import -v $HOME/neo4j/plugins:/plugins --env NEO4J_AUTH=neo4j/test neo4j:3.5

### How to start all services via main file
start the analyser service:
- go run analyser/main.go

start the frontend service:
- go run frontend/main.go

start the sherlockcrawler service:
- go run sherlockcrawler/main.go

start the webserver service:
- go run webserver/main.go

## How to start all services via dockerfile - Docker Start Script

Alternatively you can use the start.sh script. The script will stop all running docker containers related to this project, build them and execute them afterwards. Also it tries to install curl and starts the neo4j docker container.
The script can be executed on the command line via: ./start.sh

There are 2 possible arguments that the script accepts. These arguments will change the functionality of the script:
* ./start.sh prune 
  * with this argument the script will only prune all unused containers
* ./start.sh frontend 
  * with this argument the script will only start the frontend docker container