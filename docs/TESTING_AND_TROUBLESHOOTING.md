# Security Monitoring Suite - Testing & Troubleshooting Guide

This guide provides instructions for testing and troubleshooting the Security Monitoring Suite services through the ingress controller.

## Prerequisites

1. Add the following entry to your hosts file (`/etc/hosts` on Unix-like systems):

   ```
   127.0.0.1 security.dev
   ```

2. Ensure all services are running:
   ```bash
   kubectl get pods
   ```

## Testing Services

### Security Dashboard (Frontend)

```bash
# Access the dashboard (redirects to login)
curl -k https://security.dev/

# With verbose output to check TLS and headers
curl -k -v https://security.dev/
```

Expected response: 307 redirect to /login with proper security headers

### System Monitoring Gateway

```bash
# Test metrics endpoint (requires authentication)
curl -k -X GET https://security.dev/api/v1/metrics \
  -H "Authorization: Bearer YOUR_API_KEY"

# Test health endpoint
curl -k -X GET https://security.dev/health/gateway \
  -H "Authorization: Bearer YOUR_API_KEY"

# Submit metrics (requires authentication)
curl -k -X POST https://security.dev/api/v1/metrics \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "cpu_usage": 45.2,
    "memory_usage": 1024,
    "timestamp": "2025-02-01T22:00:00Z"
  }'
```

Expected responses:

- 401 Unauthorized without API key
- 200 OK with valid API key
- JSON response with metrics data

### Log Aggregator

```bash
# Query logs with time range (requires authentication)
curl -k -X GET "https://security.dev/api/v1/logs/range?start=2025-02-01T00:00:00Z&end=2025-02-01T23:59:59Z" \
  -H "Authorization: Bearer YOUR_API_KEY"

# Submit logs (requires authentication)
curl -k -X POST https://security.dev/api/v1/logs \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "level": "error",
    "message": "Test error message",
    "timestamp": "2025-02-01T22:00:00Z",
    "source": "test-service"
  }'

# Submit batch logs
curl -k -X POST https://security.dev/api/v1/logs/batch \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "level": "info",
      "message": "Test info message",
      "timestamp": "2025-02-01T22:00:00Z",
      "source": "test-service"
    },
    {
      "level": "error",
      "message": "Test error message",
      "timestamp": "2025-02-01T22:00:00Z",
      "source": "test-service"
    }
  ]'

# Query alerts
curl -k -X GET https://security.dev/api/v1/alerts \
  -H "Authorization: Bearer YOUR_API_KEY"

# Get alert by ID
curl -k -X GET https://security.dev/api/v1/alerts/123 \
  -H "Authorization: Bearer YOUR_API_KEY"

# Update alert status
curl -k -X PUT https://security.dev/api/v1/alerts/123/status \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status": "resolved"}'
```

## Common Issues and Solutions

### TLS Certificate Issues

**Symptom**: SSL certificate errors when accessing services

**Solution**:

1. Verify the TLS secret exists:

   ```bash
   kubectl get secret tls-secret
   ```

2. Recreate the certificate if needed:

   ```bash
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
     -keyout security.dev.key -out security.dev.crt \
     -subj "/CN=security.dev" \
     -addext "subjectAltName = DNS:security.dev"

   kubectl delete secret tls-secret
   kubectl create secret tls tls-secret --key=security.dev.key --cert=security.dev.crt
   ```

### Authentication Issues

**Symptom**: Getting 401 Unauthorized errors

**Solution**:

1. Verify your API key is valid
2. Check the API key format in the Authorization header
3. Ensure the API key has the correct permissions

### Service Unavailable (503) Errors

**Symptom**: Services returning 503 errors

**Solution**:

1. Check if pods are running:

   ```bash
   kubectl get pods
   ```

2. Check pod logs:

   ```bash
   kubectl logs -l app=siem-dashboard
   kubectl logs -l app=system-monitoring-gateway
   kubectl logs -l app=log-aggregator
   ```

3. Check service endpoints:
   ```bash
   kubectl get endpoints siem-dashboard-srv
   kubectl get endpoints system-monitoring-gateway
   kubectl get endpoints log-aggregator-srv
   ```

### Ingress Issues

**Symptom**: Routing or access problems

**Solution**:

1. Check ingress configuration:

   ```bash
   kubectl get ingress ingress-srv -o yaml
   ```

2. Check ingress controller logs:

   ```bash
   kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller
   ```

3. Verify ingress controller is running:
   ```bash
   kubectl get pods -n ingress-nginx
   ```

## Monitoring and Debugging

### Check Service Status

```bash
# Get all resources
kubectl get all

# Get specific service logs
kubectl logs -l app=siem-dashboard --tail=100
kubectl logs -l app=system-monitoring-gateway --tail=100
kubectl logs -l app=log-aggregator --tail=100
```

### Monitor Ingress Traffic

```bash
# Watch ingress controller logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller -f

# Get ingress controller metrics
curl -k https://security.dev/nginx_status
```

### Debug Network Issues

```bash
# Test DNS resolution
ping security.dev

# Check service connectivity
kubectl exec -it $(kubectl get pod -l app=siem-dashboard -o jsonpath='{.items[0].metadata.name}') -- curl -v http://system-monitoring-gateway:3000/health
```

## Performance Testing

### Load Testing Example

Using [hey](https://github.com/rakyll/hey) for load testing:

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Run load test (50 concurrent users, 1000 requests)
hey -n 1000 -c 50 -k https://security.dev/api/v1/metrics
```

## Best Practices

1. Always use HTTPS for all requests
2. Include proper headers:
   - `Authorization: Bearer YOUR_API_KEY`
   - `Content-Type: application/json`
3. Handle rate limiting appropriately (default: 300 requests per minute)
4. Monitor service health endpoints regularly
5. Keep TLS certificates up to date

## Getting Help

If you encounter issues:

1. Check the logs of the affected service
2. Verify all prerequisites are met
3. Ensure all services are running properly
4. Check the ingress controller logs for routing issues
5. Verify network policies and firewall rules
6. Review recent changes or deployments

For additional assistance, consult the following resources:

- [Ingress Nginx Controller Documentation](https://kubernetes.github.io/ingress-nginx/)
- [Kubernetes Ingress Documentation](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- Project-specific documentation in the `/docs` directory
