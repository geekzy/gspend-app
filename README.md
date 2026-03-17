# gSpend - Family Financial Management App

gSpend is a modern, full-stack microservices application designed for families to manage their personal finances. It allows users to track income, plan budgets, and monitor transactions with a beautiful, responsive user interface.

## 🚀 Features

- **User Authentication**: Secure registration and login with JWT.
- **Dashboard**: Overview of total balance, active budgets, and recent transactions.
- **Income Management**: Track various income sources (salary, side-hustles, etc.).
- **Budget Planning**: Define monthly budgets by category and track spending progress.
- **Transaction Tracking**: Categorized income and expense logging with real-time budget updates.
- **Family Focused**: Clean Architecture designed for multi-user family groups (expandable).

## 🛠 Tech Stack

### Backend (Go)

- **Microservices Architecture**: Auth Service & Financial Service.
- **API Framework**: Echo (REST) & gRPC (Inter-service).
- **Data Storage**: MongoDB (Primary) & Redis (Caching).
- **Communication**: Protocol Buffers / gRPC.

### Frontend (Vue.js)

- **Framework**: Vue 3 with TypeScript (Vite).
- **State Management**: Pinia.
- **Styling**: Vanilla CSS with Tailwind-like utility patterns.
- **Icons**: Lucide Vue Next.

### Infrastructure

- **API Gateway**: Nginx.
- **Containerization**: Docker & Docker Compose.

---

## 💻 Local Development Setup

To run the entire gSpend stack locally:

### Prerequisites

