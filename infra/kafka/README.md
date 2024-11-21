# Kafka Infrastructure Configuration

## Current Setup

This directory contains Kafka deployment configurations using Strimzi operator for Kubernetes. Currently using a single configuration file (`kafka-depl.yaml`) that includes:

- Kafka cluster configuration
- Metrics configuration
- User authentication
- Basic security settings

## Recommended Future Improvements

### 1. Directory Structure

```
infra/kafka/
├── base/ # Base configurations
│ ├── kafka-metrics.yaml # Shared metrics configuration
│ └── kafka-common.yaml # Common labels, annotations
├── dev/ # Development environment
│ ├── kafka-depl.yaml
│ └── kustomization.yaml
├── prod/ # Production environment
│ ├── kafka-depl.yaml
│ └── kustomization.yaml
└── local/ # Local development
├── kafka-depl.yaml
└── kustomization.yaml

```

Consider splitting into the following structure for better organization:

### 2. Environment-Specific Configurations

#### Local Development

- Single replica
- Plain listener only (no TLS)
- Ephemeral storage
- Minimal resources:
  - Memory: 512Mi-1Gi
  - CPU: 200m-500m

#### Development Environment

- Single replica
- Both plain and TLS listeners
- Ephemeral storage
- Moderate resources:
  - Memory: 1Gi-2Gi
  - CPU: 500m-1

#### Production Environment

- 3+ replicas
- TLS listener only
- Persistent storage
- Full security configuration
- Higher resources:
  - Memory: 2Gi-4Gi
  - CPU: 1-2

### 3. Security Recommendations

- Remove plain listener in production
- Implement proper SSL/TLS authentication
- Use external secrets management (e.g., HashiCorp Vault)
- Add network policies
- Implement proper backup strategies
- Configure audit logging

### 4. Monitoring Improvements

- Add Prometheus monitoring
- Configure Grafana dashboards
- Set up alerting rules
- Enhanced metrics collection

### 5. Implementation Tools

Consider using:

- Kustomize for environment management
- Helm charts for deployment
- GitOps workflows (e.g., ArgoCD, Flux)

## Migration Steps

To migrate from current single-file to multi-environment setup:

1. Create directory structure
2. Split configurations by environment
3. Implement Kustomize
4. Test in development
5. Gradually roll out to production

## Usage Examples (Future)

### Local Development

```bash
kubectl apply -k local/
```

### Development Environment

```bash
kubectl apply -k dev/
```

### Production Environment

```bash
kubectl apply -k prod/
```

## Current Limitations

- Single configuration file for all environments
- Manual environment switching
- Basic security configuration
- Limited monitoring setup

## Notes

- Current setup uses both plain and TLS listeners for development convenience
- Production deployment should be secured appropriately
- Consider implementing proper backup and disaster recovery procedures
- Monitor resource usage to adjust limits and requests

## Related Documentation

- [Strimzi Documentation](https://strimzi.io/documentation/)
- [Kafka Security](https://kafka.apache.org/documentation/#security)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
