# Troubleshooting Kubernetes Health Endpoints

This guide explains how to troubleshoot and fix authentication issues with health endpoints in a Go service using Gin framework, particularly when deployed in Kubernetes.

## Problem: Health Endpoints Returning 401 Unauthorized

### Symptoms

1. Kubernetes pods failing readiness/liveness probes
2. Logs showing 401 Unauthorized errors for /health endpoint:
   ```
   [GIN] 2025/02/01 - 16:30:02 | 401 | 103.5μs | 10.1.0.1 | GET "/health"
   ```
3. Pods restarting frequently due to failed health checks
4. Debug logs showing too many middleware handlers on health endpoints

### Diagnosis Steps

1. Check the number of middleware handlers on routes:

   ```bash
   kubectl logs -l app=log-aggregator
   ```

   Look for output like:

   ```
   [GIN-debug] GET /health --> handler.healthCheck (8 handlers)  # Too many handlers
   ```

   Health endpoints should have minimal middleware (typically 1-2 handlers).

2. Verify Kubernetes probe configuration:

   ```bash
   kubectl describe pod -l app=log-aggregator
   ```

   Look for:

   ```yaml
   Liveness: http-get http://:8080/health
   Readiness: http-get http://:8080/health
   ```

3. Check if authentication middleware is being applied globally:
   ```go
   // Problem: All routes including health get auth middleware
   router.Use(authMiddleware)
   router.GET("/health", healthCheck)
   ```

## Solution

### 1. Restructure Router Setup

Modify your main.go to separate health endpoints from authenticated routes:

```go
// Initialize HTTP server with minimal middleware
router := gin.New()
router.Use(gin.Recovery()) // Only essential middleware

// Register health check endpoints first (no auth required)
handler.RegisterHealthRoutes(router)

// Create API router group with full middleware stack
apiRouter := router.Group("/api/v1")
apiRouter.Use(
    middleware.CORS(),
    middleware.RequestID(),
    middleware.Logger(),
    middleware.Recovery(),
    middleware.Tenant(),
)

// Register API routes on the authenticated router group
logs := apiRouter.Group("/logs")
{
    logs.GET("", logHandler.ListLogs)
    // ... other routes
}
```

### 2. Update Health Handler

Ensure your health handler is simple and doesn't require authentication:

```go
func RegisterHealthRoutes(router gin.IRouter) {
    router.GET("/health", healthCheck)
    router.GET("/readiness", readinessCheck)
}

func healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func readinessCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
```

### 3. Rebuild and Deploy

1. Build new Docker image:

   ```bash
   cd log-aggregator
   docker build -t your-registry/log-aggregator:latest .
   docker push your-registry/log-aggregator:latest
   ```

2. Update Kubernetes deployment:

   ```bash
   kubectl rollout restart deployment log-aggregator-depl
   ```

3. Verify the fix:

   ```bash
   # Watch pod status
   kubectl get pods -l app=log-aggregator -w

   # Check logs for middleware handlers
   kubectl logs -l app=log-aggregator
   ```

   Look for:

   ```
   [GIN-debug] GET /health --> handler.healthCheck (2 handlers)  # Correct
   ```

### 4. Troubleshooting Deployment Issues

If you encounter issues during deployment:

1. Clean up stuck pods:

   ```bash
   kubectl delete pod -l app=log-aggregator --field-selector status.phase=Pending
   ```

2. Scale deployment down and up:

   ```bash
   kubectl scale deployment log-aggregator-depl --replicas=0
   kubectl scale deployment log-aggregator-depl --replicas=1
   ```

3. Check pod events:
   ```bash
   kubectl describe pod -l app=log-aggregator
   ```

## Best Practices

1. **Minimal Health Endpoints**: Health endpoints should be simple and have minimal middleware.
2. **Separate Route Groups**: Use route groups to separate authenticated and non-authenticated endpoints.
3. **Explicit Middleware**: Be explicit about which middleware applies to which routes.
4. **Logging**: Include enough logging to debug middleware and authentication issues.
5. **Kubernetes Probes**: Configure appropriate timing for liveness and readiness probes.

