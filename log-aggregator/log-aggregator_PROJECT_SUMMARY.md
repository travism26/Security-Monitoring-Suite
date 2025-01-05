# Log Aggregator with SIEM-Like Features: Project Summary

## 1. Project Overview

The **Log Aggregator with SIEM-Like Features** is a system designed to pull system and security metrics from a Kafka topic, store them in a relational database, and provide basic SIEM (Security Information and Event Management) functionalities. By centralizing logs from multiple sources, it enables early threat detection, trend analysis, and real-time alerting.

---

## 2. Technical Architecture

### 2.1 Backend Stack

- **Language/Framework**: Go (Golang) for core log aggregation service
- **Message Broker**: Kafka for real-time ingestion of log data
- **Database**: PostgreSQL for structured storage (tables: `logs`, `process_logs`, and optional `alerts`)
- **API Style**: RESTful endpoints for querying and retrieving logs and metrics
- **Authentication**: JSON Web Tokens (JWT) or Basic Auth (configurable)
- **Testing**: Go’s built-in testing framework, plus additional libraries (e.g., testify)

### 2.2 SIEM-Like Features

- **Rule-Based Detection**: High CPU usage, memory usage thresholds, malicious processes
- **Alerts & Severity**: Generate alerts with severity levels (LOW, MEDIUM, HIGH, CRITICAL)
- **Trend Analysis**: Basic queries to analyze CPU/memory trends or process anomalies over time
- **Extended Integrations**: Future scope for security event correlation or ML-based anomaly detection

### 2.3 Data Pipeline & Services

- **Data Ingestion**:
  - Multiple producers (e.g., a system-monitoring-agent) publish metrics to Kafka topics
  - Log Aggregator consumes these metrics and processes them
- **Normalization & Enrichment**:
  - Standardize log format and add metadata (e.g., host info, severity scoring)
- **Storage Layer**:
  - PostgreSQL schemas for fast querying and indexing (timestamp, host, process name)
  - Optionally extend to Elasticsearch for full-text search and advanced analytics
- **Alert Management**:
  - Store triggered alerts (optional table for easy reference and auditing)
  - Future integration with real-time notification systems (e.g., Slack, email)

### 2.4 DevOps and Infrastructure

- **Containerization**: Docker for building and running the aggregator
- **Orchestration**: Docker Compose for local dev; Kubernetes (optional for production scaling)
- **CI/CD**: GitHub Actions (lint, unit tests, integration tests)
- **Monitoring & Metrics**:
  - Prometheus and Grafana for service metrics
  - Logs viewable in Kibana if Elasticsearch is used
- **Logging**: Leveraging the aggregator’s own pipeline for inbound logs and a separate pipeline for system logs

### 2.5 Cross-Platform Considerations

- **Development Environment**: Compatible with Linux, macOS, or Windows
- **Command Line**: Cross-platform builds using Makefile or shell scripts
- **Network Configuration**: Kafka and PostgreSQL accessible locally or via Docker networks
- **File System**: Ensure correct handling of file paths and environment variables across OSes

### 2.6 Expanded Feature Set

- **Alert Retention Policies**: Regular cleanups or archiving for older alerts
- **Role-Based Access Control (RBAC)**: Limit log and alert visibility to authorized roles
- **Visualization Integration**: Provide optional web UI or integrate with third-party dashboards

---

## 3. Data Flow Architecture

```mermaid
flowchart TD
    A[System-Monitoring-Agent(s)] -->|Produce Metrics| B[Kafka Topic(s)]
    B -->|Consume Messages| C[Log Aggregator]
    C -->|Store Logs/Alerts| D[PostgreSQL]
    C -->|API| E[REST Endpoints]
    E -->|Query Logs/Alerts| F[Frontend/Dashboard]
    C -->|Optional| G[Alert Notification Service]
```

## 4. Performance Optimization

### 4.1 Aggregator Optimization

- Batch Consumption: Process multiple Kafka messages at once to reduce overhead.
- Indexing Strategy: Ensure PostgreSQL indexes are created for commonly queried fields (e.g., timestamp, host).
- Connection Pooling: Optimize database connections for high-volume writes.

### 4.2 Analysis Optimization

- Caching Layer: Use Redis for short-term caching of frequently queried data (e.g., summary statistics).
- Asynchronous Processing: Offload heavy analysis tasks to background workers.
- Load Balancing: Deploy multiple instances of the Log Aggregator behind a load balancer to handle large workloads.

## 5. Future Roadmap

### 5.1 Short-Term Goals

- Complete REST API for querying logs and alerts.
- Implement basic alerting system with configurable rules.
- Create initial API documentation for ease of use.

### 5.2 Mid-Term Goals

- Add support for more complex rule-based detection and correlations.
- Integrate notification services (e.g., email, Slack) for critical alerts.
- Build a basic frontend dashboard for log visualization.

### 5.3 Long-Term Vision

- Enhance threat detection with machine learning for anomaly detection.
- Expand log aggregation capabilities to handle multi-tenant environments.
- Introduce integrations with third-party SIEM tools (e.g., Splunk, Datadog).
