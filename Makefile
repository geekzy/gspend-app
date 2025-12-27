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

# =============================================================================
# Demo Environment with Dummy Data
# =============================================================================

.PHONY: demo-start
demo-start:
	@echo "🚀 Starting GSpend Demo Environment..."
	@echo "This will create a complete demo environment with dummy data."
	@echo ""
	@if ! docker info > /dev/null 2>&1; then \
		echo "❌ Docker is not running. Please start Docker and try again."; \
		exit 1; \
	fi
	@echo "🧹 Cleaning up existing containers..."
	@docker-compose -f docker-compose.demo.yml down -v > /dev/null 2>&1 || true
	@echo "🔨 Building and starting services..."
	@docker-compose -f docker-compose.demo.yml up --build -d
	@echo "⏳ Waiting for services to be ready..."
	@echo "This may take a few minutes on first run..."
	@$(MAKE) demo-wait-for-services
	@echo ""
	@echo "🎉 Demo environment is ready!"
	@echo ""
	@echo "📊 Dashboard Preview:"
	@echo "   URL: http://localhost"
	@echo ""
	@echo "🔐 Demo Login Credentials:"
	@echo "   Email: demo@gspend.com"
	@echo "   Password: password"
	@echo ""
	@echo "📈 What's included:"
	@echo "   • 3 months of sample transactions"
	@echo "   • Multiple expense categories (Food, Transport, Shopping, etc.)"
	@echo "   • Income records (Salary, Freelance, Investments)"
	@echo "   • Monthly budget with spending tracking"
	@echo "   • Dashboard with charts and analytics"
	@echo ""
	@echo "🛠️  Useful commands:"
	@echo "   • View logs: make demo-logs"
	@echo "   • Stop demo: make demo-stop"
	@echo "   • Restart: make demo-restart"

.PHONY: demo-stop
demo-stop:
	@echo "🛑 Stopping GSpend Demo Environment..."
	@docker-compose -f docker-compose.demo.yml down
	@echo "✅ Demo environment stopped."
	@echo ""
	@echo "💡 To completely remove demo data:"
	@echo "   make demo-clean"
	@echo ""
	@echo "🚀 To start again:"
	@echo "   make demo-start"

.PHONY: demo-restart
demo-restart: demo-stop demo-start

.PHONY: demo-clean
demo-clean:
	@echo "🧹 Cleaning up demo environment and data..."
	@docker-compose -f docker-compose.demo.yml down -v
	@docker system prune -f
	@echo "✅ Demo environment and data cleaned."

.PHONY: demo-logs
demo-logs:
	@echo "=== Demo Environment Logs ==="
	@docker-compose -f docker-compose.demo.yml logs --tail=50

.PHONY: demo-logs-follow
demo-logs-follow:
	@echo "=== Following Demo Environment Logs ==="
	@docker-compose -f docker-compose.demo.yml logs -f

.PHONY: demo-status
demo-status:
	@echo "=== Demo Environment Status ==="
	@docker-compose -f docker-compose.demo.yml ps

.PHONY: demo-wait-for-services
demo-wait-for-services:
	@echo "⏳ Checking Auth Service..."
	@timeout=300; \
	while [ $$timeout -gt 0 ]; do \
		if curl -s -f http://localhost:8081/api/v1/auth/health > /dev/null 2>&1; then \
			echo "✅ Auth Service is ready"; \
			break; \
		fi; \
		sleep 5; \
		timeout=$$((timeout-5)); \
	done
	@echo "⏳ Checking Financial Service..."
	@timeout=300; \
	while [ $$timeout -gt 0 ]; do \
		if curl -s -f http://localhost:8082/api/v1/health > /dev/null 2>&1; then \
			echo "✅ Financial Service is ready"; \
			break; \
		fi; \
		sleep 5; \
		timeout=$$((timeout-5)); \
	done
	@echo "⏳ Checking Frontend..."
	@timeout=300; \
	while [ $$timeout -gt 0 ]; do \
		if curl -s -f http://localhost/ > /dev/null 2>&1; then \
			echo "✅ Frontend is ready"; \
			break; \
		fi; \
		sleep 5; \
		timeout=$$((timeout-5)); \
	done

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
	@echo "🧪 Running Integration Tests..."
	@echo "Setting up test environment..."
	@docker-compose -f docker-compose.test.yml down --volumes --remove-orphans > /dev/null 2>&1 || true
	@docker-compose -f docker-compose.test.yml build
	@docker-compose -f docker-compose.test.yml up -d mongodb-test redis-test auth-service-test financial-service-test
	@echo "⏳ Waiting for services to be healthy..."
	@$(MAKE) test-wait-for-services
	@echo "🚀 Running integration test suite..."
	@if docker run --rm --network gspend-test-network \
		-e AUTH_SERVICE_URL=http://auth-service-test:8081 \
		-e FINANCIAL_SERVICE_URL=http://financial-service-test:8082 \
		-e MONGODB_URI=mongodb://mongodb-test:27017 \
		-e MONGODB_DATABASE=gspend_test \
		-v $(PWD)/test-results:/app/test-results \
		$$(docker build -q -f Dockerfile.test .); then \
		echo "✅ Integration tests PASSED"; \
	else \
		echo "❌ Integration tests FAILED"; \
		$(MAKE) test-show-results; \
		exit 1; \
	fi
	@$(MAKE) test-show-results
	@docker-compose -f docker-compose.test.yml down --volumes

