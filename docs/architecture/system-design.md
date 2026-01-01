# System Design

## Overview

gSpend is a family financial management application built with a microservices architecture. It helps families track income, plan budgets, monitor spending, and analyze transactions.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│                   (Mobile / Desktop Browser)                    │
└───────────────────────────┬─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                      Vue.js SPA Frontend                        │
│                    (TypeScript, Pinia, Vite)                    │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTP/REST
┌───────────────────────────▼─────────────────────────────────────┐
│                    Nginx (API Gateway)                          │
│            Route: /api/v1/auth → Auth Service                   │
│            Route: /api/v1/* → Financial Service                 │
└───────────┬───────────────────────────────────┬─────────────────┘
            │                                   │
┌───────────▼───────────┐         ┌─────────────▼─────────────────┐
│    Auth Service       │◄──gRPC──►     Financial Service         │
│    (Go, Echo)         │         │      (Go, Echo)               │
│    Port 8081 (HTTP)   │         │      Port 8082 (HTTP)         │
│    Port 9091 (gRPC)   │         │      Port 9092 (gRPC)         │
└───────────┬───────────┘         └─────────────┬─────────────────┘
            │                                   │
            └─────────────┬─────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                       Data Layer                                │
│    ┌─────────────────┐              ┌─────────────────┐         │
│    │    MongoDB 7    │              │    Redis 7      │         │
│    │  (Primary DB)   │              │   (Cache)       │         │
│    └─────────────────┘              └─────────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Services

### Auth Service
- **Purpose**: Authentication and user management
- **Technology**: Go, Echo framework
- **Responsibilities**:
  - User registration and login
  - JWT token generation/validation
  - Password management
  - Profile management
- **Ports**: HTTP 8081, gRPC 9091

### Financial Service
- **Purpose**: All financial operations
- **Technology**: Go, Echo framework
- **Responsibilities**:
  - Income management
  - Budget planning
  - Transaction recording
  - Category management
  - Dashboard and reports
- **Ports**: HTTP 8082, gRPC 9092

---

## Communication

| Type | Protocol | Use Case |
|------|----------|----------|
| Frontend ↔ Backend | HTTP/REST | Client API requests |
| Service ↔ Service | gRPC | Internal communication |

---

## Technology Stack

### Frontend
- Vue.js 3 (Composition API)
- TypeScript
- Pinia (State Management)
- Vue Router 4
- Tailwind CSS
- Chart.js
- Vite

### Backend
- Go 1.21+
- Echo v4 (HTTP)
- gRPC + Protocol Buffers
- JWT Authentication
- MongoDB Driver
- Zap Logger

### Infrastructure
- Docker & Docker Compose
- Nginx (Reverse Proxy)
- MongoDB 7
- Redis 7

---

## Security

1. **Authentication**: JWT tokens with refresh mechanism
2. **Password**: Bcrypt hashing (cost 12)
3. **Authorization**: User ID from JWT claims
4. **Token TTL**: Access 15min, Refresh 7 days
5. **Communication**: HTTPS in production

---

## Scalability

- Stateless services (horizontal scaling)
- Redis for session/cache data
- MongoDB replica set support
- Kubernetes-ready containers
