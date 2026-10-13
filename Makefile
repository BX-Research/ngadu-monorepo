.DEFAULT_GOAL := help

BUN ?= bun
GO ?= go
DOCKER ?= docker

.PHONY: help install install-frontend install-backend hooks dev dev-frontend dev-backend db-up db-stop db-down db-logs db-migrate db-status db-revision db-seed seed lint typecheck build start test-backend vet-backend test-dev check

help:
	@printf '%s\n' \
	  'Local development commands (no production database targets)' \
	  'make install          Download frontend/Go dependencies and enable Git hooks' \
	  'make dev              Start Docker PostgreSQL, Go backend, and frontend' \
	  'make seed / db-seed   Start, migrate, and seed local PostgreSQL' \
	  'make db-migrate       Apply local database migrations using the Go Goose library' \
	  'make db-status        Show local migration status' \
	  'make db-revision msg=NAME  Create a SQL migration using the Go Goose library' \
	  'make db-up            Start Docker PostgreSQL and wait until healthy' \
	  'make db-stop          Stop PostgreSQL, preserving its data' \
	  'make db-logs          Follow PostgreSQL logs' \
	  'make lint             Run frontend lint' \
	  'make typecheck        Run frontend TypeScript checks' \
	  'make build            Build the frontend for production' \
	  'make start            Start the built frontend' \
	  'make test-dev         Run development supervisor tests' \
	  'make test-backend     Run Go tests (DB tests need TEST_DATABASE_URL)' \
	  'make vet-backend      Run Go static analysis' \
	  'make check            Run supervisor tests, frontend checks, and Go tests/vet' \
	  'make hooks            Enable shared Git hooks in this clone'

install: install-frontend install-backend hooks

install-frontend:
	$(BUN) run setup:frontend

install-backend:
	$(GO) -C backend mod download

hooks:
	$(BUN) run setup:hooks

dev:
	BUN_BIN="$(BUN)" GO_BIN="$(GO)" DOCKER_BIN="$(DOCKER)" $(BUN) backend/scripts/dev.mjs dev

dev-frontend:
	$(BUN) run dev

dev-backend:
	$(GO) -C backend run ./cmd/api -dev

db-up:
	$(DOCKER) compose up -d --wait --wait-timeout 60 db

db-stop:
	$(DOCKER) compose stop db

db-down: db-stop

db-logs:
	$(DOCKER) compose logs --follow db

db-migrate: db-up
	$(GO) -C backend run ./scripts/migrate up

db-status: db-up
	$(GO) -C backend run ./scripts/migrate status

# Pass the name through the environment rather than interpolating shell code.
export msg
db-revision:
	@test -n "$$msg" || { echo 'Usage: make db-revision msg=NAME'; exit 1; }
	$(GO) -C backend run ./scripts/migrate create "$$msg"

db-seed: seed

seed: db-migrate
	$(GO) -C backend run ./scripts/seed

lint:
	$(BUN) run lint

typecheck:
	$(BUN) run typecheck

build:
	$(BUN) run build

start:
	$(BUN) run start

test-dev:
	$(BUN) test backend/scripts/dev.test.mjs

test-backend:
	$(GO) -C backend test -v ./...

vet-backend:
	$(GO) -C backend vet ./...

# Keep checks sequential, even with make -j: build and typecheck share generated files.
check:
	$(MAKE) test-dev
	$(MAKE) lint
	$(MAKE) typecheck
	$(MAKE) build
	$(MAKE) test-backend
	$(MAKE) vet-backend
