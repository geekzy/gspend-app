# gSpend - Family Financial Management App

gSpend is a modern, full-stack microservices application designed for families to manage their personal finances. It allows users to track income, plan budgets, and monitor transactions with a beautiful, responsive user interface.

## 🚀 Features

-   **User Authentication**: Secure registration and login with JWT.
-   **Dashboard**: Overview of total balance, active budgets, and recent transactions.
-   **Income Management**: Track various income sources (salary, side-hustles, etc.).
-   **Budget Planning**: Define monthly budgets by category and track spending progress.
-   **Transaction Tracking**: Categorized income and expense logging with real-time budget updates.
-   **Family Focused**: Clean Architecture designed for multi-user family groups (expandable).

## 🛠 Tech Stack

### Backend (Go)
-   **Microservices Architecture**: Auth Service & Financial Service.
-   **API Framework**: Echo (REST) & gRPC (Inter-service).
-   **Data Storage**: MongoDB (Primary) & Redis (Caching).
-   **Communication**: Protocol Buffers / gRPC.

### Frontend (Vue.js)
-   **Framework**: Vue 3 with TypeScript (Vite).
-   **State Management**: Pinia.
-   **Styling**: Vanilla CSS with Tailwind-like utility patterns.
-   **Icons**: Lucide Vue Next.

### Infrastructure
-   **API Gateway**: Nginx.
-   **Containerization**: Docker & Docker Compose.

---

## 💻 Local Development Setup

To run the entire gSpend stack locally:

### Prerequisites
-   [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
-   Go 1.21+ (if running services outside Docker).

### Quick Start
1.  Clone the repository.
2.  Run the application using the root Makefile:
    ```bash
    make docker-up
    ```
3.  Access the application:
    -   **Frontend**: [http://localhost](http://localhost) (via Nginx)
    -   **Auth API**: `http://localhost/api/v1/auth`
    -   **Finance API**: `http://localhost/api/v1`

### Useful Makefile Commands
-   `make test`: Run all backend tests.
-   `make test-coverage`: Run all backend tests and display code coverage percentages (aims for 80%+ on critical logic).
-   `make build`: Build backend binaries.
-   `make docker-down`: Stop and remove containers.
-   `make generate`: Re-generate gRPC code from proto files.

---

## 🌐 Production Deployment Guidelines

For deploying gSpend to a production environment, consider the following best practices:

### 1. Environment Variables
Ensure all sensitive keys and production-specific values are set via environment variables. Key variables include:
-   `JWT_SECRET`: A strong, unique secret key for token signing.
-   `MONGODB_URI`: Connection string to a production MongoDB cluster (e.g., MongoDB Atlas).
-   `REDIS_HOST` & `REDIS_PORT`: Production Redis instance.
-   `APP_ENV`: Set to `production`.

### 2. Infrastructure
-   **Orchestration**: Use Kubernetes or Docker Swarm for better scalability and health monitoring.
-   **Managed Databases**: Use managed services for MongoDB (Atlas) and Redis (Elasticache/Upstash) to ensure high availability and backups.
-   **SSL/TLS**: Configure Nginx (or a cloud load balancer like AWS ALB/GCP Load Balancer) to handle HTTPS certificates via Let's Encrypt or ACM.

### 3. Monitoring & Logging
-   Integrate monitoring tools like Prometheus/Grafana or Datadog.
-   Centralize logs using ELK stack or cloud-native logging services.

### 4. Continuous Integration (CI)
Automate tests and builds using GitHub Actions. A typical flow:
1.  Run `go test` and frontend linting on every PR.
2.  Build Docker images and push to a registry (Docker Hub, ECR).
3.  Deploy to a staging/production cluster upon merging to `main`.

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
├── docs/                   # Architecture and Design docs
└── Makefile                # Centralized command runner
```

## 📝 License
This project is licensed under the MIT License.
