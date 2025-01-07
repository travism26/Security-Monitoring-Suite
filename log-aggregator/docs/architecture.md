# Log Aggregator System Architecture

## Overview

The Log Aggregator is a robust, scalable system designed to collect, process, and analyze system metrics and logs from multiple sources. It provides real-time monitoring capabilities with an alert system for detecting anomalies in system performance.

## Core Components

### 1. Main Server (cmd/server/main.go)

The main server acts as the orchestrator for the entire system, responsible for:

- Initializing and coordinating all system components
- Managing database connections
- Setting up Kafka consumer
- Starting the HTTP server with middleware
- Handling graceful shutdowns

Key features:

- Configurable through YAML configuration
- Graceful shutdown handling
- Connection pool management for database
- Middleware support (CORS, Request ID, Logging, Recovery)

### 2. Domain Models (internal/domain/)

The domain package contains the core business logic and entities:

#### Log Entity

- Represents system metrics and log data
- Contains fields for:
  - System metrics (CPU, Memory usage)
  - Process information
  - Metadata and tags
  - Timestamps and identifiers
- Supports enrichment with additional context

#### Alert Entity

- Represents system alerts based on metric thresholds
- Tracks:
  - Alert severity levels
  - Alert status (Open/Resolved)
  - Related logs and metadata
  - Timestamps for creation/resolution

### 3. Kafka Integration (internal/kafka/)

The Kafka consumer component handles:

- Real-time message consumption from configured topics
- Message deserialization and validation
- Processing of system metrics and logs
- Integration with alert and log services

Features:

- Robust error handling
- Message batching for efficiency
- Automatic reconnection handling
- Parallel processing of partitions

### 4. Alert Service (internal/service/alert_service.go)

The alert service provides:

- Real-time monitoring of system metrics
- Configurable thresholds for:
  - CPU usage (default: 80%)
  - Memory usage (default: 85%)
  - Process count (default: 1000)
- Alert generation and management
- Trend analysis and reporting

Key features:

- Customizable alert thresholds
- Multiple severity levels
- Status management (Open/Resolved)
- Trend analysis capabilities

## Data Flow

1. **Data Ingestion**

   - System metrics and logs are received through Kafka
   - Messages are validated and parsed
   - Data is enriched with additional context

2. **Processing**

   - Metrics are analyzed against thresholds
   - Alerts are generated if thresholds are exceeded
   - Logs are stored in the database
   - Process information is tracked

3. **Storage**

   - PostgreSQL database stores:
     - System logs
     - Process information
     - Generated alerts
     - Historical metrics

4. **API Access**
   - RESTful API endpoints for:
     - Retrieving logs and metrics
     - Managing alerts
     - Accessing trend analysis
     - Health monitoring

## Design Decisions

### 1. Architecture Patterns

- Clean Architecture principles
- Domain-Driven Design (DDD) concepts
- Repository pattern for data access
- Service layer for business logic

### 2. Scalability Considerations

- Connection pooling for database
- Kafka consumer group support
- Configurable batch processing
- Parallel processing capabilities

### 3. Reliability Features

- Graceful shutdown handling
- Error recovery middleware
- Robust error handling
- Transaction support for data consistency

### 4. Monitoring and Maintenance

- Health check endpoints
- Detailed logging
- Request ID tracking
- Performance metrics collection

## Configuration

The system is highly configurable through:

- YAML configuration files
- Environment variables
- Runtime adjustable thresholds
- API key authentication

## Deployment

The application is containerized and can be deployed using:

- Kubernetes (k8s) configurations provided
- Docker containers
- Support for different environments (dev, staging, prod)

## Security Considerations

- API key authentication
- CORS middleware
- Request validation
- Secure database connections
- Error handling that prevents information leakage

## Future Enhancements

1. **Scalability**

   - Implement caching layer
   - Add support for multiple database types
   - Enhance batch processing capabilities

2. **Monitoring**

   - Add more sophisticated alerting rules
   - Implement machine learning for anomaly detection
   - Enhance trend analysis capabilities

3. **Integration**
   - Add support for more message queue systems
   - Implement additional alert notification channels
   - Add support for more metric collection sources
