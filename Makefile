# .PHONY: build run migrate test clean docker-build docker-run

# SHELL := /bin/bash

# build:
#     CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o student-api-binary .

# run: build
# 	@./run.sh

# migrate:
# 	@migrate -database ${DB_URL} -path ./migrations up

# test:
# 	@go test ./...

# clean:
# 	@rm -f student-api-binary

# docker-build:
# 	@docker build -t sre-bootcamp-api:$(VERSION) .

# docker-run:
# 	@docker run --rm -e DB_URL=postgres://studentapiuser:password@db/studentapidb?sslmode=disable -e BIND_ADDRESS=:8080 -p 8080:8080 sre-bootcamp-api:$(VERSION)

# .PHONY: build run migrate test clean docker-build docker-run

# SHELL := /bin/bash

# build:
#     CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o student-api-binary .

# run: build
# 	./run.sh

# migrate:
# 	migrate -database ${DB_URL} -path ./migrations up

# test:
# 	go test ./...

# clean:
# 	rm -f student-api-binary

# docker-build:
# 	docker build -t sre-bootcamp-api:$(VERSION) .

# docker-run:
# 	docker run --rm --name sre-bootcamp-api-container --network sre-bootcamp-network -e DB_URL=postgres://studentapiuser:password@postgres-db/studentapidb?sslmode=disable -e BIND_ADDRESS=:8080 -e DB_HOST=postgres-db -e DB_PORT=5432 -p 8080:8080 sre-bootcamp-api:$(VERSION)

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
	docker run --name postgres-db --network sre-bootcamp-network -e POSTGRES_USER=studentapiuser -e POSTGRES_PASSWORD=password -e POSTGRES_DB=studentapidb -p 5432:5432 -d postgres:latest || true

docker-run: setup-network start-db
	docker run --rm --name sre-bootcamp-api-container --network sre-bootcamp-network -e DB_URL=postgres://studentapiuser:password@postgres-db/studentapidb?sslmode=disable -e BIND_ADDRESS=:8080 -e DB_HOST=postgres-db -e DB_PORT=5432 -p 8080:8080 sre-bootcamp-api:$(VERSION)
