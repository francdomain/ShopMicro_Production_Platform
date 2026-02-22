# ShopMicro Local Setup Guide

A complete step-by-step guide to run the ShopMicro application locally using Docker Compose.

## Prerequisites

### Required Software

```bash
# Check versions
docker --version        # Docker 24+
docker compose version  # Docker Compose 2.0+
git --version          # Any recent version
```

### Installation (if needed)

- **Docker Desktop**: [Download here](https://www.docker.com/products/docker-desktop)
- **Git**: [Download here](https://git-scm.com/)

### Verify Installation

```bash
docker run hello-world
docker compose version
```

## Step 1: Clone and Navigate to Project

```bash
# Clone the repository (if not already done)
git clone <repository-url>
cd capstone-project

# Verify you're in the correct directory
pwd  # Should show: .../capstone-project
ls   # Should show: docker-compose.yml, backend/, frontend/, ml-service/, etc.
```

## Step 2: Prepare the Environment

```bash
# Navigate to project root
cd /Users/francdomain/Desktop/Dev-foundry/k8s/capstone-project

# Verify docker-compose.yml exists
ls -la docker-compose.yml

# (Optional) Clean up any previous containers
docker compose down --volumes
```

## Step 3: Start All Services

```bash
# Build and start all services in the background
docker compose up --build -d

# Expected output:
# [+] Running 11/11
#  ✔ Network shopmicro-network Created
#  ✔ Container shopmicro-postgres Started
#  ✔ Container shopmicro-redis Started
#  ✔ Container shopmicro-alloy Started
#  ✔ Container shopmicro-mimir Started
#  ✔ Container shopmicro-loki Started
#  ✔ Container shopmicro-tempo Started
#  ✔ Container shopmicro-backend Started
#  ✔ Container shopmicro-ml-service Started
#  ✔ Container shopmicro-grafana Started
#  ✔ Container shopmicro-frontend Started
```

**Explanation**: The `--build` flag rebuilds images, `-d` runs in detached mode (background)

## Step 4: Verify All Services Are Running

```bash
# Check container status (all should show "healthy" or "running")
docker compose ps

# Expected output:
# NAME                  STATUS              PORTS
# shopmicro-postgres    healthy             0.0.0.0:5432->5432/tcp
# shopmicro-redis       healthy             0.0.0.0:6379->6379/tcp
# shopmicro-backend     healthy (running)   0.0.0.0:3001->3001/tcp
# shopmicro-frontend    running             0.0.0.0:8080->3000/tcp
# shopmicro-ml-service  healthy (running)   0.0.0.0:3002->3002/tcp
# shopmicro-grafana     running             0.0.0.0:3000->3000/tcp
# ... and more
```

## Step 5: Test Application Endpoints

### Backend Service

```bash
# Test backend health
curl http://localhost:3001/health

# Expected response:
# {"status":"ok","service":"backend"}
```

### Frontend Service

```bash
# Open in browser
open http://localhost:8080
# or
curl http://localhost:8080
```

### ML Service

```bash
# Test ML service health
curl http://localhost:3002/health

# Expected response:
# {"status":"ok","service":"ml-service"}
```

### Grafana Dashboards

```bash
# Open in browser
open http://localhost:3000
# or
# Username: admin
# Password: admin
```

## Step 6: View Live Logs

```bash
# View logs from all services
docker compose logs -f

# View logs from specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f ml-service
docker compose logs -f postgres

# View last 50 lines of backend logs
docker compose logs --tail=50 backend
```

## Step 7: Access Application Components

| Component   | URL                    | Purpose                  |
| ----------- | ---------------------- | ------------------------ |
| Frontend    | http://localhost:8080  | React web application    |
| Backend API | http://localhost:3001  | REST API server          |
| ML Service  | http://localhost:3002  | Recommendation engine    |
| Grafana     | http://localhost:3000  | Dashboards (admin/admin) |
| Mimir       | http://localhost:9009  | Metrics storage          |
| Loki        | http://localhost:3100  | Log aggregation          |
| Tempo       | http://localhost:3200  | Distributed tracing      |
| Alloy UI    | http://localhost:12345 | Telemetry collector      |
| PostgreSQL  | localhost:5432         | Database                 |
| Redis       | localhost:6379         | Cache                    |

## Step 8: Interact with the Application

### Test Backend API

```bash
# Get products
curl http://localhost:3001/products

# Get ML recommendations
curl http://localhost:3002/recommendations/42
```

### Access Database

```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U postgres -d shopmicro

# Common commands in psql:
# \dt                     - List all tables
# SELECT * FROM products; - View products
# \q                      - Quit
```

### Access Redis Cache

```bash
# Connect to Redis
docker compose exec redis redis-cli

# Common commands:
# PING                    - Test connection (responds: PONG)
# KEYS *                  - List all keys
# GET <key>              - Get value
# EXIT                   - Quit
```

## Step 9: View Logs in Grafana

1. Open Grafana: http://localhost:3000
2. Login: admin / admin
3. Navigate to "Explore" → Select "Loki"
4. View logs from services
5. Select "Tempo" to view traces
6. Navigate to "Dashboards" to see pre-built dashboards

## Step 10: Run Health Check

```bash
# Build and run the health check CLI tool
cd tools
go run health-check.go

# Expected output:
# ShopMicro Health Check
# ======================
# Checking http://localhost:3001/health... OK
# Checking http://localhost:8080... OK
# Checking http://localhost:3002/health... OK
# Checking http://localhost:3000/login... OK
#
# All services are healthy!
```

## Troubleshooting

### Services Not Starting

```bash
# Check for errors
docker compose logs

# Rebuild images
docker compose down
docker compose up --build -d
```

### Port Already in Use

```bash
# Find what's using a port (e.g., 3001)
lsof -i :3001

# Kill the process
kill -9 <PID>

# Or change port in docker-compose.yml
```

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker compose ps postgres

# Check database initialization
docker compose logs postgres

# Reset database
docker compose exec postgres psql -U postgres -d shopmicro -f /docker-entrypoint-initdb.d/init.sql
```

### High Memory Usage

```bash
# Check resource usage
docker stats

# Reduce memory by stopping unused services
docker compose stop grafana loki tempo mimir alloy
```

## Step 11: Stop Services

```bash
# Stop all containers (keeps data)
docker compose stop

# Stop and remove containers (keeps volumes)
docker compose down

# Stop and remove containers and volumes (clean slate)
docker compose down --volumes
```

## Step 12: Restart Services

```bash
# Start stopped services
docker compose start

# Or start from scratch with rebuild
docker compose down --volumes
docker compose up --build -d
```

## Running Individual Services

```bash
# Start only specific services
docker compose up postgres redis backend -d

# Stop specific service
docker compose stop backend

# Restart specific service
docker compose restart backend

# View logs for specific service
docker compose logs -f backend
```

## Environment Variables

Edit `docker-compose.yml` to change:

- `POSTGRES_PASSWORD`: Database password
- `NODE_ENV`: production/development
- `LOG_LEVEL`: debug/info/warn/error

## Common Commands Reference

```bash
# List all containers
docker compose ps

# View all logs
docker compose logs

# View logs with follow (live)
docker compose logs -f

# View specific service logs
docker compose logs backend

# Execute command in container
docker compose exec backend npm list

# Stop all services
docker compose stop

# Remove all containers
docker compose down

# Remove containers and volumes
docker compose down --volumes

# Rebuild all images
docker compose build

# Rebuild specific image
docker compose build backend

# Pull latest images
docker compose pull

# Scale a service
docker compose up -d --scale backend=3
```

## Development Workflow

### Making Code Changes

1. Edit files in `backend/`, `frontend/`, or `ml-service/`
2. Changes auto-reload via volume mounts
3. Check logs for errors: `docker compose logs -f backend`

### Viewing logs in real-time

```bash
# Terminal 1: Watch backend logs
docker compose logs -f backend

# Terminal 2: Run tests or curl commands
curl http://localhost:3001/products
```

### Debugging

```bash
# Connect to backend container and inspect
docker compose exec backend bash

# Inside container:
# npm list             - check dependencies
# ps aux              - see running processes
# curl localhost:3001 - test endpoints
# exit               - quit container
```

## Next Steps

- **View Metrics**: Open http://localhost:3000 (Grafana)
- **Check Traces**: Grafana → Explore → Select Tempo
- **View Logs**: Grafana → Explore → Select Loki
- **Read Documentation**: See README.md for full project info

## Complete Local Development Checklist

- [ ] Docker and Docker Compose installed
- [ ] Project cloned and navigated to root
- [ ] `docker compose up --build -d` executed
- [ ] All containers running: `docker compose ps`
- [ ] Frontend accessible: http://localhost:8080
- [ ] Backend health check: `curl http://localhost:3001/health`
- [ ] Grafana accessible: http://localhost:3000
- [ ] Health check tool passed: `cd tools && go run health-check.go`
- [ ] Logs viewable: `docker compose logs -f`

You're all set! The ShopMicro platform is running locally.
