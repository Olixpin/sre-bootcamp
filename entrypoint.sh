#!/bin/sh

# Wait for the database to be ready
until nc -z -v -w30 ${DB_HOST} ${DB_PORT}
do
  echo "Waiting for the database to be available..."
  sleep 1
done

# Run the application
./student-api-binary
