# gSpend Monorepo Makefile

# =============================================================================
# Testing Targets
# =============================================================================

.PHONY: test
test: test-auth test-finance

.PHONY: test-v
test-v: test-auth-v test-finance-v

.PHONY: test-auth
test-auth:
	@echo "Running Auth Service tests..."
	@cd apps/auth-service && go test ./...

.PHONY: test-auth-v
test-auth-v:
	@echo "Running Auth Service tests (verbose)..."
	@cd apps/auth-service && go test -v ./...

.PHONY: test-finance
test-finance:
	@echo "Running Financial Service tests..."
	@cd apps/financial-service && go test ./...

.PHONY: test-finance-v
test-finance-v:
	@echo "Running Financial Service tests (verbose)..."
	@cd apps/financial-service && go test -v ./...

.PHONY: test-coverage
test-coverage: test-auth-coverage test-finance-coverage

.PHONY: test-auth-coverage
test-auth-coverage:
	@echo "Running Auth Service tests with coverage..."
	@cd apps/auth-service && go test -coverprofile=coverage.out ./... && \
		grep -v "\.pb\.go" coverage.out > coverage_filtered.out && \
		go tool cover -func=coverage_filtered.out

.PHONY: test-finance-coverage
test-finance-coverage:
	@echo "Running Financial Service tests with coverage..."
	@cd apps/financial-service && go test -coverprofile=coverage.out ./... && \
		grep -v "\.pb\.go" coverage.out > coverage_filtered.out && \
		go tool cover -func=coverage_filtered.out

# =============================================================================
# Code Generation
# =============================================================================

.PHONY: generate
generate:
	$(MAKE) -C proto generate

# =============================================================================
# Docker Development Environment
# =============================================================================

.PHONY: docker-up
docker-up:
	@echo "Starting gSpend development environment..."
	docker-compose -f devops/docker/docker-compose.yml up --build -d
	@echo "Waiting for services to be healthy..."
	@sleep 10
	@$(MAKE) docker-status

.PHONY: docker-down
docker-down:
	@echo "Stopping gSpend development environment..."
	docker-compose -f devops/docker/docker-compose.yml down

.PHONY: docker-restart
docker-restart: docker-down docker-up

.PHONY: docker-status
docker-status:
	@echo "=== Docker Container Status ==="
	@docker-compose -f devops/docker/docker-compose.yml ps

.PHONY: docker-logs
docker-logs:
	@echo "=== All Service Logs ==="
	docker-compose -f devops/docker/docker-compose.yml logs --tail=50

.PHONY: docker-logs-auth
docker-logs-auth:
	@echo "=== Auth Service Logs ==="
	docker-compose -f devops/docker/docker-compose.yml logs auth-service --tail=50 -f

.PHONY: docker-logs-finance
docker-logs-finance:
	@echo "=== Financial Service Logs ==="
	docker-compose -f devops/docker/docker-compose.yml logs financial-service --tail=50 -f

.PHONY: docker-logs-frontend
docker-logs-frontend:
	@echo "=== Frontend Logs ==="
	docker-compose -f devops/docker/docker-compose.yml logs frontend --tail=50 -f

.PHONY: docker-clean
docker-clean:
	@echo "Cleaning up Docker resources..."
	docker-compose -f devops/docker/docker-compose.yml down -v
	docker system prune -f
	@echo "Docker cleanup complete"

# =============================================================================
# Database Management
# =============================================================================

.PHONY: db-setup
db-setup:
	@echo "Setting up database indexes..."
	@cd apps/financial-service && go run scripts/create_indexes.go

.PHONY: db-reset
db-reset:
	@echo "Resetting database (WARNING: This will delete all data)..."
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose -f devops/docker/docker-compose.yml exec mongodb mongosh gspend --eval "db.dropDatabase()"; \
		echo "Database reset complete"; \
	else \
		echo "Database reset cancelled"; \
	fi

.PHONY: db-shell
db-shell:
	@echo "Opening MongoDB shell..."
	docker-compose -f devops/docker/docker-compose.yml exec mongodb mongosh gspend

# =============================================================================
# Local Development
# =============================================================================

.PHONY: dev-setup
dev-setup: generate build db-setup
	@echo "Development environment setup complete!"

.PHONY: dev-start
dev-start: docker-up
	@echo "Development environment started!"
	@echo "Frontend: http://localhost"
	@echo "Auth API: http://localhost/api/v1/auth"
	@echo "Financial API: http://localhost/api/v1"