.PHONY: test-integration-quick
test-integration-quick:
	@echo "🧪 Running Quick Integration Tests (reuse existing containers)..."
	@if docker ps | grep -q gspend-auth-test && docker ps | grep -q gspend-financial-test; then \
		echo "♻️  Using existing test containers..."; \
		docker run --rm --network gspend-test-network \
			-e AUTH_SERVICE_URL=http://auth-service-test:8081 \
			-e FINANCIAL_SERVICE_URL=http://financial-service-test:8082 \
			-e MONGODB_URI=mongodb://mongodb-test:27017 \
			-e MONGODB_DATABASE=gspend_test \
			-v $(PWD)/test-results:/app/test-results \
			$$(docker build -q -f Dockerfile.test .); \
	else \
		echo "⚠️  No existing test containers found. Running full integration test..."; \
		$(MAKE) test-integration; \
	fi

.PHONY: test-integration-setup
test-integration-setup:
	@echo "🔧 Setting up integration test environment..."
	@docker-compose -f docker-compose.test.yml down --volumes --remove-orphans > /dev/null 2>&1 || true
	@docker-compose -f docker-compose.test.yml build
	@docker-compose -f docker-compose.test.yml up -d mongodb-test redis-test auth-service-test financial-service-test
	@$(MAKE) test-wait-for-services
	@echo "✅ Integration test environment ready"

.PHONY: test-integration-teardown
test-integration-teardown:
	@echo "🧹 Tearing down integration test environment..."
	@docker-compose -f docker-compose.test.yml down --volumes --remove-orphans
	@echo "✅ Integration test environment cleaned up"

.PHONY: test-wait-for-services
test-wait-for-services:
	@echo "⏳ Waiting for Auth Service..."
	@timeout=180; \
	while [ $$timeout -gt 0 ]; do \
		if curl -s -f http://localhost:8083/api/v1/auth/health > /dev/null 2>&1; then \
			echo "✅ Auth Service is ready"; \
			break; \
		fi; \
		sleep 3; \
		timeout=$$((timeout-3)); \
	done; \
	if [ $$timeout -le 0 ]; then echo "❌ Auth Service failed to start"; exit 1; fi
	@echo "⏳ Waiting for Financial Service..."
	@timeout=180; \
	while [ $$timeout -gt 0 ]; do \
		if curl -s -f http://localhost:8084/api/v1/health > /dev/null 2>&1; then \
			echo "✅ Financial Service is ready"; \
			break; \
		fi; \
		sleep 3; \
		timeout=$$((timeout-3)); \
	done; \
	if [ $$timeout -le 0 ]; then echo "❌ Financial Service failed to start"; exit 1; fi

.PHONY: test-show-results
test-show-results:
	@if [ -f test-results/integration-test-results.json ]; then \
		echo ""; \
		echo "📊 Integration Test Results:"; \
		echo "============================"; \
		cat test-results/integration-test-results.json | jq -r '. | "Total Tests: \(.total_tests)\nPassed: \(.tests_passed)\nFailed: \(.tests_failed)\nSuccess Rate: \(.success_rate)%\nStatus: \(.status)"'; \
		echo ""; \
		echo "📄 Full report: test-results/integration-test-results.json"; \
	else \
		echo "⚠️  No test results found"; \
	fi

.PHONY: test-integration-logs
test-integration-logs:
	@echo "=== Integration Test Environment Logs ==="
	@docker-compose -f docker-compose.test.yml logs --tail=50

.PHONY: test-integration-logs-follow
test-integration-logs-follow:
	@echo "=== Following Integration Test Environment Logs ==="
	@docker-compose -f docker-compose.test.yml logs -f

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
	@echo "  dev-setup                  - Setup development environment"
	@echo "  dev-start                  - Start development environment"
	@echo "  dev-stop                   - Stop development environment"
	@echo "  dev-restart                - Restart development environment"
	@echo "  dev-status                 - Show development environment status"
	@echo ""
	@echo "Demo Environment:"
	@echo "  demo-start                 - Start demo environment with dummy data"
	@echo "  demo-stop                  - Stop demo environment"
	@echo "  demo-restart               - Restart demo environment"
	@echo "  demo-clean                 - Clean demo environment and data"
	@echo "  demo-logs                  - Show demo environment logs"
	@echo "  demo-status                - Show demo environment status"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up                  - Start all services with Docker"
	@echo "  docker-down                - Stop all services"
	@echo "  docker-restart             - Restart all services"
	@echo "  docker-status              - Show container status"
	@echo "  docker-logs                - Show all service logs"
	@echo "  docker-clean               - Clean up Docker resources"
	@echo ""
	@echo "Testing:"
	@echo "  test                       - Run all unit tests"
	@echo "  test-v                     - Run all unit tests (verbose)"
	@echo "  test-coverage              - Run tests with coverage"
	@echo "  test-integration           - Run full integration tests"
	@echo "  test-integration-quick     - Run integration tests (reuse containers)"
	@echo "  test-integration-setup     - Setup integration test environment"
	@echo "  test-integration-teardown  - Teardown integration test environment"
	@echo "  test-integration-logs      - Show integration test logs"
	@echo "  health-check               - Check service health"
	@echo ""
	@echo "Database:"
	@echo "  db-setup                   - Create database indexes"
	@echo "  db-reset                   - Reset database (WARNING: deletes data)"
	@echo "  db-shell                   - Open MongoDB shell"
	@echo ""
	@echo "Building:"
	@echo "  build                      - Build all services"
	@echo "  generate                   - Generate protobuf code"
	@echo "  clean                      - Clean build artifacts"
