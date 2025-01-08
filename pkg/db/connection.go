package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConnectPostgres(dsn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return conn
}
