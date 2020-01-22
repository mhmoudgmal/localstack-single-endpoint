IMAGE_NAME ?= mhmoudgmal/localstack-single-endpoint
IMAGE_TAG ?= "${$(git describe --tags):1}"

docker-build:
	docker build -t $(IMAGE_NAME) .
	docker tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(IMAGE_TAG)

docker-push:
	docker login -u $$DOCKER_USERNAME -p $$DOCKER_PASSWORD
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
