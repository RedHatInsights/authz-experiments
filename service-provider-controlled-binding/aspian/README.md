Access Control for Cross-Workspace API Bindings
======================

Setup
---------------
1. Start KCP. These examples will assume a local instance of KCP with root workspace access
    * Start from the service-provider-controlled-binding folder as your working directory
    * Run: `kcp start --token-auth-file=tokens`
        * The tokens file allows us to log in as different users described in that file
1. In another terminal, navigate to the root workspace: `kubectl kcp workspace root`
    * If you have a workspace named 'aspian' in your root, see the note on the next item
1. Run: `./setup.sh`
    * Note: This will create a workspace in your root called 'aspian' - if you don't want this (ex: you already have a workspace named aspian that you want to keep), you can update the name of the org workspace on line 1 of setup.sh and 12 of aspian-employees.yaml
1. Your current workspace should be root:aspian

That created a miniature Aspian organization such that Alex is an employee of Aspian and a member of team observability, and team telemetry exports a Kafka service for which they've created a telemetry-consumers role in their workspace to control who can bind it.

Testing
-------
If you now navigate to the observability workspace (`kubectl kcp workspace observability`), you can try to bind the exported Kafka as Alex by running `kubectl apply -f aspian/observability/kafka-binding.yaml --token alex` which should fail with the error: `unable to create APIImport: missing verb='bind' permission on apiexports`

This is because team telemetry hasn't authorized the binding yet1

To authorize the binding, navigate to the telemetry workspace `kubectl kcp workspace root:aspian:telemetry` and apply one of the binding yamls. For example, the following will allow any member of the observability team to bind this API: `kubectl apply -f aspian/telemetry/telemetry-consumers-pergroup.yaml`

To verify this, navigate back to the observability workspace and bind the API as alex:
```
kubectl kcp workspace root:aspian:observability
kubectl apply -f aspian/observability/kafka-binding.yaml --token alex
```
..which should succeed.

Further Experimentation
---------------------------
The goal of this exercise was to set up a known-good lab based on KCP's documentation which could then serve as a foundation for testing variations on it, including webhook and custom authorizers.