# Shared Configuration Service Technical Design

## 1. Technical Stack

### 1.1 Core Technologies

- **Language**: Go 1.21+
- **Framework**: Gin (REST API)
- **Primary Database**: PostgreSQL 15+
- **Cache**: Redis 7.0+
- **Authentication**: JWT + API Keys
- **Documentation**: OpenAPI 3.0
- **Testing**: Go testing + testify

### 1.2 Dependencies

```go
require (
    github.com/gin-gonic/gin v1.9.0
    github.com/go-redis/redis/v8 v8.11.5
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/lib/pq v1.10.9
    github.com/prometheus/client_golang v1.16.0
    github.com/spf13/viper v1.16.0
    gorm.io/gorm v1.25.0
    gorm.io/driver/postgres v1.5.0
)
```

## 2. Database Schema

### 2.1 PostgreSQL Tables

```sql
-- Configurations
CREATE TABLE configurations (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    namespace VARCHAR(100) NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(100) NOT NULL,
    metadata JSONB,
    UNIQUE(key, namespace)
);

-- Configuration History
CREATE TABLE configuration_history (
    id SERIAL PRIMARY KEY,
    configuration_id INT REFERENCES configurations(id),
    value JSONB NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100) NOT NULL
);

-- API Keys
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    roles TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE
);

-- Audit Logs
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    changes JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## 3. API Design

### 3.1 REST Endpoints

#### Configuration Management

```
GET    /api/v1/configs/{key}
POST   /api/v1/configs
PUT    /api/v1/configs/{key}
DELETE /api/v1/configs/{key}
GET    /api/v1/configs/namespace/{namespace}
GET    /api/v1/configs/{key}/history
POST   /api/v1/configs/{key}/rollback
```

#### Access Control

```
POST   /api/v1/auth/token
POST   /api/v1/apikeys
DELETE /api/v1/apikeys/{id}
GET    /api/v1/apikeys
```

#### Health & Metrics

```
GET    /health
GET    /metrics
```

### 3.2 WebSocket Endpoints

```
WS     /ws/v1/configs/updates
```

## 4. Component Design

### 4.1 Configuration Manager

```go
type ConfigManager interface {
    Get(ctx context.Context, key string) (*ConfigEntry, error)
    Set(ctx context.Context, entry *ConfigEntry) error
    Delete(ctx context.Context, key string) error
    List(ctx context.Context, namespace string) ([]*ConfigEntry, error)
    GetHistory(ctx context.Context, key string) ([]*ConfigHistory, error)
    Rollback(ctx context.Context, key string, version int) error
    Watch(ctx context.Context) (<-chan ConfigUpdate, error)
}
```

### 4.2 Cache Manager

```go
type CacheManager interface {
    Get(ctx context.Context, key string) (*ConfigEntry, error)
    Set(ctx context.Context, key string, entry *ConfigEntry) error
    Delete(ctx context.Context, key string) error
    Invalidate(ctx context.Context, pattern string) error
}

type RedisCacheManager struct {
    client *redis.Client
    ttl    time.Duration
}
```

### 4.3 Database Manager

```go
type DBManager struct {
    db *gorm.DB
}

func (d *DBManager) CreateConfig(ctx context.Context, cfg *ConfigEntry) error {
    return d.db.Transaction(func(tx *gorm.DB) error {
        // Store config and create history entry
    })
}
```

## 5. Security Implementation

### 5.1 Authentication Middleware

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT validation
        // API key validation
        // Role-based access control
    }
}
```

### 5.2 Encryption

```go
type Encryptor interface {
    Encrypt(data []byte) ([]byte, error)
    Decrypt(data []byte) ([]byte, error)
}

type AESEncryptor struct {
    key [32]byte
}
```

## 6. SDK Design

### 6.1 Go SDK

```go
type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
    cache      CacheManager
}

func (c *Client) GetConfig(ctx context.Context, key string) (*ConfigEntry, error) {
    // Check cache first
    // Fallback to API call
}
```

### 6.2 TypeScript SDK

```typescript
class ConfigClient {
  constructor(baseUrl: string, apiKey: string) {
    // Initialize client
  }

  async getConfig(key: string): Promise<ConfigEntry> {
    // Implementation
  }

  watchUpdates(): WebSocket {
    // WebSocket implementation
  }
}
```

## 7. Deployment

### 7.1 Kubernetes Resources

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: config-service
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: config-service
          image: config-service:latest
          ports:
            - containerPort: 8080
          env:
            - name: DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: config-service-cm
                  key: db_host
```

### 7.2 Database Resources

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres
  replicas: 3
  template:
    spec:
      containers:
        - name: postgres
          image: postgres:15
```

### 7.3 Cache Resources

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-cache
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: redis
          image: redis:7.0
```

## 8. Monitoring

### 8.1 Prometheus Metrics

```go
var (
    configRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "config_requests_total",
            Help: "Total number of configuration requests",
        },
        []string{"method", "status"},
    )

    cacheHitRatio = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "config_cache_hit_ratio",
            Help: "Cache hit ratio for configuration requests",
        },
        []string{"namespace"},
    )
)
```

### 8.2 Logging

```go
type LogEntry struct {
    Timestamp   time.Time     `json:"timestamp"`
    Level       string        `json:"level"`
    Message     string        `json:"message"`
    Operation   string        `json:"operation"`
    Key         string        `json:"key"`
    Namespace   string        `json:"namespace"`
    UserID      string        `json:"user_id"`
    Duration    time.Duration `json:"duration"`
    CacheHit    bool          `json:"cache_hit"`
    Error       string        `json:"error,omitempty"`
}
```

## 9. Testing Strategy

### 9.1 Unit Tests

```go
func TestConfigManager_Set(t *testing.T) {
    // Test cases for configuration management
}

func TestCacheManager_Get(t *testing.T) {
    // Test cases for cache operations
}
```

### 9.2 Integration Tests

```go
func TestConfigService_Integration(t *testing.T) {
    // End-to-end test scenarios
}
```

## 10. Implementation Tasks

### Phase 1: Core Service (2 weeks)

1. Set up project structure and dependencies
2. Implement PostgreSQL schema and migrations
3. Implement basic CRUD operations
4. Add authentication and authorization
5. Set up Redis caching
6. Create API documentation

### Phase 2: Advanced Features (2 weeks)

1. Implement real-time updates
2. Develop SDKs
3. Add monitoring and metrics
4. Implement security features

### Phase 3: Integration (2 weeks)

1. Create Kubernetes resources
2. Set up CI/CD pipeline
3. Write documentation
4. Conduct security audit

### Phase 4: Performance (1 week)

1. Optimize database queries
2. Fine-tune caching strategy
3. Load testing
4. Performance monitoring
