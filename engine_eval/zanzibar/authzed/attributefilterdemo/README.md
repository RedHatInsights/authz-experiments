# Attribute Filter Demo in SpiceDB (PIP-only version)
## Goals
* Store and retrieve permissions and attribute filters in SpiceDB
* Minimize calls to SpiceDB
## Non-goals
* Production security
* Performance/scalability

## How to run:

* Run a local instance of SpiceDB: ```$ docker run --name spicedb     -p 50051:50051     --rm     authzed/spicedb serve     --grpc-preshared-key "somerandomkeyhere"```
* Install Zed (https://github.com/authzed/zed)
* Configure Zed: ```$ zed context set local localhost:50051 "somerandomkeyhere" --insecure```
* Import the model: ```$ zed import https://play.authzed.com/s/xzsTAbFqjoJn/relationships```
* Compile the code: ```$ go build src/main.go```
* Run the program to analyze a user: ```$ ./main {user}``` where {user} is a username from the model (currently alec or eddy)