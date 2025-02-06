# PostgreSQL Schema Management

This directory contains the PostgreSQL schemas for the log aggregator service.

## Initial Schema

The initial schema is embedded in the Kubernetes ConfigMap and consists of:

- `001_initial_schema.sql`: Core logging tables and indexes
- `002_alerts_schema.sql`: Alert management tables and indexes

These are automatically applied when the system is first deployed.

## Adding Schema Updates

To add new schema changes:

1. Create a new migration file in the `log-aggregator/migrations/` directory with a sequential number:

```sql
-- Example: 005_add_new_feature.sql
ALTER TABLE logs ADD COLUMN new_column TEXT;
```

2. Update the ConfigMap with the new migration:

```bash
# Add your new migration to the ConfigMap
kubectl delete configmap postgres-migrations
kubectl create configmap postgres-migrations --from-file=001_initial_schema.sql --from-file=002_alerts_schema.sql --from-file=005_add_new_feature.sql

# Run the migration job
kubectl delete job postgres-migrations
kubectl apply -f ../../../infra/k8s/postgres-migrations-job.yaml
```

The migration job will:

- Wait for PostgreSQL to be ready
- Apply migrations in numerical order
- Stop and report errors if any migration fails

## Best Practices

1. Always use sequential numbering for migrations
2. Include both "up" and "down" migrations in comments
3. Test migrations on a development database first
4. Back up production database before applying migrations
5. Document schema version changes

## Troubleshooting

If migrations fail:

1. Check the migration job logs:

```bash
kubectl logs job/postgres-migrations
```

2. Verify PostgreSQL is running:

```bash
kubectl exec -it postgres-srv -- psql -U postgres -d logdb -c "\dt"
```

3. Check ConfigMap content:

```bash
kubectl get configmap postgres-migrations -o yaml
```
