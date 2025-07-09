.PHONY: help up down logs ps migrate seed \
        up-api down-api logs-api \
        up-db down-db logs-db \
		 up-sync down-sync logs-sync \
        up-redis down-redis logs-redis \
        up-web down-web logs-web

# 既存の help に追記
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  up           Start the development environment"
	@echo "  down         Stop the development environment"
	@echo "  logs         Show logs of all services"
	@echo "  ps           Show container status"
	@echo "  migrate      Run DB migration"
	@echo "  seed         Insert seed data"
	@echo ""
	@echo "  up-api       Start API only"
	@echo "  down-api     Stop API only"
	@echo "  logs-api     Show logs for API"
	@echo ""
	@echo "  up-db        Start DB only"
	@echo "  down-db      Stop DB only"
	@echo "  logs-db      Show logs for DB"
	@echo ""
	@echo "  up-sync    Start sync only"
	@echo "  down-sync  Stop sync only"
	@echo "  logs-sync  Show logs for sync"
	@echo ""
	@echo "  up-redis     Start Redis only"
	@echo "  down-redis   Stop Redis only"
	@echo "  logs-redis   Show logs for Redis"
	@echo ""
	@echo "  up-web       Start Web only"
	@echo "  down-web     Stop Web only"
	@echo "  logs-web     Show logs for Web"

up:
	docker compose -f ops/docker-compose.dev.yml up --build -d

down:
	docker compose -f ops/docker-compose.dev.yml down

logs:
	docker compose -f ops/docker-compose.dev.yml logs -f

ps:
	docker compose -f ops/docker-compose.dev.yml ps

migrate:
	docker exec -i semantra-db psql -U dev -d semantra < db/init/01_create_tables.sql

seed:
	docker exec -i semantra-db psql -U dev -d semantra < db/init/02_insert_sample_data.sql

# API only
up-api:
	docker compose -f ops/docker-compose.dev.yml up --build -d api

down-api:
	docker compose -f ops/docker-compose.dev.yml stop api

logs-api:
	docker compose -f ops/docker-compose.dev.yml logs -f api

# DB only
up-db:
	docker compose -f ops/docker-compose.dev.yml up --build -d db

down-db:
	docker compose -f ops/docker-compose.dev.yml stop db

logs-db:
	docker compose -f ops/docker-compose.dev.yml logs -f db

# sync only
up-sync:
	docker compose -f ops/docker-compose.dev.yml up --build -d sync

down-sync:
	docker compose -f ops/docker-compose.dev.yml stop sync

logs-sync:
	docker compose -f ops/docker-compose.dev.yml logs -f sync

# Redis only
up-redis:
	docker compose -f ops/docker-compose.dev.yml up --build -d redis

down-redis:
	docker compose -f ops/docker-compose.dev.yml stop redis

logs-redis:
	docker compose -f ops/docker-compose.dev.yml logs -f redis

# Web only
up-web:
	docker compose -f ops/docker-compose.dev.yml up --build -d web

down-web:
	docker compose -f ops/docker-compose.dev.yml stop web

logs-web:
	docker compose -f ops/docker-compose.dev.yml logs -f web