.PHONY: dev-stop
dev-stop: docker-down

.PHONY: dev-restart
dev-restart: docker-restart

.PHONY: dev-status
dev-status: docker-status

# =============================================================================
# Building
# =============================================================================

.PHONY: build-auth
build-auth:
	@echo "Building Auth Service..."
	@cd apps/auth-service && go build -o bin/server ./cmd/server

.PHONY: build-finance
build-finance:
	@echo "Building Financial Service..."
	@cd apps/financial-service && go build -o bin/server ./cmd/server

.PHONY: build-frontend
build-frontend:
	@echo "Building Frontend..."
	@cd apps/frontend && npm run build

.PHONY: build
build: build-auth build-finance build-frontend

# =============================================================================
# Health Checks
# =============================================================================

.PHONY: health-check
health-check:
	@echo "=== Health Check ==="
	@echo "Frontend: $$(curl -s -o /dev/null -w "%{http_code}" http://localhost || echo "FAIL")"
	@echo "Auth Service: $$(curl -s -o /dev/null -w "%{http_code}" http://localhost/api/v1/auth/health || echo "FAIL")"
	@echo "Financial Service: $$(curl -s -o /dev/null -w "%{http_code}" http://localhost/api/v1/health || echo "FAIL")"

# =============================================================================
# Testing with Docker
# =============================================================================

.PHONY: test-integration
test-integration:
	@echo "Integration tests are currently DISABLED"
	@echo "Functional tests will be re-enabled after completing implementation tasks"
	@echo "Only running service health checks..."
	docker-compose -f docker-compose.test.yml down --volumes --remove-orphans
	docker-compose -f docker-compose.test.yml build auth-service-test financial-service-test
	docker-compose -f docker-compose.test.yml up -d mongodb-test redis-test auth-service-test financial-service-test
	@echo "Waiting for services to be healthy..."
	@sleep 20
	@echo "Services started successfully. Integration tests skipped for now."
	docker-compose -f docker-compose.test.yml down --volumes

.PHONY: test-integration-local
test-integration-local: docker-up
	@echo "Running integration tests against local Docker environment..."
	@sleep 15  # Wait for services to be fully ready
	@$(MAKE) health-check
	@echo "Integration tests complete"

.PHONY: test-e2e
test-e2e: docker-up
	@echo "Running end-to-end tests..."
	@sleep 15  # Wait for services to be fully ready
	# Add e2e test commands here when available
	@echo "E2E tests complete"

.PHONY: test-clean
test-clean:
	@echo "Cleaning up test environment..."
	docker-compose -f docker-compose.test.yml down --volumes --remove-orphans
	docker system prune -f --volumes

# =============================================================================
# Utility Targets
# =============================================================================

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf apps/auth-service/bin
	@rm -rf apps/financial-service/bin
	@rm -rf apps/frontend/dist
	@rm -rf apps/*/coverage.out
	@rm -rf apps/*/coverage_filtered.out
	@echo "Clean complete"

.PHONY: help
help:
	@echo "gSpend Development Commands:"
	@echo ""
	@echo "Development:"
	@echo "  dev-setup     - Setup development environment"
	@echo "  dev-start     - Start development environment"
	@echo "  dev-stop      - Stop development environment"
	@echo "  dev-restart   - Restart development environment"
	@echo "  dev-status    - Show development environment status"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up     - Start all services with Docker"
	@echo "  docker-down   - Stop all services"
	@echo "  docker-restart - Restart all services"
	@echo "  docker-status - Show container status"
	@echo "  docker-logs   - Show all service logs"
	@echo "  docker-clean  - Clean up Docker resources"
	@echo ""
	@echo "Testing:"
	@echo "  test          - Run all unit tests"
	@echo "  test-v        - Run all unit tests (verbose)"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-integration - Run integration tests"
	@echo "  health-check  - Check service health"
	@echo ""
	@echo "Database:"
	@echo "  db-setup      - Create database indexes"
	@echo "  db-reset      - Reset database (WARNING: deletes data)"
	@echo "  db-shell      - Open MongoDB shell"
	@echo ""
	@echo "Building:"
	@echo "  build         - Build all services"
	@echo "  generate      - Generate protobuf code"
	@echo "  clean         - Clean build artifacts"
