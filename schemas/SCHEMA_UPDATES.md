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
# 1. Navigate to schema directory
cd schemas/postgresql/log-aggregator

# 2. Update both ConfigMaps with SQL content
# Delete existing ConfigMaps
kubectl delete configmap postgres-schemas
kubectl delete configmap postgres-migrations

# Create new ConfigMaps with direct SQL content
kubectl create configmap postgres-schemas --from-file=logs.sql --from-file=alerts.sql
kubectl create configmap postgres-migrations --from-file=logs.sql --from-file=alerts.sql

# 3. Run a new migration job
kubectl delete job postgres-migrations
kubectl apply -f ../../../infra/k8s/postgres-migrations-job.yaml
```

Important Note: Both ConfigMaps (postgres-schemas and postgres-migrations) need to be created with the actual SQL content directly from the files. Do not use PostgreSQL's \i command to import files, as this will cause the migrations to fail. The migration job expects to find the SQL content directly in the mounted ConfigMap files.

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

## Troubleshooting Guide

### Understanding Common Migration Failures

The most common migration failure occurs due to misconfigurations in how SQL files are mounted and accessed in the Kubernetes environment. Here's a detailed explanation of a typical failure scenario:

1. **Initial Error Scenario**

   ```
   error: /schemas/postgresql/log-aggregator/logs.sql: No such file or directory
   ```

   This error occurs when the postgres-migrations ConfigMap is configured to use PostgreSQL's \i command to import files:

   ```yaml
   # Incorrect configuration in postgres-migrations ConfigMap
   data:
     logs.sql: |
       \i /schemas/postgresql/log-aggregator/logs.sql
   ```

   The problem is two-fold:

   - The \i command looks for files in the filesystem
   - The path `/schemas/postgresql/log-aggregator/` doesn't exist in the container

2. **Solution Explanation**
   Instead of using \i commands, we should directly include the SQL content in the ConfigMap. This works because:

   - The SQL files are mounted directly in the /migrations directory
   - The migration job can execute the SQL directly without needing to import files

   The correct approach is to create the ConfigMap with the actual SQL content:

   ```bash
   kubectl create configmap postgres-migrations --from-file=logs.sql --from-file=alerts.sql
   ```

### Common Issues and Solutions

1. **ConfigMap and File Mounting Issues**

   Problem: Files not found in expected locations or incorrect paths

   ```
   error: /schemas/postgresql/log-aggregator/logs.sql: No such file or directory
   ```

   Solutions:

   - Verify ConfigMap creation:
     ```bash
     # Check ConfigMap contents
     kubectl get configmap postgres-schemas -o yaml
     kubectl get configmap postgres-migrations -o yaml
     ```
   - Create ConfigMaps with direct file content:

     ```bash
     # Navigate to schema directory
     cd schemas/postgresql/log-aggregator

     # Recreate ConfigMaps with proper content
     kubectl delete configmap postgres-schemas
     kubectl create configmap postgres-schemas --from-file=logs.sql --from-file=alerts.sql

     kubectl delete configmap postgres-migrations
     kubectl create configmap postgres-migrations --from-file=logs.sql --from-file=alerts.sql
     ```

   - Verify volume mounts in migration job:
     ```bash
     kubectl describe pod <postgres-migrations-pod-name>
     ```

2. **Database Connectivity Issues**

   Problem: Migration job can't connect to database

   ```
   error: connection refused
   ```

   Solutions:

   - Check if PostgreSQL is running:
     ```bash
     kubectl get pods | grep postgres
     kubectl logs postgres-srv
     ```
   - Verify service connectivity:
     ```bash
     kubectl describe service postgres-srv
     ```
   - Check credentials and environment variables:
     ```bash
     kubectl describe configmap postgres-config
     kubectl get secret postgres-secret -o yaml
     ```

3. **SQL Syntax and Migration Errors**

   Problem: SQL errors during migration

   ```
   ERROR: syntax error at or near...
   ```

   Solutions:

   - Validate SQL syntax locally:
     ```bash
     # Test against local PostgreSQL
     psql -d testdb -f logs.sql
     ```
   - Check migration order:
     ```bash
     # View migration job logs
     kubectl logs <postgres-migrations-pod-name>
     ```
   - Verify schema versions match:
     ```bash
     # Check current database version
     psql -d logdb -c "SELECT version FROM schema_versions;"
     ```

4. **Permission Issues**

   Problem: Insufficient privileges

   ```
   ERROR: permission denied for...
   ```

   Solutions:

   - Verify database user permissions:
     ```bash
     # Check user roles
     psql -d logdb -c "\du"
     ```
   - Grant necessary permissions:
     ```sql
     GRANT ALL PRIVILEGES ON DATABASE logdb TO postgres;
     ```
   - Check Kubernetes service account permissions:
     ```bash
     kubectl describe rolebinding postgres-role-binding
     ```

### Step-by-Step Migration Verification

1. **Pre-migration Checks**

   ```bash
   # Verify PostgreSQL is running
   kubectl get pods | grep postgres

   # Check existing schemas
   kubectl exec -it postgres-srv -- psql -d logdb -c "\dt"
   ```

2. **Apply Migrations**

   ```bash
   # Delete existing resources
   kubectl delete configmap postgres-schemas
   kubectl delete configmap postgres-migrations
   kubectl delete job postgres-migrations

   # Create new ConfigMaps
   cd schemas/postgresql/log-aggregator
   kubectl create configmap postgres-schemas --from-file=logs.sql --from-file=alerts.sql
   kubectl create configmap postgres-migrations --from-file=logs.sql --from-file=alerts.sql

   # Apply migration job
   kubectl apply -f ../../../infra/k8s/postgres-migrations-job.yaml
   ```

3. **Verify Migration Success**

   ```bash
   # Check migration job status
   kubectl get pods | grep postgres-migrations

   # View migration logs
   kubectl logs <postgres-migrations-pod-name>

   # Verify tables were created
   kubectl exec -it postgres-srv -- psql -d logdb -c "\dt"
   ```

### Best Practices for Troubleshooting

1. **Always Check Logs First**

   - Migration job logs
   - PostgreSQL server logs
   - Kubernetes events

2. **Verify Resources in Order**

   - ConfigMaps existence and content
   - Pod status and health
   - Volume mounts
   - Database connectivity

3. **Use Dry Runs When Possible**

   - Test migrations locally
   - Use kubectl --dry-run
   - Validate SQL syntax before applying

4. **Document and Version Control**
   - Keep track of applied migrations
   - Document troubleshooting steps
   - Maintain rollback procedures

### Need Help?

If issues persist:

1. Check migration logs: `kubectl logs job/postgres-migrations`
2. Check PostgreSQL logs: `kubectl logs postgres-srv`
3. Review schema version history in centralized schema files
4. Verify all ConfigMaps and secrets are properly configured

## Need Help?

- Check migration logs: `kubectl logs job/postgres-migrations`
- Check MongoDB logs: `kubectl logs system-monitoring-mongodb-0`
- Review schema version history in centralized schema files
