Access Control for Cross-Workspace API Bindings
======================

# Setup
1. Start KCP. These examples will assume a local instance of KCP with root workspace access
    * Start from the service-provider-controlled-binding folder as your working directory
    * Run: `kcp start --token-auth-file=tokens`
        * The tokens file allows us to log in as different users described in that file
1. Switch to another terminal or tab
1. Run `export KUBECONFIG=.kcp/admin.kubeconfig`
    * This will configure kubectl to talk to the new KCP instance for this terminal
1. Run: `./setup.sh`
1. Your current workspace should be root:aspian

That created a miniature Aspian organization such that Alex and Aidan are employees of Aspian and members of the Telemetry and Observability teams respectively. Team Telemetry exports a Kafka service for which they've created a telemetry-consumers role in their workspace to control who can bind it.

There's also a redhat organization to sell managed services to Aspian, with a managed-kafka team within.

# Testing
## APIBinding Relationships Within a Tenant

If you now navigate to the observability workspace (`kubectl kcp workspace observability`), you can try to bind the exported Kafka as Aidan by running `kubectl apply -f aspian/observability/kafka-binding.yaml --token aidan` which should fail with the error: `unable to create APIImport: missing verb='bind' permission on apiexports`

This is because team telemetry hasn't authorized the binding yet!

To authorize the binding, navigate to the telemetry workspace `kubectl kcp workspace root:aspian:telemetry` and apply one of the binding yamls. For example, the following will allow any member of the observability team to bind this API: `kubectl apply -f aspian/telemetry/telemetry-consumers-pergroup.yaml`

To verify this, navigate back to the observability workspace and bind the API as Aidan:
```
kubectl kcp workspace root:aspian:observability
kubectl apply -f aspian/observability/kafka-binding.yaml --token aidan
```
..which should succeed.

## APIBinding Relationships Across Tenants (RBAC)
!! Note: This experiment requires latest KCP unstable !!

If you now navigate to the telemetry workspace and try to bind the managed kafka instance as Alex, it should fail due to no authorization.

```
kubectl kcp workspace root:aspian:telemetry
kubectl apply -f aspian/telemetry/kafka-binding.yaml --token alex
```

As admin, navigate to the managed kafka org and grant members of team telemetry permission to bind the API, then try it again:

```
kubectl kcp workspace root:redhat:managed-kafka
kubectl apply -f redhat/managed-kafka/kafka-consumers.yaml

kubectl kcp workspace root:aspian:telemetry
kubectl apply -f aspian/telemetry/kafka-binding.yaml --token alex
```
## APIBinding Relationships Across Tenants (Custom)
!! Note: This experiment requires a custom build of KCP. See: https://github.com/wscalf/kcp-experiments/tree/resource-configured-custom-authorizer !!

!! Note: This experiment assumes the previous one was *not* performed. You can stop KCP, run `rm -rf .kcp`, start KCP back up, and re-run `./setup.sh` to hard-reset your experimental environment, if needed. !!

As before, if you navigate to the telemetry workspace and try to bind the managed kafka instance as Alex, it should fail due to no authorization.

```
kubectl kcp workspace root:aspian:telemetry
kubectl apply -f aspian/telemetry/kafka-binding.yaml --token alex
```

As admin, navigate to the managed kafka org and create an authzconfig custom resource that grants members of org 123 permission to bind APIs, then try again:

```
kubectl kcp workspace root:redhat:managed-kafka
kubectl apply -f redhat/managed-kafka/authz-config.yaml

kubectl kcp workspace root:aspian:telemetry
kubectl apply -f aspian/telemetry/kafka-binding.yaml --token alex
```

..And it should succeed. However, because the authzconfig object exists in the managed-kafka workspace, it only applies to that workspace and does not allow users to bind APIs hosted elsewhere. For example:

```
kubectl kcp workspace root:aspian:observability
kubectl apply -f aspian/observability/kafka-binding.yaml --token aidan
```

..will fail due to not authorization. And Aidan can be granted permission to bind to telemetry's Kafka API either by using RBAC as before (to show the systems work in parallel)
```
kubectl kcp workspace root:aspian:telemetry
kubectl apply -f aspian/telemetry/telemetry-consumers-pergroup.yaml

kubectl kcp workspace root:aspian:observability
kubectl apply -f aspian/observability/kafka-binding.yaml --token aidan
```

..OR by creating an authzconfig in the telemetry workspace that specifies the 123 orgId (for brevity, we can use the same YAML as before in a different workspace)
```
kubectl kcp workspace root:aspian:telemetry
kubectl apply -f redhat/managed-kafka/authz-config.yaml

kubectl kcp workspace root:aspian:observability
kubectl apply -f aspian/observability/kafka-binding.yaml --token aidan
```
# Further Experimentation
The goal of this exercise was to set up a known-good lab based on KCP's documentation which could then serve as a foundation for testing variations on it, including webhook and custom authorizers.