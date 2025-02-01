# Database Schema Update Guide

This guide explains how to update database schemas and deploy the changes to your Kubernetes environment.

## General Process

1. Update centralized schema files
2. Create migration files for changes
3. Update version numbers
4. Deploy changes to Kubernetes

## PostgreSQL Schema Updates

### 1. Update Base Schema

1. Navigate to the appropriate schema file in `schemas/postgresql/log-aggregator/`
2. Update the schema version and modification date
3. Make your changes to the schema file

Example:

```sql
-- Schema Version: 1.1.0 (increment version)
-- Last Modified: 2025-01-31
-- Description: Added new column to logs table

ALTER TABLE logs ADD COLUMN severity INTEGER;
```

### 2. Create Migration File

1. Create a new migration file in `log-aggregator/migrations/`
2. Name it sequentially (e.g., `005_add_severity.sql`)
3. Include only the changes, not the full schema

Example:

```sql
-- Migration: 005_add_severity
-- Schema Version: 1.1.0
ALTER TABLE logs ADD COLUMN severity INTEGER;
```

### 3. Deploy Changes

```bash
# 1. Update the schema ConfigMap
kubectl delete configmap postgres-schemas
# Create ConfigMap from schema files (using valid key names)
cd schemas/postgresql/log-aggregator
kubectl create configmap postgres-schemas \
  --from-file=logs.sql \
  --from-file=alerts.sql

# 2. Run a new migration job
kubectl delete job postgres-migrations
kubectl apply -f infra/k8s/postgres-migrations-job.yaml
```

## MongoDB Schema Updates

### 1. Update TypeScript Schemas

1. Navigate to `schemas/mongodb/monitoring-gateway/schemas.ts`
2. Update interfaces and schema definitions
3. Increment the schema version

Example:

```typescript
export interface ITenant {
  organizationName: string;
  contactEmail: string;
  status: "active" | "inactive";
  tier: string; // New field
  createdAt: Date;
  updatedAt: Date;
}
```

### 2. Update MongoDB Validation

1. Update the schema validation in `infra/k8s/mongo-init-configmap.yaml`
2. Include new fields and constraints

Example:

```yaml
validator:
  {
    $jsonSchema:
      {
        bsonType: "object",
        required: ["organizationName", "contactEmail", "tier"],
        properties:
          {
            tier:
              { bsonType: "string", description: "Customer tier - required" },
          },
      },
  }
```

### 3. Deploy Changes

```bash
# 1. Update the MongoDB init ConfigMap
kubectl delete configmap mongo-init-scripts
kubectl apply -f infra/k8s/mongo-init-configmap.yaml

# 2. Restart MongoDB to apply new validation
kubectl rollout restart statefulset system-monitoring-mongodb
```

## Testing Schema Updates

### PostgreSQL

1. Test migrations locally:

```bash
# Apply migration to test database
psql -d testdb -f migrations/005_add_severity.sql

# Verify changes
psql -d testdb -c "\d+ logs"
```

### MongoDB

1. Test schema validation locally:

```bash
# Start MongoDB with test config
mongod --port 27018

# Apply and test validation
mongosh --port 27018 --eval "load('mongo-init-configmap.yaml')"
```

## Best Practices

1. **Version Control**

   - Always increment schema versions
   - Document changes in schema files
   - Keep migration files small and focused

2. **Backward Compatibility**

   - Make additive changes when possible
   - Provide default values for new required fields
   - Consider data migration needs

3. **Testing**

   - Test migrations on a copy of production data
   - Verify all applications work with new schema
   - Have a rollback plan ready

4. **Deployment**
   - Deploy during low-traffic periods
   - Monitor application logs during update
   - Have database backups ready

## Rollback Procedures

### PostgreSQL Rollback

1. Create a rollback migration:

```sql
-- rollback_005_add_severity.sql
ALTER TABLE logs DROP COLUMN severity;
```

2. Apply rollback:

```bash
kubectl exec -it postgres-pod -- psql -d logdb -f /migrations/rollback_005_add_severity.sql
```

### MongoDB Rollback

1. Keep previous ConfigMap version
2. Rollback deployment:

```bash
kubectl rollout undo statefulset system-monitoring-mongodb
```

## Common Issues

1. **Migration Failures**

   - Check database connectivity
   - Verify SQL syntax
   - Ensure migrations run in correct order

2. **Validation Errors**

   - Verify existing data meets new constraints
   - Check JSON Schema syntax
   - Monitor MongoDB logs for validation issues

3. **Application Impact**
   - Update application code for schema changes
   - Test all affected queries
   - Monitor application performance

## Need Help?

- Check migration logs: `kubectl logs job/postgres-migrations`
- Check MongoDB logs: `kubectl logs system-monitoring-mongodb-0`
- Review schema version history in centralized schema files
