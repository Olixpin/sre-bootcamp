# Student API

A simple REST API for managing student records.

## Table of Contents

- [Student API](#student-api)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Setup Instructions](#setup-instructions)
  - [To run the application, use the following command:](#to-run-the-application-use-the-following-command)
  - [To run database migrations and create the necessary tables, use the following command:](#to-run-database-migrations-and-create-the-necessary-tables-use-the-following-command)
  - [Endpoints](#endpoints)
  - [Example Requests](#example-requests)
  - [Postman Collection](#postman-collection)
  - [Importing Postman Collection](#importing-postman-collection)
  - [Makefile Commands](#makefile-commands)

## Overview

This project provides a REST API for managing student records, including functionalities to add, retrieve, update, and delete student records. It also includes a health check endpoint.

## Setup Instructions

1. **Clone the repository:**
   ```sh
   git clone https://github.com/olixpin/student-api.git
   cd student-api

2. **Install Dependencies**
   ```sh
   go mod download
3. **Set up environment variables:**
Create a .env file in the root directory and add the following variables:
   ```sh
   DB_URL=postgres://studentapiuser:password@localhost/studentapidb?sslmode=disable
   BIND_ADDRESS=:8080
4. **Run database migrations:**
   ```sh
   make migrate
5. Run the server:
   ```sh
   make run
 6. **Environment Variables**
    ```sh
    DB_URL: The URL for connecting to the database.
    BIND_ADDRESS: The address on which the server will listen (default is :8080).
    Running the Application

## To run the application, use the following command:
    Database Migrations

## To run database migrations and create the necessary tables, use the following command:
    make migrate
## Endpoints
    GET /api/v1/students: Get all students
    GET /api/v1/students/{id}: Get a student by ID
    POST /api/v1/students: Add a new student
    PUT /api/v1/students/{id}: Update a student
    DELETE /api/v1/students/{id}: Delete a student
    GET /healthcheck: Health check endpoint

## Example Requests
**Fetch All Students**

    curl -X GET http://localhost:8080/api/v1/students
    
**Fetch a student by ID**

    curl -X GET http://localhost:8080/api/v1/students/1

**Add a New Student**

    curl -X POST http://localhost:8080/api/v1/students -H "Content-Type: application/json" -d '{
      "first_name": "Alice",
      "last_name": "Johnson",
      "age": 23,
      "email": "alice.johnson@example.com",
      "enrollment_date": "2024-07-08T20:34:02.741269Z",
      "class": "C",
      "address": "789 Maple St, Anytown USA",
      "phone_number": "555-7890"
    }'

**Update an Existing Student**
    
    curl -X PUT http://localhost:8080/api/v1/students/1 -H "Content-Type: application/json" -d '{
      "first_name": "John",
      "last_name": "Doe",
      "age": 21,
      "email": "john.doe@example.com",
      "enrollment_date": "2024-07-08T20:34:02.741269Z",
      "class": "A",
      "address": "123 Main St, Anytown USA",
      "phone_number": "555-1234"
    }'

**Delete a Student Record**

    curl -X DELETE http://localhost:8080/api/v1/students/1

## Postman Collection

A Postman collection for the APIs is included in the repository. You can import it into Postman to easily test the endpoints. The collection file is named Student_API.postman_collection.json and is located in the root directory.

## Importing Postman Collection
Open Postman.
Click on the Import button in the top left corner.
Select the Choose Files button.
Browse to the location of the Student_API.postman_collection.json file and select it.
Click Open to import the file.

## Makefile Commands
**Build the application:**

    make build

**Run the application:**

    make run

**Run database migrations:**

    make migrate

**Run tests:**  

    make test

**Unit Tests**
To run the unit tests, use the following command:

    go test ./handlers