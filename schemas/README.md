# Centralized Database Schemas

This directory contains the centralized database schemas for all microservices in the Security-Monitoring-Suite.

## Directory Structure

```
schemas/
├── postgresql/           # PostgreSQL schemas
│   ├── log-aggregator/  # Log Aggregator service schemas
│   └── migrations/      # Shared migration templates
├── mongodb/             # MongoDB schemas
│   └── monitoring-gateway/ # System Monitoring Gateway schemas
└── README.md           # This file
```

## Usage

### PostgreSQL Schemas

- Located in `postgresql/` directory
- Written in standard SQL format
- Include indexes and constraints
- Referenced by service-specific migrations

### MongoDB Schemas

- Located in `mongodb/` directory
- Written in JSON Schema format
- Include TypeScript interfaces
- Referenced by service models

## Schema Versioning

Each schema file includes:

- Version number
- Creation date
- Last modified date
- Change history
- Dependencies on other schemas

## Best Practices

1. Always reference these centralized schemas in service implementations
2. Document any changes in the schema version history
3. Update dependent services when making schema changes
4. Include validation rules and constraints
5. Maintain backward compatibility when possible
