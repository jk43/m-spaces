DOCKER_IMG=hub.hotdev.com:5555/molyspace-broker-service:latest
DOCKER_CONTAINER_NAME=broker-service 
build:
	@echo "Building Go"
	env GOOS=linux CGO_ENABLED=0 go build -o app ./cmd/api
	@echo "Done"

build-docker: build
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
