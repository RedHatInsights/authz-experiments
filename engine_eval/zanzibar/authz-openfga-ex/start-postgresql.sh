
set -e

docker network create openfga

docker run -d --name postgres -p 5432:5432 --network=openfga -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password postgres:14
