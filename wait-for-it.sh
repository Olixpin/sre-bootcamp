#!/bin/sh

# Function to check for database availability
wait_for_db() {
  DB_HOST=$1
  DB_PORT=$2
  while ! nc -z $DB_HOST $DB_PORT; do 
    echo "Waiting for the database to be available..."
    sleep 1
  done
}

# Check for the availability of the database
wait_for_db $1 $2

# Shift the positional parameters to pass remaining arguments to exec
shift 2

# Execute the main application
exec "$@"
