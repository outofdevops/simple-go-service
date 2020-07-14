#/usr/bin/env bash

set -eu

LOCAL_PORT=5432
DB_CONTAINER_NAME=testdb
POSTGRES_PASSWORD=testdb
POSTGRES_USER=testdb
function cleanup {
  echo "Removing db container..."
  docker rm -fv $DB_CONTAINER_NAME || true
}

trap cleanup EXIT

echo "Starting postgres..."

docker run -d --rm -p ${LOCAL_PORT}:5432 \
  -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
  -e POSTGRES_USER=$POSTGRES_USER \
  --name=$DB_CONTAINER_NAME \
  postgres

echo "Postgres launched!"


echo "Postgres is unavailable, waiting for port\c"
until docker exec -it -e PGPASSWORD=testdb $DB_CONTAINER_NAME psql -h "localhost" -U "testdb" -c '\q' > /dev/null; do
  echo ".\c"
  sleep 0.5
done
  
>&2 echo "\n\nPostgres ready to accept connections!!!"

go test -v

