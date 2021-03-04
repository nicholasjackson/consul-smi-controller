DOCKER_REPO=nicholasjackson/smi-controller-consul
DOCKER_VERSION=dev

build_docker:
	docker build -t ${DOCKER_REPO}:${DOCKER_VERSION} .

push_docker:
	docker push ${DOCKER_REPO}:${DOCKER_VERSION}

run: fetch_certs
	go run main.go

fetch_certs:
	mkdir -p /tmp/k8s-webhook-server/serving-certs/
	
	kubectl get secret controller-webhook-certificate -n smi -o json | \
		jq -r '.data."tls.crt"' | \
		base64 -d > /tmp/k8s-webhook-server/serving-certs/tls.crt
	
	kubectl get secret controller-webhook-certificate -n smi -o json | \
		jq -r '.data."tls.key"' | \
		base64 -d > /tmp/k8s-webhook-server/serving-certs/tls.key

functional_tests: fetch_certs
	cd functional_tests	&& go run .