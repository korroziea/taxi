CURRENT_DIR=$(shell pwd)

DATABASE_STRING=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable

migrate-create: ## Create migration file  `make migrate-create name=migration_name`
	@migrate create -dir migrations -seq -ext sql $(name)

migrate-up: ## Run migrations up `make migrate-up`
	@migrate -database $(DATABASE_STRING) -path migrations -verbose up

migrate-down: ## Run one migration down `make migrate-down`
	@migrate -database $(DATABASE_STRING) -path migrations -verbose down 1

migrate-force: ## Change migration version without running migrations `make migrate-force n=14`
	@migrate -database $(DATABASE_STRING) -path migrations -verbose force $(n)

migrate-version: ## Print current migration version `make migrate-version`
	@migrate -database $(DATABASE_STRING) -path migrations version
