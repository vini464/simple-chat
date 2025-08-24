create_network:
	docker network create internet

build_images:
	docker-compose build

all: 
	docker-compose up --build 

run_server: 
	docker run --rm --network internet --name server server:v1

run_client: 
	docker run -it --rm --network internet --name client client:v1
