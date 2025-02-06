# CyberSecurity-Toolset

A comprehensive suite of security and monitoring tools designed to enhance system performance monitoring, security threat detection, and network analysis capabilities. This repository contains a collection of projects aimed at providing robust security insights and actionable monitoring for various environments, focusing on Windows, macOS systems, and Go-based solutions.

## Project Status: 🚧 Active Development

### Current Progress:

- ✅ **[System Monitoring Agent](./system-monitoring-agent/)**: Completed
- ✅ **[System Monitoring Gateway](./system-monitoring-gateway/)**: Completed
- ✅ **[SIEM Dashboard](./siem-dashboard/)**: Completed
- ✅ **[Log Aggregator](./log-aggregator/)**: Completed
- ✅ **Threat Detection Simulation**: Completed
- 🚧 **Mini XDR System**: In Development
- 📝 **Network Protocol Analyzer**: Planning Phase

---

## Project Overview

The CyberSecurity-Toolset repository is developed as an integrated security suite, with each component designed to work both independently and as part of the larger ecosystem. The system supports multi-tenancy, allowing secure isolation of data and resources between different organizations or teams. Each sub-folder represents a unique project with its own README, instructions, and usage guidelines.

---

## Components:

1. **System Monitoring Agent**

   - Lightweight Go agent for real-time system performance monitoring
   - Collects CPU, memory, disk, network, and process data
   - Cross-platform support (Windows, macOS)
   - Sends metrics to the System Monitoring Gateway
   - Kubernetes-ready deployment options

2. **System Monitoring Gateway**

   - TypeScript/Node.js backend service
   - Handles authentication and authorization
   - Multi-tenant architecture support
   - API key management for secure agent communication
   - Collects and validates metrics from the Monitoring Agent
   - Publishes metrics as Kafka events for downstream consumers
   - MongoDB integration for user and configuration management

3. **SIEM Dashboard**

   - Next.js-based frontend application
   - Real-time monitoring and analytics
   - Incident response management
   - Alert configuration and management
   - API key management interface
   - Threat analysis and visualization
   - User and tenant management

4. **Log Aggregator**

   - Centralizes logs from multiple sources
   - PostgreSQL-based storage with multi-tenancy support
   - Real-time log processing and analysis
   - Alert generation based on log patterns
   - REST API for log querying and management
   - Kafka integration for real-time event processing

5. **Threat Detection Simulation**

   - Simulates common security threats
   - Validates detection mechanisms
   - Integrates with Monitoring Agent data
   - Tests incident response pipelines

6. **Mini XDR System** (In Development)

   - Correlates events across all toolset components
   - Automates incident detection and response
   - Integrates threat intelligence feeds
   - Advanced threat hunting capabilities

7. **Network Protocol Analyzer** (Planned)
   - Deep packet inspection and traffic analysis
   - Anomaly detection in connection patterns
   - Advanced persistent threat detection
   - Malware C2 channel identification

---

## Infrastructure

The suite utilizes modern cloud-native technologies:

- Kubernetes for container orchestration
- Kafka for event streaming and processing
- PostgreSQL for log storage and analysis
- MongoDB for user and configuration management
- Ingress controllers for routing and load balancing

---

## Development Goals

This project serves as both a practical security toolkit and a learning journey into:

- Security monitoring and detection tools
- Multi-tenant architecture design
- Windows and macOS system integration
- Efficient and lightweight Go-based agent design
- Building Kubernetes-ready applications
- Kafka-powered real-time data pipelines
- Threat simulation and response systems
- Advanced network analysis and protocol inspection

---

### Future Improvements

- Enhanced threat intelligence integration
- Machine learning-based anomaly detection
- Advanced correlation rules for XDR
- Extended API capabilities
- Additional dashboard visualizations
- Advanced multi-tenancy features
- Automated incident response workflows
