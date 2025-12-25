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
