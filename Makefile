.PHONY: help install test lint lint-fix dev build docker-up docker-down docker-logs clean backend-% frontend-%

help:
	@echo "specialstandard - Therapy Management Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  make install        - Install all dependencies locally"
	@echo "  make dev            - Run with Docker (hot reload enabled)"
	@echo "  make test           - Run all tests"
	@echo "  make lint           - Check linting (frontend + backend)"
	@echo "  make lint-fix       - Fix auto-fixable linting issues"
	@echo "  make build          - Build Docker images"
	@echo "  make docker-up      - Start services without watch"
	@echo "  make docker-down    - Stop all Docker services"
	@echo "  make docker-logs    - View Docker logs"
	@echo "  make clean          - Clean all artifacts and stop Docker"
	@echo ""
	@echo "Backend-specific:"
	@echo "  make backend-test, backend-test-coverage, etc."
	@echo ""
	@echo "Frontend-specific:"
	@echo "  make frontend-dev, frontend-build, etc."

# Install all dependencies (for local development)
install:
	$(MAKE) -C backend install
	$(MAKE) -C frontend install

# Run all tests
test:
	$(MAKE) -C backend test

# Lint all code
lint:
	@echo "Linting backend..."
	$(MAKE) -C backend lint
	@echo "Linting frontend..."
	$(MAKE) -C frontend lint

# Fix auto-fixable linting issues
lint-fix:
	@echo "Fixing backend linting..."
	$(MAKE) -C backend lint-fix
	@echo "Fixing frontend linting..."
	$(MAKE) -C frontend lint-fix

# Run with Docker and hot reload (main development command)
dev:
	docker compose up --build --watch

# Build Docker images
build:
	docker compose build

# Start services without watch mode
docker-up:
	docker compose up -d

# Stop all services
docker-down:
	docker compose down

# View logs
docker-logs:
	docker compose logs -f

# Restart services
docker-restart:
	docker compose restart

# Clean everything
clean:
	docker compose down -v
	$(MAKE) -C backend clean
	$(MAKE) -C frontend clean

# Delegate backend commands (for local development)
backend-%:
	$(MAKE) -C backend $*

# Delegate frontend commands (for local development)
frontend-%:
	$(MAKE) -C frontend $*