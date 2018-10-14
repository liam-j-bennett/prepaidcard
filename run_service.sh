#!/bin/bash

docker build -f Dockerfile -t prepaidcard:latest .

NETWORK_ID=$(docker network create local-network)
DATABASE_ID=$(docker run -d --publish 5432:5432 --name ppc_database --network local-network postgres:latest)
sleep 3
API_ID=$(docker run -d --publish 8080:8080 --network local-network -e DB_HOST=ppc_database prepaidcard:latest)

read -p "Local env running... Press any key to kill."

docker kill ${API_ID}
docker kill ${DATABASE_ID}

docker rm ${DATABASE_ID}
docker rm ${API_ID}