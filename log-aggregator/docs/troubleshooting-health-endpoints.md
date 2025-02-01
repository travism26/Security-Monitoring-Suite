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