## Common Mistakes to Avoid

1. Applying authentication middleware globally
2. Using gin.Default() when you need fine-grained middleware control
3. Not testing health endpoints independently of authenticated routes
4. Overcomplicating health check logic

## Verification

After implementing the fix, verify:

1. Health endpoints return 200 OK without authentication
2. API endpoints still require proper authentication
3. Kubernetes probes succeed
4. Pods remain stable without restart loops
5. Logs show correct number of middleware handlers for each route

Remember to always test health endpoints both locally and in Kubernetes environment after making middleware changes.

## Testing Service Endpoints

### Prerequisites

- `curl` installed for HTTP requests
- `jq` installed for JSON formatting (optional but recommended)
- Valid API key for authenticated endpoints
- Access to the Kubernetes cluster

### 1. Testing Gateway Service Endpoints

#### Health and Metrics (Internal)

```bash
# Test health endpoint
curl -i https://security.dev/health

# Test metrics endpoint
curl -i https://security.dev/metrics
```

Expected response for health: HTTP 200 with `{"status": "OK"}`

#### API Endpoints (Authenticated)

```bash
# Test system metrics ingestion
curl -i -X POST https://security.dev/gateway/api/v1/system/metrics/ingest \
  -H "Content-Type: application/json" \
  -H "x-api-key: YOUR_API_KEY" \
  -d '{"data": {"metrics": {...}}, "timestamp": "2025-02-02T00:00:00Z"}'

# Test API keys endpoint
curl -i https://security.dev/gateway/api/v1/keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 2. Testing Log Aggregator Endpoints

#### Health Endpoints (Internal)

```bash
# Test health endpoint
curl -i https://security.dev/health/logs

# Test readiness endpoint
curl -i https://security.dev/readiness
```

Expected response: HTTP 200 with `{"status": "healthy"}` or `{"status": "ready"}`

#### API Endpoints (Authenticated)

```bash
# List logs
curl -i https://security.dev/logs/api/v1/logs \
  -H "x-api-key: YOUR_API_KEY"

# Store a log
curl -i -X POST https://security.dev/logs/api/v1/logs \
  -H "Content-Type: application/json" \
  -H "x-api-key: YOUR_API_KEY" \
  -d '{"message": "test log", "level": "info"}'

# List alerts
curl -i https://security.dev/logs/api/v1/alerts \
  -H "x-api-key: YOUR_API_KEY"
```

### 3. Troubleshooting Common Issues

1. **401 Unauthorized**

   - Verify API key or JWT token is valid
   - Check if token is properly formatted in header
   - Ensure you're using the correct authentication method for the endpoint

2. **404 Not Found**

   - Verify the URL path is correct with proper prefix
   - Check ingress configuration is properly routing the path
   - Ensure service and pods are running

3. **503 Service Unavailable**
   - Check if pods are running and ready
   - Verify Kafka connection for metrics ingestion
   - Check database connectivity for log aggregator

### 4. Kubernetes Verification

```bash
# Check pod status
kubectl get pods

# Check ingress configuration
kubectl describe ingress ingress-srv

# View service logs
kubectl logs -l app=system-monitoring-gateway
kubectl logs -l app=log-aggregator

# Port forward for local testing
kubectl port-forward svc/system-monitoring-gateway 3000:3000
kubectl port-forward svc/log-aggregator-srv 8080:8080
```

### 5. Best Practices for Testing

1. Test health endpoints first to ensure basic connectivity
2. Test unauthenticated endpoints before authenticated ones
3. Use proper HTTP methods (GET, POST, PUT, DELETE)
4. Include all required headers
5. Test with invalid authentication to verify security
6. Monitor logs while testing to catch issues
7. Test both success and error cases
8. Verify proper error responses and status codes

Remember to replace placeholders (YOUR_API_KEY, YOUR_JWT_TOKEN) with actual values, and adjust the security.dev domain to match your environment.
