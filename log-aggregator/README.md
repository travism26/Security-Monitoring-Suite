# Log Aggregator Service (Backend)

## Objective

Build a centralized log aggregation service to collect, store, and analyze logs from multiple sources. Implement basic SIEM-like features for detecting and analyzing potential security threats.

---

## Components

### 1. **Kafka Consumer for Log Aggregation**

- **Purpose**: Consume logs from Kafka topics and forward them for storage and analysis.
- **Tasks**:
  - Listen to Kafka topics (e.g., `system-metrics`, `threat-logs`).
  - Process logs for normalization and enrichment.
  - Forward processed logs to a database.

---

### 2. **Database for Log Storage**

- **Options**: MongoDB, Elasticsearch, or PostgreSQL.
- **Schema Design**:
  - Example fields for storing logs:
    - `timestamp`: When the log was generated.
    - `host`: Hostname or IP address of the source.
    - `cpu_usage`: CPU usage percentage.
    - `memory_usage`: Memory usage details (used, total, percentage).
    - `threats`: List of detected threats with severity and tags.
- **Tasks**:
  - Set up the database schema for logs and metrics.
  - Implement functions for storing and retrieving data.

---

### 3. **Basic SIEM-Like Features**

- **Rule-Based Detection**:
  - Implement simple rules to detect anomalies or threats.
  - Example rules:
    - "CPU usage exceeds 90%."
    - "Known malicious process detected."
- **Severity Levels**:
  - Assign severity levels to detected threats:
    - LOW, MEDIUM, HIGH, CRITICAL.
- **Trend Analysis**:
  - Aggregate metrics over time to identify patterns or trends.
- **Anomaly Detection**:
  - Implement threshold-based detection for CPU, memory, or network anomalies.

---

### 4. **REST API for Logs**

- **Purpose**: Provide endpoints for querying and analyzing logs.
- **Endpoints**:
  - `GET /logs` - Retrieve all logs with pagination and filtering.
  - `GET /logs/:id` - Fetch a specific log by ID.
  - `GET /alerts` - Retrieve a list of triggered alerts.
  - `GET /metrics/summary` - Provide aggregated metrics (e.g., average CPU usage, memory trends).

---

## Design Considerations

### Performance

- Batch-process Kafka messages to optimize consumer throughput.
- Use indexing in the database for faster log retrieval.

### Scalability

- Ensure the Kafka setup and database can handle large log volumes.
- Use a scalable database like Elasticsearch if querying is a priority.

### Security

- Implement secure access to the REST API using authentication (e.g., JWT).
- Encrypt sensitive data in transit and at rest.
- Implement CORS policies to restrict access to authorized UI clients.

---

## Deliverables

1. **Kafka Consumer**:
   - Processes logs from Kafka topics.
2. **Database**:
   - Stores aggregated logs with support for querying.
3. **REST API**:
   - Exposes logs, metrics, and alerts via endpoints.
   - Provides authentication and authorization.
   - Implements CORS policies.

## Related Repositories

- Dashboard UI: [link-to-dashboard-repo] - Web interface for visualizing logs, metrics, and alerts.

## Future Enhancements

- Integrate machine learning for advanced anomaly detection.
- Add support for external alerting systems (e.g., Slack, PagerDuty).
- Include role-based access control (RBAC) for securing log visibility.
- Implement rate limiting and request quotas.
- Add API versioning support.
