# Testing Multi-Tenancy Features in Kubernetes

This guide explains how to test the multi-tenancy features of the Log Aggregator service in a Kubernetes environment.

## Prerequisites

- A running Kubernetes cluster
- kubectl configured to access your cluster
- PostgreSQL database deployed and running
- Log Aggregator service deployed

## 1. Deploy the Services

First, ensure all required services are running:

```bash
# Apply PostgreSQL deployment and service
kubectl apply -f infra/k8s/log-aggregator-postgres-depl.yaml

# Apply database migrations
kubectl apply -f infra/k8s/postgres-migrations-configmap.yaml
kubectl apply -f infra/k8s/postgres-migrations-job.yaml

# Apply Log Aggregator deployment and service
kubectl apply -f infra/k8s/log-aggregator-depl.yaml

# Verify deployments
kubectl get pods
kubectl get services
```

## 2. Port Forward the Service

To access the Log Aggregator service locally:

```bash
kubectl port-forward service/log-aggregator-srv 8080:8080
```

## 3. Create Test Organizations

Use the following curl commands to create test organizations:

```bash
# Create a customer organization
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Customer Org",
    "description": "Test customer organization"
  }'

# Create an agent organization
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Agent Org",
    "description": "Test agent organization"
  }'
```

Save the returned organization IDs for the next steps.

## 4. Generate API Keys

Generate API keys for both organizations:

```bash
# Generate a customer API key
curl -X POST http://localhost:8080/api/v1/organizations/{customer-org-id}/apikeys \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Customer Key",
    "type": "customer"
  }'

# Generate an agent API key
curl -X POST http://localhost:8080/api/v1/organizations/{agent-org-id}/apikeys \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Agent Key",
    "type": "agent"
  }'
```

Save both API keys for testing.

## 5. Test Multi-Tenancy Features

### 5.1 Test Customer API Key

```bash
# List alerts for customer organization
curl http://localhost:8080/api/v1/alerts \
  -H "X-API-Key: {customer-api-key}"

# Attempt to access agent-only endpoint (should fail)
curl http://localhost:8080/api/v1/metrics \
  -H "X-API-Key: {customer-api-key}"
```

### 5.2 Test Agent API Key

```bash
# Submit metrics (agent-only endpoint)
curl -X POST http://localhost:8080/api/v1/metrics \
  -H "X-API-Key: {agent-api-key}" \
  -H "Content-Type: application/json" \
  -d '{
    "host": "test-host",
    "cpu_usage": 50.5,
    "memory_usage": 1024,
    "process_count": 100
  }'

# Attempt to list alerts (should fail)
curl http://localhost:8080/api/v1/alerts \
  -H "X-API-Key: {agent-api-key}"
```

### 5.3 Test Data Isolation

```bash
# Create alerts for both organizations
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "X-API-Key: {customer-1-api-key}" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Alert Org 1",
    "description": "Test alert for organization 1"
  }'

curl -X POST http://localhost:8080/api/v1/alerts \
  -H "X-API-Key: {customer-2-api-key}" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Alert Org 2",
    "description": "Test alert for organization 2"
  }'

# Verify each organization only sees their own alerts
curl http://localhost:8080/api/v1/alerts \
  -H "X-API-Key: {customer-1-api-key}"

curl http://localhost:8080/api/v1/alerts \
  -H "X-API-Key: {customer-2-api-key}"
```

## 6. Verify Headers and Logging

Check the response headers to verify tenant context:

```bash
curl -v http://localhost:8080/api/v1/alerts \
  -H "X-API-Key: {customer-api-key}"
```

Look for the following headers in the response:

- `X-Organization-ID`: Should match the organization ID
- `X-API-Key-Type`: Should be either "customer" or "agent"

## 7. Check Logs

View the Log Aggregator pod logs to verify tenant-related operations:

```bash
# Get pod name
kubectl get pods | grep log-aggregator

# View logs
kubectl logs {pod-name} -f
```

Look for log entries containing:

- Organization ID
- API key validation
- Tenant context initialization
- Access denied messages for incorrect key types

## 8. Clean Up

To clean up test data:

```bash
# Delete test organizations (this will cascade delete their API keys)
curl -X DELETE http://localhost:8080/api/v1/organizations/{org-id-1}
curl -X DELETE http://localhost:8080/api/v1/organizations/{org-id-2}
```

## Common Issues and Troubleshooting

1. API Key Invalid

   - Verify the API key hash matches what's in the database
   - Check the API key hasn't expired or been revoked
   - Ensure the key type matches the endpoint requirements

2. Missing Tenant Context

   - Verify the X-API-Key header is being sent
   - Check logs for API key validation errors
   - Ensure the organization ID exists in the database

3. Cross-Tenant Data Access
   - Verify repository queries include organization ID filters
   - Check service layer methods properly use tenant context
   - Review database indexes for organization ID columns

## Next Steps

After verifying basic multi-tenancy functionality, consider testing:

1. API key rotation
2. Rate limiting per tenant
3. Tenant-specific configurations
4. Cross-tenant isolation under load
5. Audit logging for tenant operations
