# Open FGA Experiment 

# Run Openfga server with postgres
1) Using docker compose: `docker compose up`
or using docker scripts
2) Postgresql database `./start-postgresql.sh`
2) Openfga playground `./start-openfga.sh`

# Run
```
export OPENFGA_API_SCHEME="http" 
export OPENFGA_API_HOST="0.0.0.0:8080" //openfga running on localhost:8001 (playground localhost:3000)
```

Run on non-go path enable the go modules
`export GO111MODULE=on`

## Sample Auth model
```
export OPENFGA_AUTH_MODEL="{\"type_definitions\":[{\"type\":\"document\",\"relations\":{\"reader\":{\"this\":{}},\"writer\":{\"this\":{}},\"owner\":{\"this\":{}}}}]}"
```
## Create a store auth_usermgmt and create OPENFGA_AUTH_MODEL
```
go run main.go
```

Run test against local openfga server
```
go test authz-openfga-ex/test -test.v
```