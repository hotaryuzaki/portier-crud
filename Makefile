run:
	go run cmd/app/main.go

migrate-up:
	psql -d portier -f db/migrations/*.sql
