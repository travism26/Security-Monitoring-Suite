# Schema Integration Guide

This guide explains how to integrate the centralized schemas into your microservices.

## PostgreSQL Services

### 1. Reference Schema in Migrations

Instead of maintaining schemas in your service's migration files, reference the centralized schema:

```sql
-- In your service's migration file
\i ../../schemas/postgresql/log-aggregator/logs.sql
```

### 2. Version Control

When making schema changes:

1. Update the centralized schema file
2. Create a new migration in your service that applies only the changes
3. Document the schema version being used in your service

Example migration for schema changes:

```sql
-- 005_update_logs_table.sql
ALTER TABLE logs ADD COLUMN severity INTEGER;
-- Schema version: 1.1.0
```

## MongoDB Services

### 1. Import Schemas

```typescript
// In your service's model file
import {
  mongooseTenantConfig,
  mongooseApiKeyConfig,
  tenantIndexes,
  apiKeyIndexes,
} from "../../../schemas/mongodb/monitoring-gateway/schemas";
import mongoose from "mongoose";

// Use the centralized schema configuration
const tenantSchema = new mongoose.Schema(mongooseTenantConfig, {
  timestamps: true,
  toJSON: {
    transform(doc, ret) {
      ret.id = ret._id;
      delete ret._id;
      delete ret.__v;
    },
  },
});

// Apply indexes
tenantIndexes.forEach((index) =>
  tenantSchema.index(index.key, { unique: index.unique })
);
```

### 2. Type Safety

```typescript
import {
  ITenant,
  IApiKey,
} from "../../../schemas/mongodb/monitoring-gateway/schemas";

// Use interfaces in your service
interface TenantDoc extends mongoose.Document, ITenant {}
interface ApiKeyDoc extends mongoose.Document, IApiKey {}
```

## Setup

1. Install dependencies:

```bash
cd schemas
npm install
```

2. Add schemas directory to your service's package.json:

```json
{
  "dependencies": {
    "security-monitoring-schemas": "file:../schemas"
  }
}
```

## Best Practices

1. **Schema Updates**

   - Always update the centralized schema first
   - Create service-specific migrations for changes
   - Update schema version number
   - Document breaking changes

2. **Version Control**

   - Track schema versions in your service
   - Plan for backward compatibility
   - Include upgrade guides for breaking changes

3. **Testing**

   - Test migrations against the centralized schema
   - Validate data against schema constraints
   - Include schema validation in integration tests

4. **Documentation**
   - Document schema dependencies
   - Keep change history updated
   - Include examples for common operations

## Troubleshooting

1. **Schema Import Issues**

   - Verify relative paths to schema files
   - Check package.json dependencies
   - Ensure all required dependencies are installed

2. **Migration Failures**

   - Verify schema version compatibility
   - Check for missing dependencies
   - Validate SQL syntax and MongoDB schema structure

3. **Type Errors**
   - Update TypeScript interfaces
   - Check for breaking changes in schema
   - Verify proper type exports/imports
