# Authzed SpiceDB integration test environment
This demonstrates how to perform integration tests against SpiceDB using docker and testcontainers.


1) Spin up a SpiceDB container running the serve-testing command.
2) For each independent test, create a SpiceDB client with a random key.
3) Run tests. Tests with different keys are safe to run in parallel.

## TODO
- add scheme from YAML, see how to automate e.g. running assertions etc
- add GHA Workflow
## Based on:
https://github.com/authzed/examples/tree/main/integration-testing


Test 1.1: build should run as texxt is altered.