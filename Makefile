-include .env

.SILENT:

DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

tidy:
	@go mod tidy
	@go mod vendor

migration:
	@migrate create -ext sql -dir ./migrations -seq $(name)

migrateup:
	@migrate --path ./migrations --database "$(DB_URL)" --verbose up

migratedown:
	@migrate --path ./migrations --database "$(DB_URL)" --verbose down