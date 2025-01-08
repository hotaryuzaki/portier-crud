# Load the environment variables from the .env file
include .env

# Default target to run the application
run:
	go run cmd/app/main.go

# Target to run with nodemon for live reload
dev:
	nodemon --exec "go run cmd/app/main.go" --watch . --ext go

# Target to build the application
build:
	go build -o app cmd/app/main.go

# Target to clean up build artifacts
clean:
	rm -f app

# Target to run tests
test:
	go test ./...

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
