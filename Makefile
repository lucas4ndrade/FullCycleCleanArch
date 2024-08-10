run-infra:
	docker-compose up

run:
	go run ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go 

# Definir valores padrão
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_USER ?= root
DB_PASS ?= root
MIGRATION_UP_FILE ?= ./cmd/migrations/migration_up.sql
MIGRATION_DOWN_FILE ?= ./cmd/migrations/migration_down.sql

# Comando migration-up
migration-up:
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) < $(MIGRATION_UP_FILE)

migration-down:
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) < $(MIGRATION_DOWN_FILE)
