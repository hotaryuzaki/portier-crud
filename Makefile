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

# Command to run Docker Compose
docker-compose-up:
	docker-compose up -d

# Command to build Docker images
docker-build:
	docker-compose build

# Command to stop Docker Compose
docker-compose-down:
	docker-compose down

# Command to restart Docker Compose
docker-compose-restart:
	docker-compose down && docker-compose up -d

# Command to view Docker Compose logs
docker-compose-logs:
	docker-compose logs -f

# Command to view logs of the app container
docker-compose-logs-app:
	docker-compose logs -f app

# Command to view logs of the frontend container
docker-compose-logs-frontend:
	docker-compose logs -f frontend

# Command to view logs of the db container
docker-compose-logs-db:
	docker-compose logs -f db

# Command to remove Docker containers, networks, and volumes
docker-clean:
	docker-compose down -v --rmi all --remove-orphans

# Command to build the backend
build-backend:
	docker-compose run --rm app go build -o main .

# Command to build the frontend
build-frontend:
	docker-compose run --rm frontend yarn build

# Command to run migrations
migrate-up:
	for file in db/migrations/*.sql; do \
		docker-compose run --rm db bash -c "PGPASSWORD=${DB_PASS} psql -h db -U ${DB_USER} -d ${DB_NAME} -f /docker-entrypoint-initdb.d/migrations/$$(basename $$file)"; \
	done

# Command to open psql session in the PostgreSQL container
psql:
	docker-compose run --rm db psql -h db -U ${DB_USER} -d ${DB_NAME}
