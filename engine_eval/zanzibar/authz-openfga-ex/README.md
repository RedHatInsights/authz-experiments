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

## Authmodel Data

You can use the existing Json file in the modeldata folder, or create a new one for your test model (for model, tuples, assertions)
and export the same as follows
```
export OPENFGA_AUTH_MODEL_JSON_FILE="path-to-model.json-file"
export OPENFGA_AUTH_MODEL_TUPLES_JSON_FILE="path=to-tuples.json-file"
export OPENFGA_AUTH_MODEL_ASSERTION_JSON_FILE="path-to-assertions.json-file"
```
## Create a store auth_usermgmt and create OPENFGA_AUTH_MODEL
```
make run 
or
go run ./authz-openfga/main.go
```
Note: 
You will get errors when trying to add the same tuples again  (key - conflict), you can try to remove the values from table (postgres docker) or choose to ignore it
Run test against local openfga server
```
To run a specific test (since multiple models are there - you need to run only your test file/package)
Go to the test folder ex: authz-openfga-ex/test
go test <filename-to-test>
```