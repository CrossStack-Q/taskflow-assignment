# ==============================
# LOAD ENV
# ==============================
include .env

MIGRATIONS_PATH=./internal/database/migrations

# ==============================
# APP
# ==============================

.PHONY: run
run:
	go run ./backend/cmd/api

.PHONY: dev
dev:
	air

.PHONY: build
build:
	go build -o bin/app ./backend/cmd/api

# ==============================
# DOCKER (DB + REDIS ONLY)
# ==============================

.PHONY: up
up:
	docker compose up -d postgres redis

.PHONY: down
down:
	docker compose down

.PHONY: logs
logs:
	docker compose logs -f

# ==============================
# MIGRATIONS
# ==============================

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide name: make migrate-create name=your_migration"; \
	else \
		migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name); \
	fi

.PHONY: migrate-up
migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database "$(TASKFLOW_DATABASE_ADDR)" up

.PHONY: migrate-down
migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database "$(TASKFLOW_DATABASE_ADDR)" down 1

.PHONY: migrate-force
migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "❌ Please provide version: make migrate-force version=1"; \
	else \
		migrate -path=$(MIGRATIONS_PATH) -database "$(TASKFLOW_DATABASE_ADDR)" force $(version); \
	fi

.PHONY: migrate-version
migrate-version:
	migrate -path=$(MIGRATIONS_PATH) -database "$(TASKFLOW_DATABASE_ADDR)" version

.PHONY: migrate-reset
migrate-reset:
	migrate -path=$(MIGRATIONS_PATH) -database "$(TASKFLOW_DATABASE_ADDR)" down
	migrate -path=$(MIGRATIONS_PATH) -database "$(TASKFLOW_DATABASE_ADDR)" up

# ==============================
# SEED
# ==============================

.PHONY: seed
seed:
	go run ./backend/cmd/seed/main.go

# ==============================
# TEST
# ==============================

.PHONY: test
test:
	go test -v ./...

# ==============================
# HELP
# ==============================

.PHONY: help
help:
	@echo "Available commands:"
	@echo " make up                → start postgres + redis"
	@echo " make down              → stop all services"
	@echo " make run               → run API locally"
	@echo " make migrate-create    → create migration (name=xyz)"
	@echo " make migrate-up        → apply migrations"
	@echo " make migrate-down      → rollback last migration"
	@echo " make migrate-reset     → reset DB"
	@echo " make migrate-version   → show current version"
	@echo " make seed              → seed data"
	@echo " make test              → run tests"