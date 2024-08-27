run-infra:
	docker-compose up

run:
	go run ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go 

DB_HOST ?= localhost
DB_PORT ?= 3306
DB_USER ?= root
DB_PASS ?= root
DB_NAME ?= fullcycle

migration-up:
	migrate -path=cmd/migrations -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" -verbose up

migration-down:
	migrate -path=cmd/migrations -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" -verbose down
