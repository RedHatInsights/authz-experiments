# Authzed SpiceDB integration test environment
This demonstrates how to perform integration tests against SpiceDB using docker and testcontainers.


1) Spin up a SpiceDB container running the serve-testing command.
2) For each independent test, create a SpiceDB client with a random key.
3) Run tests. Tests with different keys are safe to run in parallel.

## TODO
- container spins up, test seems to be executed, but afterwards endless loop.
- add scheme from YAML, see how to automate e.g. running assertions or sth

## Based on:
https://github.com/authzed/examples/tree/main/integration-testing