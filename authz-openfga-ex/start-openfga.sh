docker run --rm --network=openfga openfga/openfga migrate \
    --datastore-engine postgres \
    --datastore-uri 'postgres://postgres:password@postgres:5432/postgres?sslmode=disable'
    
docker run --name openfga --network=openfga -p 3000:3000 -p 8080:8080 -p 8081:8081 openfga/openfga run \
    --datastore-engine postgres \
    --datastore-uri 'postgres://postgres:password@postgres:5432/postgres?sslmode=disable'