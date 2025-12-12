.PHONY: help build up down logs clean test verify seed seed-bash demo-gateway demo-sorting lint tidy fmt

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d
	@echo "Waiting for services to start..."
	@sleep 5
	@echo "Services started! Check health at:"
	@echo "  Console: http://localhost:8080/healthz"
	@echo "  ODFI:    http://localhost:8081/healthz"
	@echo "  RDFI:    http://localhost:8082/healthz"
	@echo "  Ledger:  http://localhost:8083/healthz"
	@echo "  EIP:     http://localhost:8084/healthz"

down: ## Stop all services
	docker-compose down

logs: ## Follow logs for all services
	docker-compose logs -f

logs-%: ## Follow logs for a specific service (e.g., make logs-odfi)
	docker-compose logs -f $*

clean: ## Stop services and remove volumes
	docker-compose down -v

test: ## Run Go tests
	go test -v ./...

verify: ## Verify system is working
	./verify.sh

demo-gateway: ## Run complete gateway demo (all operations via gateway)
	./demo-gateway.sh

demo-sorting: ## Demo all sorting options for unified ACH items
	./demo-sorting.sh

# Development targets
dev-odfi: ## Run ODFI service locally (requires PostgreSQL)
	@echo "Make sure PostgreSQL is running locally"
	PORT=8081 go run cmd/odfi/main.go

dev-rdfi: ## Run RDFI service locally (requires PostgreSQL)
	@echo "Make sure PostgreSQL is running locally"
	PORT=8082 go run cmd/rdfi/main.go

dev-ledger: ## Run Ledger service locally (requires PostgreSQL)
	@echo "Make sure PostgreSQL is running locally"
	PORT=8083 go run cmd/ledger/main.go

dev-eip: ## Run EIP service locally (requires PostgreSQL)
	@echo "Make sure PostgreSQL is running locally"
	PORT=8084 go run cmd/eip/main.go

dev-console: ## Run Console service locally
	PORT=8080 \
	ODFI_BASE_URL=http://localhost:8081 \
	RDFI_BASE_URL=http://localhost:8082 \
	LEDGER_BASE_URL=http://localhost:8083 \
	EIP_BASE_URL=http://localhost:8084 \
	go run cmd/console/main.go

tidy: ## Tidy Go modules
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

seed: ## Seed databases with 3700 records (1200 ODFI + 1200 RDFI + 800 Ledger + 500 EIP)
	@echo "🌱 Seeding databases (1200 ODFI, 1200 RDFI, 800 Ledger, 500 EIP)..."
	@go run cmd/seed/main.go

seed-bash: ## Seed databases with 620 records (slower bash version)
	@./seed.sh

lint: ## Run linters (requires golangci-lint)
	golangci-lint run

# Docker targets
docker-build-%: ## Build specific service Docker image (e.g., make docker-build-odfi)
	docker-compose build $*

docker-restart-%: ## Restart specific service (e.g., make docker-restart-odfi)
	docker-compose restart $*

# Database targets
psql-odfi: ## Connect to ODFI database
	docker-compose exec odfi-db psql -U odfi_user -d odfi_db

psql-rdfi: ## Connect to RDFI database
	docker-compose exec rdfi-db psql -U rdfi_user -d rdfi_db

psql-ledger: ## Connect to Ledger database
	docker-compose exec ledger-db psql -U ledger_user -d ledger_db

psql-eip: ## Connect to EIP database
	docker-compose exec eip-db psql -U eip_user -d eip_db

