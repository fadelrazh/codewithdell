# CodeWithDell Deployment Guide

## Overview

This guide covers deploying CodeWithDell to various environments, from local development to production.

## Prerequisites

- Docker & Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)
- Git

## Local Development

### Quick Start

1. **Clone the repository**

```bash
git clone https://github.com/your-username/codewithdell.git
cd codewithdell
```

2. **Run the development script**

```bash
./scripts/dev.sh setup
./scripts/dev.sh start
```

3. **Access the application**

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Documentation: http://localhost:8080/swagger/
- Database GUI: http://localhost:8081
- Redis GUI: http://localhost:8082

### Manual Setup

1. **Environment Configuration**

```bash
# Copy environment files
cp .env.example .env
cp frontend/.env.local.example frontend/.env.local

# Edit configuration
nano .env
nano frontend/.env.local
```

2. **Start Infrastructure**

```bash
docker-compose up -d postgres redis
```

3. **Install Dependencies**

```bash
# Backend
cd backend
go mod download
cd ..

# Frontend
cd frontend
npm install
cd ..
```

4. **Run Migrations**

```bash
cd backend
go run main.go migrate
cd ..
```

5. **Start Development Servers**

```bash
# Backend (in one terminal)
cd backend
go run main.go

# Frontend (in another terminal)
cd frontend
npm run dev
```

## Staging Environment

### Docker Deployment

1. **Build Images**

```bash
# Build all images
docker-compose -f docker-compose.staging.yml build

# Or build individually
docker build -f backend/Dockerfile.prod -t codewithdell-backend:staging ./backend
docker build -f frontend/Dockerfile.prod -t codewithdell-frontend:staging ./frontend
```

2. **Deploy**

```bash
docker-compose -f docker-compose.staging.yml up -d
```

### Environment Variables (Staging)

```bash
# Application
ENVIRONMENT=staging
PORT=8080
CORS_ORIGIN=https://staging.codewithdell.com

# Database
DB_HOST=staging-db.example.com
DB_PORT=5432
DB_NAME=codewithdell_staging
DB_USER=codewithdell_staging
DB_PASSWORD=secure_password

# Redis
REDIS_HOST=staging-redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=secure_redis_password

# JWT
JWT_SECRET=your-staging-jwt-secret
JWT_EXPIRATION=24h

# Email
EMAIL_PROVIDER=sendgrid
EMAIL_API_KEY=your_sendgrid_api_key
EMAIL_FROM=noreply@staging.codewithdell.com

# Storage
STORAGE_PROVIDER=s3
STORAGE_BUCKET=codewithdell-staging
STORAGE_REGION=us-east-1
STORAGE_ACCESS_KEY=your_aws_access_key
STORAGE_SECRET_KEY=your_aws_secret_key
```

## Production Environment

### Infrastructure Setup

#### Option 1: Cloud Provider (AWS)

1. **Create VPC and Security Groups**

```bash
# Using AWS CLI or Terraform
aws ec2 create-vpc --cidr-block 10.0.0.0/16
aws ec2 create-security-group --group-name codewithdell-prod --description "CodeWithDell Production"
```

2. **Set up RDS PostgreSQL**

```bash
aws rds create-db-instance \
  --db-instance-identifier codewithdell-prod \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --master-username codewithdell \
  --master-user-password secure_password \
  --allocated-storage 20
```

3. **Set up ElastiCache Redis**

```bash
aws elasticache create-cache-cluster \
  --cache-cluster-id codewithdell-redis \
  --engine redis \
  --cache-node-type cache.t3.micro \
  --num-cache-nodes 1
```

4. **Set up S3 for File Storage**

```bash
aws s3 mb s3://codewithdell-prod
aws s3api put-bucket-versioning --bucket codewithdell-prod --versioning-configuration Status=Enabled
```

#### Option 2: Kubernetes

