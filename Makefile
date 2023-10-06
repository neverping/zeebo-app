.PHONY: run create-minikube install-ingress create-go-container create-python-container

create-minikube:
	minikube start

install-ingress:
	minikube addons enable ingress

create-go-container:
	cd services/go && \
	DOCKER_BUILDKIT=1 docker image build -t zeebo-go:latest --no-cache .

create-python-container:
	cd services/python && \
	DOCKER_BUILDKIT=1 docker image build -t zeebo-python:latest --no-cache .

deploy-helm-chart:
	cd helm/charts && \
	helm upgrade -i --cleanup-on-fail v1 zeebo/

run-docker-compose:
	docker-compose up -d

stop-docker-compose:
	docker-compose down

destroy-docker-compose:
	docker-compose down --rmi all --remove-orphans
