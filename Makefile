# Monorepo Makefile

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
	@cd apps/auth-service && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep -v "\.pb\.go" | grep -v "total:"

.PHONY: test-finance-coverage
test-finance-coverage:
	@echo "Running Financial Service tests with coverage..."
	@cd apps/financial-service && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep -v "\.pb\.go" | grep -v "total:"

.PHONY: generate
generate:
	$(MAKE) -C proto generate

.PHONY: docker-up
docker-up:
	docker-compose -f devops/docker/docker-compose.yml up --build -d

.PHONY: docker-down
docker-down:
	docker-compose -f devops/docker/docker-compose.yml down

.PHONY: build-auth
build-auth:
	cd apps/auth-service && go build -o bin/server ./cmd/server

.PHONY: build-finance
build-finance:
	cd apps/financial-service && go build -o bin/server ./cmd/server

.PHONY: build
build: build-auth build-finance
