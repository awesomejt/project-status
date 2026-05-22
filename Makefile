.PHONY: help validate lint format format-check typecheck build build-cli test test-api test-web test-cli smoke integration-test clean clean-all migrations

## Show this help message
help:
	@echo "Project Status - Makefile Targets"
	@echo "================================"
	@echo ""
	@echo "Validation:"
	@echo "  validate      - Run all validation (lint, format-check, typecheck, build, test)"
	@echo "  lint          - Run linters for all modules"
	@echo "  format        - Format all code"
	@echo "  format-check  - Check formatting without changes"
	@echo "  typecheck     - Run type checkers for all modules"
	@echo ""
	@echo "Build:"
	@echo "  build         - Build all modules (web, cli)"
	@echo "  build-cli     - Build CLI to build/project-status"
	@echo ""
	@echo "Tests:"
	@echo "  test          - Run all tests"
	@echo "  test-api      - Run API tests"
	@echo "  test-web      - Run web tests"
	@echo "  test-cli      - Run CLI tests"
	@echo "  smoke         - Run smoke tests"
	@echo "  integration-test - Run integration tests"
	@echo ""
	@echo "Development:"
	@echo "  dev           - Start all development services"
	@echo "  db            - Start PostgreSQL database only"
	@echo "  migrations    - Run database migrations"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean         - Clean build artifacts"
	@echo "  clean-all     - Clean all artifacts including Docker volumes"
	@echo ""

## Run all validation
validate: lint format-check typecheck build test

## Run linters for all modules
lint: lint-api lint-web

lint-api:
	cd api && uv run ruff check .

lint-web:
	cd web && npm run lint

## Format all code
format: format-api format-web

format-api:
	cd api && uv run ruff format .

format-web:
	@echo "Web formatting: no specific formatter configured"

## Check formatting without changes
format-check: format-check-api format-check-web

format-check-api:
	cd api && uv run ruff format --check .

format-check-web:
	@echo "Web format check: no specific formatter configured"

## Run type checkers for all modules
typecheck: typecheck-api typecheck-web

typecheck-api:
	cd api && uv run mypy . || true

typecheck-web:
	cd web && npm run typecheck

## Build all modules
build: build-web build-cli

## Build web module
build-web:
	cd web && npm run build

## Build CLI to build/project-status
build-cli:
	mkdir -p build
	cd cli && go build -o ../build/project-status .

## Run all tests
test: test-api test-web test-cli

## Run API tests
test-api:
	cd api && uv run pytest

## Run web tests
test-web:
	cd web && npm test

## Run CLI tests
test-cli:
	cd cli && go test ./...

## Run smoke tests
smoke:
	./scripts/smoke-curl.sh

## Run integration tests
integration-test:
	docker compose run --rm integration-test

## Start all development services
dev:
	docker compose up --build

## Start PostgreSQL database only
db:
	docker compose up -d db

## Run database migrations
migrations:
	docker compose up migrations

## Clean build artifacts
clean: clean-api clean-web clean-cli

clean-api:
	cd api && rm -rf __pycache__ .mypy_cache .ruff_cache .pytest_cache .coverage htmlcov

clean-web:
	cd web && rm -rf dist .vite

clean-cli:
	rm -rf cli/build

## Clean all artifacts including Docker volumes
clean-all: clean
	docker compose down -v
