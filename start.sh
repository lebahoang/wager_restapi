#!/usr/bin/env bash
set -e
set -x

# start docker containers
docker-compose down && docker-compose build && docker-compose up -d
# init db
until cat db_schema.sql |  docker exec -i wager_restapi_db_1 psql -U postgres
do
  echo "Waiting for postgres"
  sleep 2;
done
