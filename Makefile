.ONESHELL:
.SHELL := /bin/bash
.PHONY: init build deploy get-url get-state list-containers delete-container destroy cover start test test-integration

init:
	@aws lightsail create-container-service --service-name canvas --power micro --scale 1 --tags key=label,value=app \
		--query 'containerService.{state: state, url: url, region: location.regionName, power: powerId, scale: scale}' --output table

build:
	@docker build -q -t canvas .

deploy: build
	@aws lightsail push-container-image --service-name canvas --image canvas --label app
	@IMG_NAME=$$(aws lightsail get-container-images --service-name canvas --query 'reverse(sort_by(containerImages, &createdAt))[0].image' --output text)
	@CONTAINER_PARAM=$$(jq --arg IMG_NAME "$${IMG_NAME}" -n '{"app":{"image":$$IMG_NAME,"environment":{"HOST":"","PORT":"8080","LOG_ENV":"production"},"ports":{"8080":"HTTP"}}}')
	@ENDPOINT_PARAM=$$(jq -n '{"containerName":"app","containerPort":8080,"healthCheck":{"path":"/health"}}')
	@aws lightsail create-container-service-deployment --service-name canvas \
		--containers "$${CONTAINER_PARAM}" \
		--public-endpoint "$${ENDPOINT_PARAM}" \
		--query 'containerService.{state: state, url: url, name: containerServiceName}' --output table

get-url:
	@aws lightsail get-container-services --service-name canvas --query 'containerServices[*].{region: location.regionName, power: powerId, scale: scale, name: containerServiceName, state: state, url: url, currentDeployment: {version: currentDeployment.version, state: currentDeployment.state, image: currentDeployment.containers.app.image}, nextDeployment: {version: nextDeployment.version, state: nextDeployment.state, image: nextDeployment.containers.app.image}}' --output table

get-state:
	@aws lightsail get-container-service-deployments --service-name canvas --query 'reverse(sort_by(deployments, &createdAt))[*].{version: version, state: state, container: containers.app.image}' --output table

list-containers:
	@aws lightsail get-container-images --service-name canvas --query containerImages[*].image --output table

delete-container:
ifndef version
	$(error version is undefined, set it as version=3)
endif
	@aws lightsail delete-container-image --service-name canvas --image ":canvas.app.$(version)"

destroy:
	@aws lightsail delete-container-service --service-name canvas

cover:
	@go tool cover -html=cover.out

start:
	@go run cmd/server/*.go

test:
	@go test -coverprofile=cover.out -short ./...

test-integration:
	@go test -coverprofile=cover.out -p 1 ./...
