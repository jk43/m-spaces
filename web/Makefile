DOCKER_IMG=hub.hotdev.com:5555/molyspace-web:latest
DOCKER_CONTAINER_NAME=web 
build:
	@echo "Building Vue"
	npm install
	BROKER_URL="https://api.molyspace.com" quasar build
	@echo "Done"

build-stg:
	@echo "Building Vue"
	npm install
	BROKER_URL="https://api.stg.molyspace.com" quasar build
	@echo "Done"

docker: build
	@echo "Building docker image"
	docker build -f Dockerfile -t ${DOCKER_IMG} .
	docker push ${DOCKER_IMG}

docker-stg: build-stg
	@echo "Building docker image"
	docker build -f Dockerfile -t ${DOCKER_IMG} .
	docker push ${DOCKER_IMG}

docker-up-local: build-prod
	@echo "Removing container"
	docker stop ${DOCKER_CONTAINER_NAME} 
	docker rm ${DOCKER_CONTAINER_NAME} 
	@echo "Running container"
	docker run -it -p 8099:8040 --name ${DOCKER_CONTAINER_NAME} ${DOCKER_IMG} 

docker-down-local:
	docker stop ${DOCKER_CONTAINER_NAME} 
	docker rm ${DOCKER_CONTAINER_NAME}
