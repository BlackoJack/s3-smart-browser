# S3 Smart Browser Deployment Examples

## Docker Compose Examples

### Basic Setup
```yaml
version: '3.8'

services:
  s3-smart-browser:
    image: ghcr.io/blackojack/s3-smart-browser:latest
    ports:
      - "8080:8080"
    environment:
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
      - AWS_BUCKET=${AWS_BUCKET}
      - AWS_REGION=${AWS_REGION}
    restart: unless-stopped
```

### With MinIO
```yaml
version: '3.8'

services:
  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin123
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data

  s3-smart-browser:
    image: ghcr.io/blackojack/s3-smart-browser:latest
    ports:
      - "8080:8080"
    environment:
      - AWS_ACCESS_KEY_ID=minioadmin
      - AWS_SECRET_ACCESS_KEY=minioadmin123
      - AWS_BUCKET=test-bucket
      - AWS_REGION=us-east-1
      - AWS_ENDPOINT=http://minio:9000
    depends_on:
      - minio
    restart: unless-stopped

volumes:
  minio_data:
```

### Production Setup with Nginx
```yaml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - s3-smart-browser

  s3-smart-browser:
    image: ghcr.io/blackojack/s3-smart-browser:latest
    environment:
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
      - AWS_BUCKET=${AWS_BUCKET}
      - AWS_REGION=${AWS_REGION}
      - SERVER_PORT=8080
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Kubernetes Examples

### Basic Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: s3-smart-browser
spec:
  replicas: 2
  selector:
    matchLabels:
      app: s3-smart-browser
  template:
    metadata:
      labels:
        app: s3-smart-browser
    spec:
      containers:
      - name: s3-smart-browser
        image: ghcr.io/blackojack/s3-smart-browser:latest
        ports:
        - containerPort: 8080
        env:
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              name: s3-credentials
              key: access-key-id
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: s3-credentials
              key: secret-access-key
        - name: AWS_BUCKET
          value: "my-bucket"
        - name: AWS_REGION
          value: "us-east-1"
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: s3-smart-browser-service
spec:
  selector:
    app: s3-smart-browser
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: Secret
metadata:
  name: s3-credentials
type: Opaque
data:
  access-key-id: <base64-encoded-access-key>
  secret-access-key: <base64-encoded-secret-key>
```

### With Ingress
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: s3-smart-browser-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: s3-browser.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: s3-smart-browser-service
            port:
              number: 80
  tls:
  - hosts:
    - s3-browser.example.com
    secretName: s3-browser-tls
```

## Environment Variables Reference

### Required Variables
- `AWS_ACCESS_KEY_ID` - AWS Access Key ID
- `AWS_SECRET_ACCESS_KEY` - AWS Secret Access Key  
- `AWS_BUCKET` - S3 Bucket name

### Optional Variables
- `AWS_REGION` - AWS Region (default: us-east-1)
- `AWS_ENDPOINT` - Custom S3 endpoint (for MinIO, etc.)
- `SERVER_PORT` - Server port (default: 8080)
- `BASE_DIRECTORY` - Base directory in bucket

## Security Best Practices

### Using IAM Roles (AWS)
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: s3-smart-browser
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::ACCOUNT:role/s3-smart-browser-role
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: s3-smart-browser
spec:
  template:
    spec:
      serviceAccountName: s3-smart-browser
      containers:
      - name: s3-smart-browser
        image: ghcr.io/blackojack/s3-smart-browser:latest
        env:
        - name: AWS_BUCKET
          value: "my-bucket"
        - name: AWS_REGION
          value: "us-east-1"
        # No AWS credentials needed with IAM roles
```

### Using External Secrets Operator
```yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-secrets-manager
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-east-1
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: s3-credentials
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: s3-credentials
    creationPolicy: Owner
  data:
  - secretKey: access-key-id
    remoteRef:
      key: s3-smart-browser/credentials
      property: access-key-id
  - secretKey: secret-access-key
    remoteRef:
      key: s3-smart-browser/credentials
      property: secret-access-key
```
