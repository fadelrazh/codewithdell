# CodeWithDell Makefile
# Advanced level development tasks

.PHONY: help setup dev build test clean docker-build docker-run docker-stop lint format migrate seed backup restore

# Default target
help: ## Show this help message
	@echo "CodeWithDell Development Commands"
	@echo "=================================="
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development setup
setup: ## Initial setup of the development environment
	@echo "🚀 Setting up CodeWithDell development environment..."
	@./scripts/dev.sh setup

dev: ## Start development servers
	@echo "🔥 Starting development servers..."
	@./scripts/dev.sh start

stop: ## Stop all development services
	@echo "🛑 Stopping all services..."
	@./scripts/dev.sh stop

status: ## Show service status
	@echo "📊 Checking service status..."
	@./scripts/dev.sh status

logs: ## Show logs (backend|frontend|postgres|redis|all)
	@echo "📝 Showing logs..."
	@./scripts/dev.sh logs $(filter-out $@,$(MAKECMDGOALS))

# Database operations
migrate: ## Run database migrations
	@echo "🗄️  Running database migrations..."
	@cd backend && go run main.go migrate

seed: ## Seed database with sample data
	@echo "🌱 Seeding database..."
	@cd backend && go run main.go seed

reset-db: ## Reset database (WARNING: This will delete all data)
	@echo "⚠️  Resetting database..."
	@./scripts/dev.sh reset-db

backup: ## Create database backup
	@echo "💾 Creating database backup..."
	@mkdir -p backups
	@docker exec codewithdell_postgres pg_dump -U codewithdell codewithdell > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

restore: ## Restore database from backup (usage: make restore BACKUP=backup_file.sql)
	@echo "📥 Restoring database from $(BACKUP)..."
	@docker exec -i codewithdell_postgres psql -U codewithdell codewithdell < backups/$(BACKUP)

# Building and testing
build: ## Build all applications
	@echo "🔨 Building applications..."
	@cd backend && go build -o bin/backend .
	@cd frontend && npm run build

test: ## Run all tests
	@echo "🧪 Running tests..."
	@cd backend && go test ./...
	@cd frontend && npm test

test-backend: ## Run backend tests
	@echo "🧪 Running backend tests..."
	@cd backend && go test ./...

test-frontend: ## Run frontend tests
	@echo "🧪 Running frontend tests..."
	@cd frontend && npm test

# Code quality
lint: ## Run linting for all code
	@echo "🔍 Running linters..."
	@cd backend && golangci-lint run
	@cd frontend && npm run lint

lint-backend: ## Run backend linting
	@echo "🔍 Running backend linter..."
	@cd backend && golangci-lint run

lint-frontend: ## Run frontend linting
	@echo "🔍 Running frontend linter..."
	@cd frontend && npm run lint

format: ## Format all code
	@echo "✨ Formatting code..."
	@cd backend && go fmt ./...
	@cd frontend && npm run format

format-backend: ## Format backend code
	@echo "✨ Formatting backend code..."
	@cd backend && go fmt ./...

format-frontend: ## Format frontend code
	@echo "✨ Formatting frontend code..."
	@cd frontend && npm run format

# Docker operations
docker-build: ## Build Docker images
	@echo "🐳 Building Docker images..."
	@docker-compose build

docker-run: ## Run with Docker Compose
	@echo "🐳 Running with Docker Compose..."
	@docker-compose up -d

docker-stop: ## Stop Docker Compose services
	@echo "🐳 Stopping Docker Compose services..."
	@docker-compose down

docker-logs: ## Show Docker logs
	@echo "🐳 Showing Docker logs..."
	@docker-compose logs -f

docker-clean: ## Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	@docker-compose down -v --remove-orphans
	@docker system prune -f

# Dependencies
deps: ## Install all dependencies
	@echo "📦 Installing dependencies..."
	@cd backend && go mod download
	@cd frontend && npm install

deps-backend: ## Install backend dependencies
	@echo "📦 Installing backend dependencies..."
	@cd backend && go mod download

deps-frontend: ## Install frontend dependencies
	@echo "📦 Installing frontend dependencies..."
	@cd frontend && npm install

# Security
security-check: ## Run security checks
	@echo "🔒 Running security checks..."
	@cd backend && go mod verify
	@cd frontend && npm audit
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image codewithdell/backend:latest
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image codewithdell/frontend:latest

# Performance
benchmark: ## Run performance benchmarks
	@echo "⚡ Running benchmarks..."
	@cd backend && go test -bench=. ./...

profile: ## Run profiling
	@echo "📊 Running profiling..."
	@cd backend && go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...

