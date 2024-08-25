ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	air -- --config=.env

# --- Database ---
db-up:
	docker compose -f=docker/docker-compose.db.yml -p=pg_bot up -d

db-down:
	docker compose -f=docker/docker-compose.db.yml -p=pg_bot down

# --- App ---
app-up:
	make db-up
	docker-compose up -d --build

app-down:
	docker-compose down

app-restart:
	make app-down
	make app-up

# --- Migrations ---

DB_HOST=localhost

db-migrate:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):${DB_PORT}/$(DB_NAME)?sslmode=disable" -path migrations up

db-migrate-down:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):${DB_PORT}/$(DB_NAME)?sslmode=disable" -path migrations down