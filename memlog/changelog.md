# Changelog

All notable changes to the Log Aggregator project will be documented in this file.

## [Unreleased]

### Added

- Initial project structure with domain models
- Kafka consumer implementation for log ingestion
- Basic log service with CRUD operations
- Process tracking and storage capabilities
- Alert domain model with severity levels
- PostgreSQL integration for log storage

### Work in Progress

- Alert system implementation
- REST API development
- Performance optimizations
- Monitoring and metrics setup

## [0.1.0] - 2024-01-17

### Added

- Basic project scaffolding
- Domain models for logs, processes, and alerts
- Kafka consumer setup with message processing
- PostgreSQL repository implementations
- Initial service layer implementation

### Technical Details

- Implemented log ingestion pipeline
- Added process tracking capabilities
- Set up basic error handling
- Created database schema for logs and processes

### Infrastructure

- Docker configuration for local development
- Kubernetes deployment manifests
- PostgreSQL database setup
- Kafka integration for message streaming

## Migration Notes

- Initial database schema in migrations/001_initial_schema.sql
- Kafka topics need to be created before running the service
- PostgreSQL connection requires proper credentials in config

## Dependencies

- Go 1.x
- Kafka
- PostgreSQL
- Docker/Kubernetes for deployment
