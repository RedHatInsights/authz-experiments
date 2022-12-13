set -e

docker stop openfga
docker stop postgres

docker rm openfga
docker rm postgres

docker network rm openfga