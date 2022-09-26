Access Control for Cross-Workspace API Bindings
======================

Setup
---------------
1. Start KCP. These examples will assume a local instance of KCP with root workspace access
    * Start from the service-provider-controlled-binding folder as your working directory
    * Run: `kcp start --token-auth-file=tokens`
        * The tokens file allows us to log in as different users described in that file
1. Switch to another terminal or tab
1. Run `export KUBECONFIG=.kcp/admin.kubeconfig`
    * This will configure kubectl to talk to the new KCP instance for this terminal
1. Navigate to the root workspace: `kubectl kcp workspace root`
1. Run: `./setup.sh`
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