# Documentation
docs: ## Generate documentation
	@echo "📚 Generating documentation..."
	@cd backend && swag init -g main.go
	@cd frontend && npm run docs

docs-serve: ## Serve documentation locally
	@echo "📚 Serving documentation..."
	@cd docs && python3 -m http.server 8000

# Monitoring
monitor: ## Start monitoring tools
	@echo "📊 Starting monitoring tools..."
	@docker-compose -f docker-compose.monitoring.yml up -d

monitor-stop: ## Stop monitoring tools
	@echo "📊 Stopping monitoring tools..."
	@docker-compose -f docker-compose.monitoring.yml down

# Production
prod-build: ## Build for production
	@echo "🏭 Building for production..."
	@docker build -f backend/Dockerfile.prod -t codewithdell/backend:latest ./backend
	@docker build -f frontend/Dockerfile.prod -t codewithdell/frontend:latest ./frontend

prod-deploy: ## Deploy to production
	@echo "🚀 Deploying to production..."
	@docker-compose -f docker-compose.prod.yml up -d

prod-logs: ## Show production logs
	@echo "📝 Showing production logs..."
	@docker-compose -f docker-compose.prod.yml logs -f

# Utilities
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf backend/bin
	@rm -rf frontend/.next
	@rm -rf frontend/out
	@rm -rf node_modules
	@go clean -cache -testcache

check: ## Run all checks (lint, test, security)
	@echo "✅ Running all checks..."
	@make lint
	@make test
	@make security-check

ci: ## Run CI pipeline
	@echo "🔄 Running CI pipeline..."
	@make deps
	@make lint
	@make test
	@make security-check
	@make build

# Database utilities
db-shell: ## Open database shell
	@echo "🗄️  Opening database shell..."
	@docker exec -it codewithdell_postgres psql -U codewithdell -d codewithdell

redis-cli: ## Open Redis CLI
	@echo "🔴 Opening Redis CLI..."
	@docker exec -it codewithdell_redis redis-cli

# Development tools
pgadmin: ## Start pgAdmin4 (Database GUI)
	@echo "🖥️  Starting pgAdmin4..."
	@docker-compose --profile tools up -d pgadmin

adminer: ## Start Adminer (Database GUI)
	@echo "🖥️  Starting Adminer..."
	@docker-compose --profile tools up -d adminer

redis-commander: ## Start Redis Commander (Redis GUI)
	@echo "🖥️  Starting Redis Commander..."
	@docker-compose --profile tools up -d redis-commander

# Git utilities
git-hooks: ## Install Git hooks
	@echo "🔗 Installing Git hooks..."
	@cd frontend && npx husky install

commit: ## Make a conventional commit
	@echo "💬 Making conventional commit..."
	@cd frontend && npx git-cz

# Environment
env-example: ## Create example environment files
	@echo "📝 Creating example environment files..."
	@cp .env.example .env
	@cp frontend/.env.local.example frontend/.env.local
	@echo "✅ Environment files created. Please edit them with your configuration."

# Health checks
health: ## Check application health
	@echo "🏥 Checking application health..."
	@curl -f http://localhost:8080/health || echo "❌ Backend health check failed"
	@curl -f http://localhost:3000/api/health || echo "❌ Frontend health check failed"

# Quick development
quick: ## Quick development setup (minimal)
	@echo "⚡ Quick development setup..."
	@docker-compose up -d postgres redis
	@sleep 5
	@cd backend && go run main.go &
	@cd frontend && npm run dev &
	@echo "✅ Quick setup complete!"
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🔧 Backend: http://localhost:8080"

# Development workflow
workflow: ## Complete development workflow
	@echo "🔄 Running complete development workflow..."
	@make setup
	@make deps
	@make migrate
	@make seed
	@make dev

# Help for specific commands
help-setup: ## Show setup help
	@echo "Setup Commands:"
	@echo "  make setup      - Initial development environment setup"
	@echo "  make dev        - Start development servers"
	@echo "  make stop       - Stop all services"
	@echo "  make status     - Show service status"

help-db: ## Show database help
	@echo "Database Commands:"
	@echo "  make migrate    - Run database migrations"
	@echo "  make seed       - Seed database with sample data"
	@echo "  make reset-db   - Reset database (WARNING: deletes all data)"
	@echo "  make backup     - Create database backup"
	@echo "  make restore    - Restore database from backup"

help-docker: ## Show Docker help
	@echo "Docker Commands:"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker Compose services"
	@echo "  make docker-logs  - Show Docker logs"
	@echo "  make docker-clean - Clean Docker resources"

# Catch-all for unknown targets
%:
	@echo "❌ Unknown target: $@"
	@echo "💡 Run 'make help' to see available commands" 