# Student API

A simple REST API for managing student records.

## Table of Contents

- [Overview](#overview)
- [Setup Instructions](#setup-instructions)
- [Environment Variables](#environment-variables)
- [Running the Application](#running-the-application)
- [Database Migrations](#database-migrations)
- [Endpoints](#endpoints)
- [Postman Collection](#postman-collection)
- [Makefile Commands](#makefile-commands)
- [Unit Tests](#unit-tests)

## Overview

This project provides a REST API for managing student records, including functionalities to add, retrieve, update, and delete student records. It also includes a health check endpoint.

## Setup Instructions

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/student-api.git
   cd student-api

2. Install dependencies:
   go mod download

3. Set up environment variables:
Create a .env file in the root directory and add the following variables:

DB_URL=your-database-url
BIND_ADDRESS=:8080

