# Log Aggregator Components Detail

## Service Layer Components

### Alert Service

The Alert Service (`internal/service/alert_service.go`) is responsible for monitoring system metrics and generating alerts based on predefined thresholds.

#### Key Features:

1. **Configurable Thresholds**

   ```go
   type AlertThresholds struct {
       CPUUsagePercent    float64 // Default: 80%
       MemoryUsagePercent float64 // Default: 85%
       ProcessCount       int     // Default: 1000
   }
   ```

2. **Alert Generation**

   - Monitors CPU usage
   - Tracks memory consumption
   - Watches process count
   - Generates appropriate severity alerts

3. **Trend Analysis**
   - Tracks alert patterns
   - Provides statistical analysis
   - Offers time-based distribution data

### Log Service

The Log Service handles the processing and storage of system logs and metrics.

#### Capabilities:

1. **Log Enrichment**

   - Adds environmental context
   - Includes application metadata
   - Generates correlation IDs
   - Adds standardized tags

2. **Batch Processing**
   - Efficient bulk log storage
   - Transaction support
   - Error handling and retries

## Data Layer Components

### Kafka Consumer

The Kafka Consumer (`internal/kafka/consumer.go`) is designed to handle real-time data ingestion.

#### Features:

1. **Message Processing**

   ```go
   func (c *Consumer) processMessage(msg *sarama.ConsumerMessage) error {
       // Unmarshal raw message
       // Create log entry
       // Extract processes
       // Store data
       // Process metrics for alerts
   }
   ```

2. **Data Extraction**

   - Parses host information
   - Extracts system metrics
   - Processes threat indicators
   - Handles process information

3. **Error Handling**
   - Robust message validation
   - Error recovery mechanisms
   - Detailed error logging

### Domain Models

#### Log Entity

The Log entity (`internal/domain/log.go`) represents the core data structure:

```go
type Log struct {
    ID              string
    Timestamp       time.Time
    Host            string
    Message         string
    Level           string
    Metadata        map[string]interface{}
    ProcessCount    int
    TotalCPUPercent float64
    TotalMemoryUsage int64
    Environment     string
    Application     string
    Component       string
    CorrelationID   string
    Tags            []string
    EnrichedAt      time.Time
    ProcessedCount  int
}
```

#### Process Entity

Represents individual system processes with metrics:

- Process ID
- Name
- CPU usage
- Memory consumption
- Status
- Timestamps

## Infrastructure Components

### Database Layer

1. **Repositories**

   - Log Repository
   - Process Repository
   - Alert Repository

2. **Connection Management**
   - Connection pooling
   - Retry mechanisms
   - Transaction handling

### API Layer

1. **Handlers**

   - Log Handler
   - Alert Handler
   - Health Handler

2. **Middleware**
   - CORS support
   - Request ID generation
   - Logging
   - Error recovery

## Integration Points

### 1. Kafka Integration

- Consumer group support
- Topic management
- Message serialization/deserialization
- Error handling and recovery

### 2. Database Integration

- PostgreSQL connection management
- Query optimization
- Transaction management
- Connection pooling

### 3. API Integration

- RESTful endpoints
- Authentication
- Rate limiting
- Response formatting

## Configuration Management

### 1. Application Configuration

```yaml
server:
  host: localhost
  port: 8080

database:
  host: localhost
  port: 5432
  name: logdb
  user: postgres
  password: secret

kafka:
  brokers:
    - localhost:9092
  topic: system_metrics
  groupID: log_aggregator

logService:
  environment: production
  application: log_aggregator
  component: processor
```

### 2. Alert Configuration

- Configurable thresholds
- Severity levels
- Notification settings
- Alert rules

## Error Handling

### 1. Error Types

- Database errors
- Kafka errors
- Processing errors
- Validation errors

### 2. Recovery Mechanisms

- Automatic retries
- Circuit breaking
- Fallback strategies
- Error logging

## Monitoring and Metrics

### 1. Health Checks

- Database connectivity
- Kafka connectivity
- API availability
- System resources

### 2. Performance Metrics

- Processing latency
- Queue depth
- Error rates
- Resource utilization

## Security Measures

### 1. Authentication

- API key validation
- Token management
- Access control

### 2. Data Protection

- Input validation
- Output sanitization
- Secure connections
- Error information hiding
