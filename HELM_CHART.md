# S3 Smart Browser Helm Chart

This Helm chart deploys S3 Smart Browser on a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- S3-compatible storage access

## Installation

### Add Helm Repository

```bash
helm repo add s3-smart-browser https://blackojack.github.io/s3-smart-browser
helm repo update
```

### Install Chart

```bash
helm install s3-smart-browser s3-smart-browser/s3-smart-browser \
  --set config.aws.bucket=your-bucket-name \
  --set config.aws.region=us-east-1 \
  --set secrets.awsAccessKeyId=your_access_key \
  --set secrets.awsSecretAccessKey=your_secret_key
```

### Install with Custom Values

```bash
helm install s3-smart-browser s3-smart-browser/s3-smart-browser -f helm-values-example.yaml
```

## Configuration

### Required Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `config.aws.bucket` | S3 bucket name | `""` |
| `secrets.awsAccessKeyId` | AWS Access Key ID | `""` |
| `secrets.awsSecretAccessKey` | AWS Secret Access Key | `""` |

### Optional Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `config.aws.region` | AWS region | `us-east-1` |
| `config.aws.endpoint` | Custom S3 endpoint | `""` |
| `config.app.baseDirectory` | Base directory in bucket | `""` |
| `replicaCount` | Number of replicas | `2` |
| `image.repository` | Image repository | `ghcr.io/blackojack/s3-smart-browser` |
| `image.tag` | Image tag | `latest` |
| `service.type` | Service type | `ClusterIP` |
| `ingress.enabled` | Enable ingress | `false` |

## Examples

### Basic Installation

```yaml
config:
  aws:
    bucket: "my-s3-bucket"
    region: "us-east-1"
secrets:
  awsAccessKeyId: "AKIAIOSFODNN7EXAMPLE"
  awsSecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
```

### With Ingress

```yaml
ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: s3-browser.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: s3-browser-tls
      hosts:
        - s3-browser.example.com
```

### With MinIO

```yaml
config:
  aws:
    bucket: "test-bucket"
    region: "us-east-1"
    endpoint: "http://minio:9000"
secrets:
  awsAccessKeyId: "minioadmin"
  awsSecretAccessKey: "minioadmin123"
```

## Upgrading

```bash
helm upgrade s3-smart-browser s3-smart-browser/s3-smart-browser
```

## Uninstalling

```bash
helm uninstall s3-smart-browser
```

## Troubleshooting

### Common Issues

1. **Pod not starting**: Check AWS credentials and bucket permissions
2. **Cannot access files**: Verify S3 bucket policy and IAM permissions
3. **Ingress not working**: Check ingress controller and DNS configuration

### Debug Commands

```bash
# Check pod status
kubectl get pods -l app=s3-smart-browser

# Check logs
kubectl logs -l app=s3-smart-browser

# Check service
kubectl get svc s3-smart-browser

# Check ingress
kubectl get ingress s3-smart-browser
```

## Security

### Best Practices

- Use IAM roles instead of access keys when possible
- Implement proper RBAC policies
- Use Kubernetes secrets for sensitive data
- Enable network policies
- Use non-root containers

### Security Context

The chart includes security contexts:

```yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 100
  runAsGroup: 101
  fsGroup: 101

securityContext:
  allowPrivilegeEscalation: false
  runAsNonRoot: true
  runAsUser: 100
  capabilities:
    drop:
      - ALL
```

## Monitoring

### Health Checks

The application includes health checks:

```yaml
livenessProbe:
  httpGet:
    path: /
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### Metrics

Integration with monitoring systems:

- Prometheus metrics
- Grafana dashboards
- AlertManager rules

## Support

For issues with the Helm chart:

- GitHub Issues: [Chart Issues](https://github.com/BlackoJack/s3-smart-browser/issues)
- Email: devopsnskru@gmail.com
- Telegram: @devopsnsk

---

⭐ **If you find this chart useful, please give it a star!** ⭐
