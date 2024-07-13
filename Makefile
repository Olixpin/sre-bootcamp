.PHONY: build run migrate test clean docker-build docker-run setup-network start-db

SHELL := /bin/bash

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o student-api-binary .

run: build
	./run.sh

migrate:
	migrate -database ${DB_URL} -path ./migrations up

test:
	go test ./...

clean:
	rm -f student-api-binary

docker-build:
	docker build -t sre-bootcamp-api:$(VERSION) .

setup-network:
	docker network create sre-bootcamp-network || true

start-db:
	set -o allexport; source .env; set +o allexport; \
	docker run --name postgres-db --network sre-bootcamp-network \
	  -e POSTGRES_USER=$${POSTGRES_USER} \
	  -e POSTGRES_PASSWORD=$${POSTGRES_PASSWORD} \
	  -e POSTGRES_DB=$${POSTGRES_DB} \
	  -p 5432:5432 -d postgres:latest || true

docker-run: setup-network start-db
	set -o allexport; source .env; set +o allexport; \
	docker run --rm --name sre-bootcamp-api-container --network sre-bootcamp-network \
	  -e DB_URL=postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@postgres-db/$${POSTGRES_DB}?sslmode=disable \
	  -e BIND_ADDRESS=:8080 \
	  -e DB_HOST=postgres-db \
	  -e DB_PORT=5432 \
	  -p 8080:8080 sre-bootcamp-api:$(VERSION)
