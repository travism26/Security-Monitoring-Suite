# Notes

- The `system-monitoring-gateway` service is exposed on port `30001` on the node port.
- The `kafka-ui` service is exposed on port `30007` on the node port.
  - http://localhost:30007/ui/clusters/systems-kafka-cluster/all-topics

## Commands to get details

```bash
# List Pods in the Namespace:
kubectl get pods

# List Services in the Namespace:
kubectl get svc

# Describe a Pod:
kubectl describe pod <pod-name>

# Get Logs for a Pod:
kubectl logs <pod-name>

# Execute a Command in a Pod:
kubectl exec -it <pod-name> -- /bin/bash

# Get Environment Variables for a Pod:
kubectl exec -it <pod-name> -- env

# Get ConfigMap:
kubectl get configmap

# Describe ConfigMap:
kubectl describe configmap system-monitoring-config

# Get Container Status for a Pod:
kubectl get pods <pod-name> -o jsonpath='{.status.containerStatuses[*].state}'

# Get Liveness Probe for a Deployment:
kubectl describe deployment system-monitoring-gateway | grep -A 10 'Liveness'

```

# Commands to hit the endpoints

```bash
# Health check (Verbose)
curl -v http://localhost:30001/health

# Metrics (Verbose)
curl -v http://localhost:30001/metrics

# Ingest metrics
curl -X POST http://localhost:30001/api/v1/system-metrics/ingest -H "Content-Type: application/json" -d '{"data":{"metrics":{"cpu":{"usage":{"total":100,"idle":50}},"memory":{"total":1000,"available":500}},"timestamp":"2024-01-01T00:00:00Z"}}'
```

# Run migrations

## Apply the ConfigMap first

```bash
kubectl apply -f infra/k8s/postgres-migrations-configmap.yaml
```

## Wait for PostgreSQL pod to be ready

```bash
kubectl wait --for=condition=ready pod -l app=postgres
```

## Run the migrations

```bash
kubectl apply -f infra/k8s/postgres-migrations-job.yaml
```

## Check migration job status

```bash
kubectl get jobs postgres-migrations
kubectl logs job/postgres-migrations
```

# Pgadmin login

## Main login

Email: admin@admin.com
Password: admin

## Connecting to PostgreSQL

When connecting to PostgreSQL from pgAdmin, use these connection details:

Host: postgres-srv
Port: 5432
Database: logdb
Username: postgres
Password: mysecretpassword (from your postgres-secret)

# When starting you need to run the following commands to apply the migrations

```bash
cd schemas/postgresql/log-aggregator && kubectl delete configmap postgres-migrations && kubectl create configmap postgres-migrations --from-file=logs.sql --from-file=alerts.sql && kubectl delete job postgres-migrations && kubectl apply -f ../../../infra/k8s/postgres-migrations-job.yaml
```
