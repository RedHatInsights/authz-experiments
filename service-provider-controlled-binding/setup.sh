kubectl kcp workspace create aspian --enter --ignore-existing
kubectl apply -f aspian/aspian-employees.yaml
kubectl apply -f aspian/observability-team.yaml

kubectl kcp workspace create telemetry --enter --ignore-existing
kubectl apply -f aspian/telemetry/kafka-export.yaml
kubectl apply -f aspian/telemetry/telemetry-consumers.yaml

kubectl kcp workspace ..
kubectl kcp workspace create observability --ignore-existing