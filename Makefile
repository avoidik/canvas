.ONESHELL:
.SHELL := /bin/bash
.PHONY: build deploy get-url cover start test test-integration

build:
	@docker build -t canvas .

deploy:
	@aws lightsail push-container-image --service-name canvas --label app --image canvas
	@IMG_NAME=$$(aws lightsail get-container-images --service-name canvas --query containerImages[0].image --output text)
	@CONTAINER_PARAM=$$(jq --arg IMG_NAME "$${IMG_NAME}" -n '{"app":{"image":$$IMG_NAME,"environment":{"HOST":"","PORT":"8080","LOG_ENV":"production"},"ports":{"8080":"HTTP"}}}')
	@ENDPOINT_PARAM=$$(jq -n '{"containerName":"app","containerPort":8080,"healthCheck":{"path":"/health"}}')
	@aws lightsail create-container-service-deployment --service-name canvas \
		--containers "$${CONTAINER_PARAM}" \
		--public-endpoint "$${ENDPOINT_PARAM}"

get-url:
	@aws lightsail get-container-services --service-name canvas --query containerServices[*].url --output text

cover:
	@go tool cover -html=cover.out

start:
	@go run cmd/server/*.go

test:
	@go test -coverprofile=cover.out -short ./...

test-integration:
	@go test -coverprofile=cover.out -p 1 ./...