- [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
- Go 1.21+ (if running services outside Docker).

### Quick Start

1. Clone the repository.
2. Run the application using the root Makefile:

```bash
make docker-up
```
  
Access the application:

- **Frontend**: [http://localhost](http://localhost) (via Nginx)
- **Auth API**: `http://localhost/api/v1/auth`

### Demo Environment with Sample Data

For a complete preview with realistic dummy data:

```bash
make demo-start
```

- **Login**: demo@gspend.com / passw0rd!
- **Includes**: 3 months of transactions, budgets, and dashboard analytics
- **See**: [DEMO.md](DEMO.md) for full details

### Configuration

Services use a hierarchical configuration system:

1. **`config.yaml`**: Main configuration file (e.g., `apps/auth-service/config.yaml`).
2. **Environment Variables**: Override `config.yaml` values using a mapping (e.g., `SMTP_USER` overrides `smtp.user`).

### Available Commands

```bash
make help                    # Show all available commands

# Development
make dev-start              # Start development environment  
make dev-stop               # Stop development environment

# Demo with dummy data
make demo-start             # Start demo with sample data
make demo-stop              # Stop demo environment

# Testing
make test                   # Run unit tests
make test-integration       # Run integration tests
make test-coverage          # Run tests with coverage

# Docker management
make docker-up              # Start services
make docker-down            # Stop services
make docker-logs            # View logs
```

---

## 🧪 Testing

GSpend includes comprehensive testing at multiple levels to ensure reliability and correctness.

### Unit Tests

Run unit tests for individual services:

```bash
make test                    # Run all unit tests
make test-v                  # Run with verbose output
make test-coverage           # Run with coverage reports
```

Individual service testing:

```bash
make test-auth              # Test auth service only
make test-finance           # Test financial service only
```

### Integration Tests

Comprehensive end-to-end testing in containerized environments:

```bash
make test-integration       # Full integration test suite
make test-integration-quick # Quick tests (reuse containers)
```

#### Integration Test Features

- **Isolated Environment**: Separate test database and Redis instances
- **Automated Data Seeding**: Categories and test data automatically created
- **Service Health Checks**: Validates all services start correctly
- **API Endpoint Testing**: Tests authentication, financial APIs, and database connectivity
- **Comprehensive Reporting**: JSON test results with detailed metrics

#### Test Coverage Areas

1. **Authentication Service**
   - User registration and login flows
   - JWT token validation and refresh
   - Profile management operations
   - Password change functionality

2. **Financial Service**
   - Categories and transactions endpoints
   - Dashboard data aggregation
   - Budget management operations
   - Income tracking functionality

3. **Database Operations**
   - MongoDB connectivity and operations
   - Data seeding and validation
   - Index creation and performance

4. **Service Integration**
   - Inter-service communication (gRPC)
   - API gateway routing
   - Health check endpoints

### Test Environment Management

```bash
# Setup test environment and keep running
make test-integration-setup

# Run quick tests against existing environment
make test-integration-quick

# View test logs for debugging
make test-integration-logs

# Clean up test environment
make test-integration-teardown
```

### Test Results

Integration tests generate detailed JSON reports:

```bash
# View latest test results
make test-show-results

# Results saved to: test-results/integration-test-results.json
```

### CI/CD Integration

Tests are designed for automated CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Integration Tests
  run: make test-integration
- name: Upload Test Results
  uses: actions/upload-artifact@v3
  with:
    name: test-results
    path: test-results/
```

### Testing Documentation

- **[TESTING-INTEGRATION.md](TESTING-INTEGRATION.md)** - Comprehensive integration testing guide
- **Test Architecture** - Containerized testing with Docker Compose
- **Debugging Guide** - Troubleshooting test failures and environment issues

---

## 🌐 Production Deployment Guidelines

For deploying gSpend to a production environment, consider the following best practices:

### 1. Environment Variables

Ensure all sensitive keys and production-specific values are set via environment variables. Key variables include:

- `JWT_SECRET`: A strong, unique secret key for token signing.
- `MONGODB_URI`: Connection string to a production MongoDB cluster (e.g., MongoDB Atlas).
- `REDIS_HOST` & `REDIS_PORT`: Production Redis instance.
- `APP_ENV`: Set to `production`.

### 2. Infrastructure

- **Orchestration**: Use Kubernetes or Docker Swarm for better scalability and health monitoring.
- **Managed Databases**: Use managed services for MongoDB (Atlas) and Redis (Elasticache/Upstash) to ensure high availability and backups.
- **SSL/TLS**: Configure Nginx (or a cloud load balancer like AWS ALB/GCP Load Balancer) to handle HTTPS certificates via Let's Encrypt or ACM.

### 3. Monitoring & Logging

- Integrate monitoring tools like Prometheus/Grafana or Datadog.
- Centralize logs using ELK stack or cloud-native logging services.

### 4. Continuous Integration (CI)

Automate tests and builds using GitHub Actions. GSpend includes comprehensive testing infrastructure:

```yaml
# Example GitHub Actions workflow
name: CI/CD Pipeline
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      # Unit Tests
      - name: Run Unit Tests
        run: make test-coverage
      
      # Integration Tests
      - name: Run Integration Tests
        run: make test-integration
      
      # Upload Results
      - name: Upload Test Results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: test-results
          path: test-results/

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Docker Images
        run: make build
      # Deploy to staging/production
```

**Testing Features:**

- **Unit Tests**: Individual service testing with coverage reports
- **Integration Tests**: End-to-end API testing in containerized environment
- **Automated Reporting**: JSON test results for CI/CD integration
- **Performance Tracking**: Test execution time and success rate monitoring

---

## 📂 Project Structure

```text
├── apps/
│   ├── auth-service/       # Identity and Access Management
│   ├── financial-service/  # Financial data and logic
│   └── frontend/           # Vue.js SPA
├── devops/
│   ├── docker/             # Docker Compose configurations
│   └── nginx/              # API Gateway configurations
├── proto/                  # gRPC protocol definitions
├── scripts/                # Utility scripts and data seeding
│   ├── seed_dummy_data.go  # Demo data generation
│   └── integration-test.sh # Integration test execution
├── docs/                   # Architecture and Design docs
├── test-results/           # Integration test results
├── docker-compose.demo.yml # Demo environment configuration
├── docker-compose.test.yml # Test environment configuration
├── DEMO.md                 # Demo environment guide
├── TESTING-INTEGRATION.md  # Integration testing guide
└── Makefile                # Centralized command runner
```

---

## 🚀 Quick Reference

### Essential Commands

```bash
# Get started quickly
make demo-start             # Complete demo with sample data
make dev-start              # Development environment
make test-integration       # Run all tests

# View all available commands
make help
```

### Access Points

- **Demo Dashboard**: http://localhost (demo@gspend.com / passw0rd!)
- **Development Frontend**: http://localhost
- **Auth API**: http://localhost/api/v1/auth  
- **Financial API**: http://localhost/api/v1

### Documentation

- **[DEMO.md](DEMO.md)** - Demo environment setup and usage
- **[TESTING-INTEGRATION.md](TESTING-INTEGRATION.md)** - Comprehensive testing guide
- **[Architecture Docs](docs/)** - Technical architecture and design documents

---

## 📝 License

This project is licensed under the MIT License.
