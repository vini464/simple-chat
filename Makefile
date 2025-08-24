create_network:
	docker network create internet

build_images:
	docker-compose build

all: build_images
	docker-compose up 

run_server: build_images
	docker run --rm --network internet --name server server:v1

run_client: build_images
	docker run -it --rm --network internet --name client client:v1
