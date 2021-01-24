DOCKER_REPO=nicholasjackson/smi-controller-consul
DOCKER_VERSION=dev

build_docker:
	docker build -t ${DOCKER_REPO}:${DOCKER_VERSION} .

push_docker:
	docker push ${DOCKER_REPO}:${DOCKER_VERSION}

run:
	ENABLE_WEBHOOKS=false go run main.go
