# Local Development Setup

## Prerequisites

- **Docker** 20.10+ with Docker Compose
- **Node.js** 18+ (for frontend development)
- **Go** 1.21+ (for backend development)
- **Git**

---

## Quick Start (Docker)

The fastest way to run gSpend locally:

```bash
# Clone the repository
git clone https://github.com/geekzy/gspend-app.git
cd gspend-app

# Start all services with Docker
make docker-up

# Or start demo with sample data
make demo-start
```

**Access the application:**
- Frontend: http://localhost
- Auth Service: http://localhost/api/v1/auth/health
- Financial Service: http://localhost/api/v1/health

**Demo credentials:**
- Email: `demo@gspend.com`
- Password: `password`

---

## Development Mode

### Frontend Development

```bash
cd apps/frontend

# Install dependencies
npm install

# Start dev server with hot reload
npm run dev
```

Frontend runs at http://localhost:5173

### Backend Development

**Auth Service:**
```bash
cd apps/auth-service

# Run the service
go run cmd/server/main.go
```

**Financial Service:**
```bash
cd apps/financial-service

# Run the service
go run cmd/server/main.go
```

---

## Environment Variables

### Auth Service
```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=gspend
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key
PORT=8081
GRPC_PORT=9091
```

### Financial Service
```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=gspend
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key
PORT=8082
AUTH_SERVICE_GRPC_ADDR=localhost:9091
```

---

## Database Setup

MongoDB and Redis are included in Docker Compose. For manual setup:

```bash
# Start MongoDB
docker run -d --name mongodb -p 27017:27017 mongo:7

# Start Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

---

## Useful Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Start all services |
| `make docker-down` | Stop all services |
| `make docker-logs` | View logs |
| `make demo-start` | Start with sample data |
| `make demo-stop` | Stop demo environment |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage |
| `make health-check` | Check service health |

---

## Troubleshooting

### Port conflicts
```bash
# Check if ports are in use
lsof -i :80 -i :8081 -i :8082 -i :27017 -i :6379
```

### Reset database
```bash
make db-reset
```

### Clean Docker resources
```bash
make docker-clean
```
