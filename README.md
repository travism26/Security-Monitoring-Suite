# CyberSecurity-Toolset

A comprehensive suite of security and monitoring tools designed to enhance system performance monitoring, security threat detection, and network analysis capabilities. This repository contains a collection of projects aimed at providing robust security insights and actionable monitoring for various environments, focusing on Windows, macOS systems, and Go-based solutions.

## Project Status: 🚧 Work In Progress

### Current Progress:

- ✅ **[System Monitoring Agent](./system-monitoring-agent/)**: Completed
- ✅ **API Gateway for Metrics**: Completed
- 📝 **Threat Detection Simulation**: Planning Phase / In Progress
- 📝 **Log Aggregator with SIEM-Like Features**: Planning Phase
- 📝 **Mini XDR System**: Planning Phase
- 📝 **Network Protocol Analyzer**: Planning Phase

---

## Project Overview

The CyberSecurity-Toolset repository is being developed as an integrated security suite, with each component designed to work both independently and as part of the larger ecosystem. Each sub-folder represents a unique project with its own README, instructions, and usage guidelines.

---

## Components:

1. **System Monitoring Agent**

   - Lightweight Go agent for real-time system performance monitoring.
   - Collects CPU, memory, disk, network, and process data.
   - Cross-platform support (Windows, macOS).
   - Sends metrics to the API Gateway.
   - Kubernetes-ready deployment options.

2. **API Gateway for Metrics**

   - Collects and validates metrics from the Monitoring Agent.
   - Publishes metrics as Kafka events for downstream consumers.
   - Acts as the central hub for metric ingestion.

3. **Threat Detection Simulation**

   - Simulates common security threats (e.g., ransomware, privilege escalation).
   - Validates detection mechanisms and incident response pipelines.
   - Integrates with Monitoring Agent data.

4. **Log Aggregator with SIEM-Like Features**

   - Centralizes logs from multiple sources (applications, OS logs, network devices).
   - Performs anomaly detection using predefined rules or ML models.
   - Provides visualization for event analysis.

5. **Mini XDR System**

   - Correlates events across all toolset components (metrics, logs, simulated threats).
   - Automates incident detection and response.
   - Integrates threat intelligence feeds.

6. **Network Protocol Analyzer**
   - Provides deep packet inspection and traffic analysis.
   - Detects anomalies in connection patterns and protocols.
   - Aims to uncover advanced persistent threats and malware C2 channels.

---

## Development Goals

This project serves as both a practical security toolkit and a learning journey into:

- Security monitoring and detection tools.
- Windows and macOS system integration.
- Efficient and lightweight Go-based agent design.
- Building Kubernetes-ready applications.
- Kafka-powered real-time data pipelines.
- Threat simulation and response systems.
- Advanced network analysis and protocol inspection.

---

### Future Improvements

- Add gRPC support to the System Monitoring Agent for metric exports.
- Extend Threat Detection Simulation to cover advanced threat scenarios.
- Build dashboards for visualizing logs and metrics.
- Enhance the Mini XDR system with fully automated response workflows.
