#!/bin/bash
declare -a arr=("frontend" "webserver" "crawler" "analyser")
declare -a web=("-it -p 8080:8080" "-it -p 8081:8081" "" "")
image_name = dockerfile
countweb = 0

echo "Start Docker Script"
sudo docker run --name neo4j3.5 -p7474:7474 -p7687:7687 -d -v $HOME/neo4j/data:/data -v $HOME/neo4j/logs:/logs -v $HOME/neo4j/import:/var/lib/neo4j/import -v $HOME/neo4j/plugins:/plugins --env NEO4J_AUTH=neo4j/test neo4j:3.5

case $1 in
    frontend)
        echo "only start frontend"
        for($z = 1 ; $z < ${#arr[@]} ; $z++)
        do
            unset($arr[$z]);
            unset($web[$z]);
        done
        ;;
    prune)
        echo "Prune docker"
        docker system prune
        exit 0
        ;;
    rebuild)
        for i in "${arr[@]}"
        do
            container_name = i
            if [[ "$(docker ps -a | grep $container_name)" = ""]]; then
                echo "Stop $container_name"
                docker stop $container_name
            fi
            cd $container_name
            echo "rebuild docker container and image for $container_name"
            sudo docker build -f ./$container_name/dockerfile .
            docker run -d ${web[$countweb]} --name $container_name $image_name
            echo "start docker container and image for $container_name"
            docker start $container_name
            cd ..
            ((countweb++))
        done
        exit 0
        ;;
    stop)
    for i in "${arr[@]}"
    do
        container_name = i
        if [[ "$(docker ps -a | grep $container_name)" = ""]]; then
        echo "Stop $container_name"
        docker stop $container_name
    done
    exit 0
    ;;
esac

for i in "${arr[@]}"
do
    container_name = i
    cd $container_name
    echo "Check if Container $container_name already exists"
    if [[ "$(docker ps -a | grep $container_name)" ]] \
    && [[ "$(docker images -q $image_name 2> /dev/null)" != "" ]]; then #unsicher
        echo "docker container  $container_name already exists"
        echo "docker image $image_name already exists"
        docker start $container_name

    elif [[ "$(docker ps -a | grep $container_name)" ]]; then
        echo "docker container  $container_name already exists,"
        echo "container image $image_name missing, building image..."
        docker run -d ${web[$countweb]} --name $container_name $image_name
        docker start $container_name

    else
        echo "docker container and image for $container_name are missing, building content..."
        sudo docker build -f ./$container_name/dockerfile .
        docker run -d ${web[$countweb]} --name $container_name $image_name
        docker start $container_name
    fi
    cd ..
    ((countweb++))
done