build:
	sudo docker-compose build
down:
	sudo docker-compose down
infup:
	sudo docker-compose up zookeeper kafka mongodb_container postgres neo4j
serup:
	sudo docker-compose up --build webserver frontend analyser crawler screenshot analysepictures
inf:
	sudo docker-compose run zookeeper kafka mongodb_container postgres neo4j
ser:
	sudo docker-compose run webserver frontend analyser crawler screenshot analysepictures