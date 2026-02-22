# ShopMicro Production Platform

## 1. Problem Statement and Architecture Summary

ShopMicro is a microservices e-commerce platform consisting of:

- **Frontend**: React.js web application for customer interactions
- **Backend**: Node.js/Express API for product and user management
- **ML Service**: Python/Flask recommendation engine
- **Database**: PostgreSQL for persistent data, Redis for caching
- **Observability**: Grafana stack (Mimir, Loki, Tempo) for monitoring

The platform implements a complete DevOps/Platform Engineering toolchain for production deployment on Kubernetes with comprehensive observability.

## 2. High-Level Architecture Diagram

```
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend API   │
│   (React)       │◄──►│   (Node.js)     │
│   Port: 8080    │    │   Port: 3001    │
└─────────────────┘    └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│   ML Service    │    │   PostgreSQL    │
│   (Python)      │    │   Port: 5432    │
│   Port: 3002    │    └─────────────────┘
└─────────────────┘             ▲
         │                     │
         ▼                     │
┌─────────────────┐             │
│     Redis       │◄────────────┘
│   Port: 6379    │
└─────────────────┘

Observability Stack:
- Grafana (Port 3000): Dashboards and visualization
- Mimir (Port 9009): Metrics storage
- Loki (Port 3100): Log aggregation
- Tempo (Port 3200): Distributed tracing
- Alloy: Telemetry collection
```

## 3. Prerequisites and Tooling Versions

### Required Tools

- Docker 24+
- Docker Compose 2.0+
- kubectl 1.28+
- Terraform 1.5+
- Ansible 2.15+
- Go 1.21+ (for CLI tools)
- Node.js 20+
- Python 3.12+

### Infrastructure Requirements

- Kubernetes cluster (EKS, GKE, or local)
- AWS account (for cloud deployment)
- Docker registry access

## 4. Exact Deploy Commands

### Local Development

```bash
# Start all services locally
docker compose up --build -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

### Production Deployment

```bash
# Install KEDA for event-driven autoscaling
kubectl apply -f https://github.com/kedacore/keda/releases/download/v2.13.0/keda-2.13.0.yaml

# Deploy infrastructure
cd infrastructure/terraform
terraform init
terraform plan
terraform apply

# Configure hosts
cd ../ansible
ansible-playbook -i inventory.ini playbook.yml

# Deploy to Kubernetes
kubectl apply -f k8s/

# Verify deployment
kubectl get pods
kubectl get services
```

## 5. Exact Test/Verification Commands

### Health Checks

```bash
# Run health check CLI
cd tools
go run health-check.go

# Manual health checks
curl http://localhost:3001/health
curl http://localhost:8080
curl http://localhost:3002/health
```

### Test Suite

```bash
# Backend tests
cd backend && npm test

# ML service tests
cd ../ml-service && python -m pytest

# Integration tests
docker compose exec backend npm run test:integration
```

## 6. Observability Usage Guide

### Accessing Dashboards

```bash
# Port forward Grafana
kubectl port-forward svc/grafana 3000:3000

# Open browser to http://localhost:3000
# Username: admin, Password: admin
```

### Viewing Metrics/Logs/Traces

- **Metrics**: Grafana dashboards (Platform Overview, Backend Health, ML Service)
- **Logs**: Loki in Grafana Explore, filter by service
- **Traces**: Tempo in Grafana, view distributed traces

### SLIs and SLOs

- **Availability SLI**: 99.9% uptime for backend service (SLO: 99.9% over 30 days)
- **Latency SLI**: 95% of requests < 500ms (SLO: 95% over 30 days)
- **Error Rate SLI**: < 1% 5xx errors (SLO: < 1% over 30 days)

## 7. Rollback Procedure

```bash
# Check rollout history
kubectl rollout history deployment/backend

# Rollback to previous version
kubectl rollout undo deployment/backend

# Verify rollback
kubectl rollout status deployment/backend

# Check application health
curl http://<ingress-url>/health
```

## 8. Security Controls Implemented

### Network Security

- Least privilege network policies
- No public SSH access
- Encrypted secrets in Kubernetes

### Application Security

- Non-root container execution
- Read-only root filesystems where possible
- Security contexts applied

### Infrastructure Security

- IAM roles with minimal permissions
- VPC isolation
- Encrypted data at rest/transit

## 9. Backup/Restore Procedure

See `runbooks/backup-restore.md` for detailed procedures.

### Quick Commands

```bash
# PostgreSQL backup
kubectl exec postgres-pod -- pg_dump -U postgres shopmicro > backup.sql

# Redis backup
kubectl exec redis-pod -- redis-cli SAVE
kubectl cp redis-pod:/data/dump.rdb ./redis-backup.rdb
```

## 10. Known Limitations and Next Improvements

### Current Limitations

- Single region deployment
- No automated scaling policies
- Basic authentication only
- Manual backup procedures

### Future Improvements

- Multi-region deployment with failover
- Advanced HPA with custom metrics
- OAuth2 authentication
- Automated backup to S3
- Service mesh (Istio) integration
- GitOps with ArgoCD
