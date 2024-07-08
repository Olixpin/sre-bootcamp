#!/bin/bash

# Source environment variables from .env file
export $(grep -v '^#' .env | xargs)

# Print the environment variables to verify
echo "DB_URL=${DB_URL}"
echo "BIND_ADDRESS=${BIND_ADDRESS}"

# Run the application
./student-api-binary