1. **Create Namespace**

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: codewithdell
```

2. **Deploy PostgreSQL**

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: codewithdell
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:15-alpine
          env:
            - name: POSTGRES_DB
              value: codewithdell
            - name: POSTGRES_USER
              value: codewithdell
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: password
          ports:
            - containerPort: 5432
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
```

3. **Deploy Redis**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: codewithdell
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          ports:
            - containerPort: 6379
          volumeMounts:
            - name: redis-storage
              mountPath: /data
      volumes:
        - name: redis-storage
          persistentVolumeClaim:
            claimName: redis-pvc
```

### Application Deployment

#### Docker Compose (Production)

1. **Create Production Compose File**

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./docker/nginx/nginx.prod.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - frontend
      - backend

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.prod
    environment:
      - NODE_ENV=production
      - NEXT_PUBLIC_API_URL=https://api.codewithdell.com
    restart: unless-stopped

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.prod
    environment:
      - ENVIRONMENT=production
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - JWT_SECRET=${JWT_SECRET}
    restart: unless-stopped
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

2. **Deploy**

```bash
# Set environment variables
export $(cat .env.prod | xargs)

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

#### Kubernetes Deployment

1. **Create ConfigMap**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: codewithdell-config
  namespace: codewithdell
data:
  ENVIRONMENT: "production"
  DB_HOST: "postgres"
  DB_PORT: "5432"
  DB_NAME: "codewithdell"
  REDIS_HOST: "redis"
  REDIS_PORT: "6379"
```

2. **Create Secrets**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: codewithdell-secrets
  namespace: codewithdell
type: Opaque
data:
  DB_PASSWORD: <base64-encoded-password>
  REDIS_PASSWORD: <base64-encoded-redis-password>
  JWT_SECRET: <base64-encoded-jwt-secret>
  EMAIL_API_KEY: <base64-encoded-email-api-key>
  AWS_ACCESS_KEY: <base64-encoded-aws-access-key>
  AWS_SECRET_KEY: <base64-encoded-aws-secret-key>
```

3. **Deploy Backend**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: codewithdell
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
        - name: backend
          image: codewithdell/backend:latest
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef:
                name: codewithdell-config
            - secretRef:
                name: codewithdell-secrets
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

4. **Deploy Frontend**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  namespace: codewithdell
spec:
  replicas: 3
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
        - name: frontend
          image: codewithdell/frontend:latest
          ports:
            - containerPort: 3000
          env:
            - name: NODE_ENV
              value: "production"
            - name: NEXT_PUBLIC_API_URL
              value: "https://api.codewithdell.com"
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
```

5. **Create Services**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-service
  namespace: codewithdell
spec:
  selector:
    app: backend
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: frontend-service
  namespace: codewithdell
spec:
  selector:
    app: frontend
  ports:
    - protocol: TCP
      port: 80
      targetPort: 3000
  type: ClusterIP
```

6. **Create Ingress**

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: codewithdell-ingress
  namespace: codewithdell
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
    - hosts:
        - codewithdell.com
        - api.codewithdell.com
      secretName: codewithdell-tls
  rules:
    - host: codewithdell.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
    - host: api.codewithdell.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: backend-service
                port:
                  number: 80
```

### SSL/TLS Configuration

#### Let's Encrypt with Certbot

1. **Install Certbot**

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install certbot

# CentOS/RHEL
sudo yum install certbot
```

2. **Obtain Certificate**

```bash
sudo certbot certonly --standalone -d codewithdell.com -d api.codewithdell.com
```

3. **Auto-renewal**

```bash
# Add to crontab
sudo crontab -e

# Add this line
0 12 * * * /usr/bin/certbot renew --quiet
```

