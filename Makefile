IMAGE_NAME ?= mhmoudgmal/localstack-single-endpoint
IMAGE_TAG ?= $(TAG)

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-push:
	docker login -u $$DOCKER_USERNAME -p $$DOCKER_PASSWORD;
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
