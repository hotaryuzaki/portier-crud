# Load the environment variables from the .env file
include .env

run:
	go run cmd/app/main.go

migrate-up:
ifeq ($(OS),Windows_NT)
	# Use Windows-compatible syntax for setting PGPASSWORD
	@for file in db/migrations/*.sql; do \
    set PGPASSWORD=$(DB_PASS) && \
    psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f "$$file"; \
	done
else
	# Use Unix-compatible syntax for setting PGPASSWORD
	@for file in db/migrations/*.sql; do \
		PGPASSWORD=$(PGPASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f "$$file"; \
	done
endif