#### Nginx SSL Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name codewithdell.com;

    ssl_certificate /etc/letsencrypt/live/codewithdell.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/codewithdell.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    location / {
        proxy_pass http://frontend:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 443 ssl http2;
    server_name api.codewithdell.com;

    ssl_certificate /etc/letsencrypt/live/api.codewithdell.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.codewithdell.com/privkey.pem;

    location / {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name codewithdell.com api.codewithdell.com;
    return 301 https://$server_name$request_uri;
}
```

## Monitoring and Logging

### Prometheus Monitoring

1. **Deploy Prometheus**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
        - name: prometheus
          image: prom/prometheus:latest
          ports:
            - containerPort: 9090
          volumeMounts:
            - name: prometheus-config
              mountPath: /etc/prometheus
      volumes:
        - name: prometheus-config
          configMap:
            name: prometheus-config
```

2. **Prometheus Configuration**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  namespace: monitoring
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
    scrape_configs:
    - job_name: 'codewithdell-backend'
      static_configs:
      - targets: ['backend:8080']
      metrics_path: /metrics
```

### Grafana Dashboard

1. **Deploy Grafana**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
        - name: grafana
          image: grafana/grafana:latest
          ports:
            - containerPort: 3000
          env:
            - name: GF_SECURITY_ADMIN_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: grafana-secret
                  key: admin-password
```

### Log Aggregation

#### ELK Stack (Elasticsearch, Logstash, Kibana)

1. **Deploy Elasticsearch**

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: elasticsearch
  namespace: logging
spec:
  serviceName: elasticsearch
  replicas: 3
  selector:
    matchLabels:
      app: elasticsearch
  template:
    metadata:
      labels:
        app: elasticsearch
    spec:
      containers:
        - name: elasticsearch
          image: docker.elastic.co/elasticsearch/elasticsearch:8.0.0
          env:
            - name: discovery.type
              value: single-node
          ports:
            - containerPort: 9200
            - containerPort: 9300
```

2. **Deploy Logstash**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: logstash
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: logstash
  template:
    metadata:
      labels:
        app: logstash
    spec:
      containers:
        - name: logstash
          image: docker.elastic.co/logstash/logstash:8.0.0
          ports:
            - containerPort: 5044
          volumeMounts:
            - name: logstash-config
              mountPath: /usr/share/logstash/pipeline
      volumes:
        - name: logstash-config
          configMap:
            name: logstash-config
```

3. **Deploy Kibana**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kibana
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kibana
  template:
    metadata:
      labels:
        app: kibana
    spec:
      containers:
        - name: kibana
          image: docker.elastic.co/kibana/kibana:8.0.0
          ports:
            - containerPort: 5601
          env:
            - name: ELASTICSEARCH_HOSTS
              value: http://elasticsearch:9200
```

## Backup and Recovery

### Database Backup

1. **Automated Backup Script**

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="codewithdell"
DB_USER="codewithdell"

# Create backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/backup_$DATE.sql

# Upload to S3
aws s3 cp $BACKUP_DIR/backup_$DATE.sql.gz s3://codewithdell-backups/

# Clean old backups (keep last 30 days)
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
```

2. **Cron Job**

```bash
# Add to crontab
0 2 * * * /path/to/backup.sh
```

### File Storage Backup

```bash
#!/bin/bash
# storage-backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/storage"

# Sync with S3
aws s3 sync /uploads $BACKUP_DIR/uploads_$DATE
aws s3 sync $BACKUP_DIR/uploads_$DATE s3://codewithdell-storage-backups/

# Clean old backups
find $BACKUP_DIR -name "uploads_*" -mtime +7 -exec rm -rf {} \;
```

## Security Hardening

### Network Security

1. **Firewall Rules**

```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# iptables
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -j DROP
```

2. **Security Groups (AWS)**

```bash
# Backend security group
aws ec2 create-security-group \
  --group-name codewithdell-backend \
  --description "CodeWithDell Backend Security Group"

aws ec2 authorize-security-group-ingress \
  --group-name codewithdell-backend \
  --protocol tcp \
  --port 8080 \
  --source-group codewithdell-frontend
```

### Application Security

1. **Secrets Management**

```bash
# Use HashiCorp Vault
vault kv put secret/codewithdell \
  db_password=secure_password \
  jwt_secret=secure_jwt_secret \
  redis_password=secure_redis_password
```

2. **Container Security**

```bash
# Scan images for vulnerabilities
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image codewithdell/backend:latest

# Run containers as non-root
docker run --user 1001:1001 codewithdell/backend:latest
```

## Performance Optimization

### Database Optimization

1. **Connection Pooling**

```go
// In Go application
db.DB().SetMaxOpenConns(25)
db.DB().SetMaxIdleConns(5)
db.DB().SetConnMaxLifetime(5 * time.Minute)
```

2. **Query Optimization**

```sql
-- Add indexes for common queries
CREATE INDEX idx_posts_status_published_at ON posts(status, published_at);
CREATE INDEX idx_projects_difficulty ON projects(difficulty);
CREATE INDEX idx_comments_post_id ON comments(post_id);
```

### Caching Strategy

1. **Redis Configuration**

```bash
# redis.conf
maxmemory 256mb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
```

2. **Application Caching**

```go
// Cache frequently accessed data
func (s *PostService) GetPopularPosts() ([]Post, error) {
    cacheKey := "popular_posts"

    // Try cache first
    if cached, err := s.cache.Get(cacheKey); err == nil {
        return cached, nil
    }

    // Fetch from database
    posts, err := s.repo.GetPopularPosts()
    if err != nil {
        return nil, err
    }

    // Cache for 1 hour
    s.cache.Set(cacheKey, posts, time.Hour)

    return posts, nil
}
```

## Troubleshooting

### Common Issues

1. **Database Connection Issues**

```bash
# Check database connectivity
docker exec -it postgres psql -U codewithdell -d codewithdell -c "SELECT 1;"

# Check logs
docker logs postgres
```

2. **Redis Connection Issues**

```bash
# Test Redis connection
docker exec -it redis redis-cli ping

# Check Redis logs
docker logs redis
```

3. **Application Issues**

```bash
# Check application logs
docker logs backend
docker logs frontend

# Check health endpoints
curl http://localhost:8080/health
curl http://localhost:3000/api/health
```

### Performance Issues

1. **High CPU Usage**

```bash
# Check resource usage
docker stats

# Profile Go application
go tool pprof http://localhost:8080/debug/pprof/profile
```

2. **High Memory Usage**

```bash
# Check memory usage
free -h
docker stats --no-stream

# Analyze memory leaks
go tool pprof http://localhost:8080/debug/pprof/heap
```

### Monitoring Alerts

1. **Prometheus Alert Rules**

```yaml
groups:
  - name: codewithdell
    rules:
      - alert: HighCPUUsage
        expr: container_cpu_usage_seconds_total{container="backend"} > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage on backend"
          description: "Backend container is using more than 80% CPU"

      - alert: HighMemoryUsage
        expr: container_memory_usage_bytes{container="backend"} > 500000000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage on backend"
          description: "Backend container is using more than 500MB memory"
```

## Maintenance

### Regular Maintenance Tasks

1. **Database Maintenance**

```sql
-- Vacuum database
VACUUM ANALYZE;

-- Update statistics
ANALYZE;

-- Check for bloat
SELECT schemaname, tablename, n_dead_tup, n_live_tup,
       round(n_dead_tup * 100.0 / nullif(n_live_tup, 0), 2) AS dead_percentage
FROM pg_stat_user_tables
WHERE n_dead_tup > 0
ORDER BY dead_percentage DESC;
```

2. **Log Rotation**

```bash
# Configure logrotate
cat > /etc/logrotate.d/codewithdell << EOF
/var/log/codewithdell/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 codewithdell codewithdell
    postrotate
        systemctl reload nginx
    endscript
}
EOF
```

3. **Security Updates**

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Update Docker images
docker-compose pull
docker-compose up -d

# Update application dependencies
cd backend && go mod tidy && cd ..
cd frontend && npm audit fix && cd ..
```

This deployment guide provides a comprehensive approach to deploying CodeWithDell in various environments. Choose the deployment method that best fits your infrastructure and requirements.
