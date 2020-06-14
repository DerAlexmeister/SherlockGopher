#!/bin/bash
declare -a arr=("frontend" "webserver" "sherlockcrawler" "analyser")
declare -a web=("-it -p 8080:8080" "-it -p 8081:8081" "" "")
image_name="dockerfile"
countweb=0

echo "[+] Start Script"

echo "[+] Installing curl"
apt install curl


if [[ "$(sudo docker ps -a | grep neo4j3.5)" != "" ]]; then
    echo "[+] Remove neo4j Docker"
    docker stop neo4j3.5
    docker rm neo4j3.5
fi

echo "[+] Start neo4j Docker"
docker run --name neo4j3.5 -p7474:7474 -p7687:7687 -d -v $HOME/neo4j/data:/data -v $HOME/neo4j/logs:/logs -v $HOME/neo4j/import:/var/lib/neo4j/import -v $HOME/neo4j/plugins:/plugins --env NEO4J_AUTH=neo4j/test neo4j:3.5

case "$1" in
    ("prune")
        echo "[+] Prune docker"
        docker system prune
        exit 0
        ;;
    ("frontend")
        echo "[+] only start frontend"
        for ((i=1; i<=3; i++)); do unset "arr[$i]"; done
        for ((i=1; i<=3; i++)); do unset "web[$i]"; done
        ;;
esac

for i in "${arr[@]}"
    do
        container_name=$i
        conid=$(sudo docker ps -a | awk '{ print $1,$2 }' | grep $container_name | awk '{print $1 }')

        if [[ "$(sudo docker ps -a | grep $container_name)" != "" ]]; then
            if [[ $(sudo docker ps -a --filter "status=running" | grep $container_name) != "" ]]; then
            echo "[+] Stop $container_name"
            docker stop $container_name
            fi
            echo "[+] Remove Container: $container_name"
            docker rm $conid  
        fi
        if [[ "$(sudo docker images | grep $container_name)" != "" ]]; then
            echo "[+] Remove Image: $container_name"
            docker rmi $container_name
        fi
        echo "[+] build and run: $container_name"          
        docker build -t $container_name -f ./$container_name/dockerfile .
        gnome-terminal -- sudo docker run ${web[$countweb]} $container_name
        ((countweb++))
    done
exit 0