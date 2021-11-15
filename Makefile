DOCKER_REPO=nicholasjackson/smi-controller-consul
DOCKER_VERSION=dev.01

build_docker:
	docker build -t ${DOCKER_REPO}:${DOCKER_VERSION} .

push_docker:
	docker push ${DOCKER_REPO}:${DOCKER_VERSION}

run: fetch_certs
	go run main.go

fetch_certs:
	mkdir -p /tmp/k8s-webhook-server/serving-certs/
	
	kubectl get secret smi-controller-webhook-certificate -n shipyard -o json | \
		jq -r '.data."tls.crt"' | \
		base64 -d > /tmp/k8s-webhook-server/serving-certs/tls.crt
	
	kubectl get secret smi-controller-webhook-certificate -n shipyard -o json | \
		jq -r '.data."tls.key"' | \
		base64 -d > /tmp/k8s-webhook-server/serving-certs/tls.key

functional_tests: fetch_certs
# currently only splits have been fully implemented
	cd functional_tests	&& go run . --godog.tags=@split