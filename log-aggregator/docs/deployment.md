# Log Aggregator Deployment Guide

## Overview

This document details the deployment process and infrastructure requirements for the Log Aggregator system. The application is containerized and designed to be deployed in a Kubernetes environment, with support for different deployment environments (development, staging, production).

## Prerequisites

### Infrastructure Requirements

1. **Kubernetes Cluster**

   - Minimum version: 1.19+
   - Resource requirements:
     - CPU: 2 cores minimum
     - Memory: 4GB minimum
     - Storage: 20GB minimum

2. **Dependencies**
   - PostgreSQL (11+)
   - Kafka (2.8+)
   - Redis (optional, for caching)

### Configuration Requirements

1. **Environment Variables**

   - Database credentials
   - Kafka configuration
   - API keys
   - Service endpoints

2. **Kubernetes Secrets**
   - Database credentials
   - API keys
   - TLS certificates

## Kubernetes Deployment

### Core Components

1. **Log Aggregator Service**

   ```yaml
   # log-aggregator-depl.yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: log-aggregator
   spec:
     replicas: 3
     selector:
       matchLabels:
         app: log-aggregator
     template:
       metadata:
         labels:
           app: log-aggregator
       spec:
         containers:
           - name: log-aggregator
             image: log-aggregator:latest
             ports:
               - containerPort: 8080
             env:
               - name: DB_HOST
                 valueFrom:
                   configMapKeyRef:
                     name: postgres-config
                     key: host
   ```

2. **PostgreSQL Database**
   ```yaml
   # log-aggregator-postgres-depl.yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: postgres
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: postgres
     template:
       metadata:
         labels:
           app: postgres
       spec:
         containers:
           - name: postgres
             image: postgres:13
             ports:
               - containerPort: 5432
             env:
               - name: POSTGRES_DB
                 value: logdb
   ```

### Supporting Services

1. **Kafka Setup**

   ```yaml
   # kafka-depl.yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: kafka
   spec:
     replicas: 3
     template:
       spec:
         containers:
           - name: kafka
             image: confluentinc/cp-kafka:latest
   ```

2. **Monitoring Tools**
   - Prometheus for metrics collection
   - Grafana for visualization
   - PgAdmin for database management

## Database Migration

### Initial Setup

1. **Migration Job**

   ```yaml
   # postgres-migrations-job.yaml
   apiVersion: batch/v1
   kind: Job
   metadata:
     name: postgres-migrations
   spec:
     template:
       spec:
         containers:
           - name: migrations
             image: log-aggregator-migrations:latest
             env:
               - name: DB_HOST
                 valueFrom:
                   configMapKeyRef:
                     name: postgres-config
                     key: host
         restartPolicy: Never
   ```

2. **Migration Scripts**
   - Located in `/migrations` directory
   - Versioned incrementally
   - Includes rollback procedures

## Scaling Considerations

### Horizontal Scaling

1. **Log Aggregator Service**

   - Scale based on CPU/Memory utilization
   - Configure HPA (Horizontal Pod Autoscaler)

   ```yaml
   apiVersion: autoscaling/v2beta1
   kind: HorizontalPodAutoscaler
   metadata:
     name: log-aggregator-hpa
   spec:
     scaleTargetRef:
       apiVersion: apps/v1
       kind: Deployment
       name: log-aggregator
     minReplicas: 3
     maxReplicas: 10
     metrics:
       - type: Resource
         resource:
           name: cpu
           targetAverageUtilization: 80
   ```

2. **Kafka Scaling**
   - Partition management
   - Consumer group configuration
   - Topic replication factors

### Vertical Scaling

1. **Resource Limits**

   ```yaml
   resources:
     requests:
       memory: "256Mi"
       cpu: "250m"
     limits:
       memory: "512Mi"
       cpu: "500m"
   ```

2. **Database Scaling**
   - Connection pool configuration
   - Resource allocation
   - Query optimization

## Monitoring Setup

### Health Checks

1. **Liveness Probe**

   ```yaml
   livenessProbe:
     httpGet:
       path: /health
       port: 8080
     initialDelaySeconds: 30
     periodSeconds: 10
   ```

2. **Readiness Probe**
   ```yaml
   readinessProbe:
     httpGet:
       path: /health/ready
       port: 8080
     initialDelaySeconds: 5
     periodSeconds: 5
   ```

### Metrics Collection

1. **Prometheus Integration**

   - System metrics
   - Application metrics
   - Custom alerts

2. **Logging**
   - Centralized logging
   - Log aggregation
   - Search and analysis

## Security Configuration

### Network Policies

1. **Service Isolation**

   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: log-aggregator-network-policy
   spec:
     podSelector:
       matchLabels:
         app: log-aggregator
     policyTypes:
       - Ingress
       - Egress
   ```

2. **TLS Configuration**
   - Certificate management
   - Secret rotation
   - Secure communication

### Access Control

1. **RBAC Configuration**

   - Service accounts
   - Role bindings
   - Access policies

2. **API Security**
   - Authentication
   - Rate limiting
   - Input validation

## Backup and Recovery

### Database Backups

1. **Automated Backups**

   - Daily snapshots
   - Transaction logs
   - Retention policies

2. **Recovery Procedures**
   - Point-in-time recovery
   - Disaster recovery
   - Backup verification

### System State

1. **Configuration Backups**

   - Kubernetes manifests
   - Application configs
   - Secrets management

2. **State Recovery**
   - Service restoration
   - Data consistency
   - Validation procedures

## Environment-Specific Configurations

### Development

```yaml
environment: development
logLevel: debug
kafka:
  replicas: 1
database:
  replicas: 1
```

### Staging

```yaml
environment: staging
logLevel: info
kafka:
  replicas: 2
database:
  replicas: 2
```

### Production

```yaml
environment: production
logLevel: warn
kafka:
  replicas: 3
database:
  replicas: 3
```

## Troubleshooting Guide

### Common Issues

1. **Database Connectivity**

   - Connection pool exhaustion
   - Network latency
   - Authentication failures

2. **Kafka Issues**
   - Consumer lag
   - Partition rebalancing
   - Message processing delays

### Resolution Steps

1. **Service Recovery**

   - Pod restart procedures
   - Log analysis
   - State verification

2. **Data Recovery**
   - Backup restoration
   - Data validation
   - Service reconciliation
