# Docker Deployment

## Overview

gSpend uses Docker Compose for containerized deployment with the following services:
- Frontend (Vue.js + Nginx)
- Auth Service (Go)
- Financial Service (Go)
- MongoDB
- Redis
- Nginx (API Gateway)

---

## Quick Start

```bash
# Development environment
make docker-up

# Demo with sample data
make demo-start

# Production build
docker compose -f docker-compose.demo.yml up --build -d
```

---

## Docker Compose Files

| File | Purpose |
|------|---------|
| `devops/docker/docker-compose.yml` | Development environment |
| `docker-compose.demo.yml` | Demo with sample data |
| `docker-compose.test.yml` | Integration testing |

---

## Service Configuration

### Frontend
```yaml
frontend:
  build:
    context: ./apps/frontend
    dockerfile: Dockerfile
  environment:
    - VITE_AUTH_SERVICE_URL=http://localhost/api/v1/auth
    - VITE_FINANCIAL_SERVICE_URL=http://localhost/api/v1
```

### Auth Service
```yaml
auth-service:
  build:
    context: ./apps/auth-service
  ports:
    - "9091:9091"  # gRPC
  environment:
    - MONGODB_URI=mongodb://mongodb:27017
    - MONGODB_DATABASE=gspend
    - REDIS_HOST=redis
    - JWT_SECRET=${JWT_SECRET:-gspend-secret-key}
    - PORT=8081
    - GRPC_PORT=9091
```

### Financial Service
```yaml
financial-service:
  build:
    context: ./apps/financial-service
  environment:
    - MONGODB_URI=mongodb://mongodb:27017
    - MONGODB_DATABASE=gspend
    - JWT_SECRET=${JWT_SECRET:-gspend-secret-key}
    - AUTH_SERVICE_GRPC_ADDR=auth-service:9091
```

---

## Network Architecture

```
┌─────────────────────────────────────────────────┐
│                  gspend-net                     │
│                                                 │
│  ┌─────────┐    ┌──────────┐    ┌───────────┐  │
│  │ nginx   │───►│   auth   │◄──►│ financial │  │
│  │  :80    │    │  :8081   │    │   :8082   │  │
│  └────▲────┘    └────┬─────┘    └─────┬─────┘  │
│       │              │                │        │
│  ┌────┴────┐    ┌────▼────────────────▼────┐   │
│  │frontend │    │        mongodb:27017     │   │
│  └─────────┘    │        redis:6379        │   │
│                 └──────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

---

## Commands

```bash
# Start services
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs
make docker-logs-auth
make docker-logs-finance
make docker-logs-frontend

# Restart
make docker-restart

# Check status
make docker-status

# Clean up
make docker-clean
```

---

## Health Checks

All services include health checks:

```bash
make health-check
```

Or manually:
```bash
curl http://localhost/api/v1/auth/health
curl http://localhost/api/v1/health
```

---

## Volumes

| Volume | Purpose |
|--------|---------|
| `mongodb-data` | MongoDB data persistence |
| `redis-data` | Redis data persistence |

---

## Environment Variables

Create a `.env` file in the project root:

```env
JWT_SECRET=your-production-secret-key
MONGODB_DATABASE=gspend
```
