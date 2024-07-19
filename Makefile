.PHONY: run migrate test clean docker-build docker-run setup-network start-db logs

SHELL := /bin/bash

run:
	docker-compose up -d

migrate:
	docker-compose run migrate

test:
	go test ./...

clean:
	rm -f student-api-binary

docker-build:
	docker-compose build

docker-run: docker-build
	docker-compose up -d

setup-network:
	docker network create sre-bootcamp-network || true

start-db:
	docker-compose up -d db

logs:
	docker-compose logs -f

.PHONY: up down build migrate logs
