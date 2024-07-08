.PHONY: build run migrate test clean

SHELL := /bin/bash

build:
	go build -o student-api-binary
	chmod +x student-api-binary

run: build
	./run.sh

migrate:
	migrate -database ${DB_URL} -path ./migrations up

test:
	go test ./...

clean:
	rm -f student-api-binary
