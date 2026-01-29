include ./.env

MIGRATION_PATH=./db/migration/
SEEDER_PATH=./db/seeder/seeder.sql
DB_URL=postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down:
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) down

seeder:
	psql -q -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USERNAME) -d $(DB_NAME) -f $(SEEDER_PATH)
	