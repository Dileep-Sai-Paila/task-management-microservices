#!/bin/bash
# While docker containerization, it couldn't recognize the task_db for the task-service, hence using this bash script
set -e

# Connect to the default database (user_db) and execute a command to create our second database.
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE task_db;
EOSQL
