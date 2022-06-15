#!/bin/sh

sleep 5s # wait for DB to get initialized

psql -c "CREATE DATABASE auri_development;"
psql -c "CREATE DATABASE auri_test;"
psql -c "CREATE DATABASE auri_production;"

psql -c "GRANT ALL PRIVILEGES ON DATABASE auri_development to postgres;"
psql -c "GRANT ALL PRIVILEGES ON DATABASE auri_test to postgres;"
psql -c "GRANT ALL PRIVILEGES ON DATABASE auri_production to postgres;"

sleep 10d
