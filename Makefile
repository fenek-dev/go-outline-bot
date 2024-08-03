
# --- Database ---
db-up:
	docker compose -f=docker/docker-compose.db.yml -p=pg_bot up

db-down:
	docker compose -f=docker/docker-compose.db.yml -p=pg_bot down

# --- Migrations ---

DB_USER ?= user
DB_PASSWORD ?= password
DB_ADDRESS ?= localhost:5432
DB_NAME ?= bot
# example: make db-migrate DB_USER=your_user DB_PASSWORD=your_password DB_ADDRESS=your_address DB_NAME=your_dbname

db-migrate:
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_ADDRESS)/$(DB_NAME)?sslmode=disable" -path migrations up

db-migrate-down:
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_ADDRESS)/$(DB_NAME)?sslmode=disable" -path migrations down