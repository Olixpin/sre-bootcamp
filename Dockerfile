# # Use the official Golang image as the base image
# FROM golang:1.22

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy go.mod and go.sum files
# COPY go.mod go.sum ./

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the rest of the source code into the container
# COPY . .

# # Install make utility
# RUN apt-get update && apt-get install -y make

# # Copy wait-for-it.sh script and make it executable
# COPY wait-for-it.sh /wait-for-it.sh
# RUN chmod +x /wait-for-it.sh

# # Run the build command using Makefile
# RUN make build

# # Expose the port the application runs on
# EXPOSE 8080

# # Command to run the executable
# CMD /wait-for-it.sh db:5432 -- ./student-api-binary


# # Stage 1: Build the application
# FROM golang:1.22-alpine as builder

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Install necessary build tools
# RUN apk add --no-cache git make gcc libc-dev

# # Copy go.mod and go.sum files
# COPY go.mod go.sum ./

# # Download all dependencies
# RUN go mod download

# # Copy the rest of the source code into the container
# COPY . .

# # Copy wait-for-it.sh script and make it executable
# COPY wait-for-it.sh /wait-for-it.sh
# RUN chmod +x /wait-for-it.sh

# # Ensure static build
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o student-api-binary .

# # Stage 2: Run the application
# FROM alpine:latest

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Install necessary utilities
# RUN apk add --no-cache bash netcat-openbsd

# # Copy the built binary and the wait-for-it.sh script from the builder stage
# COPY --from=builder /app/student-api-binary /app/student-api-binary
# COPY --from=builder /wait-for-it.sh /wait-for-it.sh

# # Ensure wait-for-it.sh is executable
# RUN chmod +x /wait-for-it.sh

# # Ensure the binary is executable
# RUN chmod +x /app/student-api-binary

# # Expose the port the application runs on
# EXPOSE 8080

# # Command to run the executable
# CMD ["/wait-for-it.sh", "db:5432", "--", "./student-api-binary"]


# Stage 1: Build the application
FROM golang:1.22-alpine as builder

WORKDIR /app

RUN apk add --no-cache git make gcc libc-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Ensure static build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o student-api-binary .

# Stage 2: Run the application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary utilities
RUN apk add --no-cache bash netcat-openbsd curl

# Copy the built binary and the wait-for-it.sh script from the builder stage
COPY --from=builder /app/student-api-binary /app/student-api-binary
COPY --from=builder /wait-for-it.sh /wait-for-it.sh

# Ensure wait-for-it.sh is executable
RUN chmod +x /wait-for-it.sh

# Ensure the binary is executable
RUN chmod +x /app/student-api-binary

# Install the migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Copy migration files
COPY ./migrations /migrations

# Expose the port the application runs on
EXPOSE 8080

# Command to run the migrations and then start the application
CMD /wait-for-it.sh $DB_HOST $DB_PORT -- /usr/local/bin/migrate -path /migrations -database $DB_URL up && ./student-api-binary
