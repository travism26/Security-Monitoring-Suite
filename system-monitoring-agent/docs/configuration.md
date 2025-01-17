# System Monitoring Agent Configuration

The agent's behavior can be configured through a YAML configuration file. Below is the complete reference of all available configuration options.

## Main Configuration

### LogFilePath

- **Type**: String
- **Default**: './agent.log'
- **Description**: Path to the log file where agent operations are recorded
- **Example**: '/var/log/monitoring-agent/agent.log'

### Interval

- **Type**: Integer
- **Default**: 10
- **Description**: Collection interval in seconds
- **Constraints**: Must be between 5 and 300 seconds
- **Example**: 30

## Kafka Configuration

### Brokers

- **Type**: List of strings
- **Default**: ['localhost:9092']
- **Description**: List of Kafka broker addresses
- **Example**:
  ```yaml
  Brokers:
    - "kafka1:9092"
    - "kafka2:9092"
  ```

### Topic

- **Type**: String
- **Default**: 'system-metrics'
- **Description**: Kafka topic where metrics will be published
- **Example**: 'prod-system-metrics'

## HTTP Configuration

### Endpoint

- **Type**: String
- **Default**: 'http://localhost:30001/api/v1/system-metrics/ingest'
- **Description**: HTTP endpoint for direct metric submission
- **Note**: Alternative to Kafka-based submission

### StorageDir

- **Type**: String
- **Default**: './storage'
- **Description**: Directory for storing temporary data
- **Example**: '/var/lib/monitoring-agent/storage'

## Monitoring Options

### CPU

- **Type**: Boolean
- **Default**: true
- **Description**: Enable CPU monitoring

### Memory

- **Type**: Boolean
- **Default**: true
- **Description**: Enable memory monitoring

### Disk

- **Type**: Boolean
- **Default**: true
- **Description**: Enable disk monitoring

### Network

- **Type**: Boolean
- **Default**: true
- **Description**: Enable network monitoring

## Threshold Configuration

### CPU

- **Type**: Integer
- **Default**: 8
- **Description**: CPU usage percentage threshold for alerts
- **Constraints**: 0-100

### Memory

- **Type**: Integer
- **Default**: 85
- **Description**: Memory usage percentage threshold for alerts
- **Constraints**: 0-100

### Disk

- **Type**: Integer
- **Default**: 90
- **Description**: Disk usage percentage threshold for alerts
- **Constraints**: 0-100

## Example Configuration

```yaml
LogFilePath: "/var/log/monitoring-agent.log"
Interval: 30
Kafka:
  Brokers:
    - "kafka1:9092"
    - "kafka2:9092"
  Topic: "prod-system-metrics"
HTTP:
  Endpoint: "http://monitoring.example.com/api/v1/system-metrics/ingest"
  StorageDir: "/var/lib/monitoring-agent/storage"
Monitors:
  CPU: true
  Memory: true
  Disk: true
  Network: true
Thresholds:
  CPU: 80
  Memory: 90
  Disk: 95
```
