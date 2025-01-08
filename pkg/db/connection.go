package db

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
)

var (
	conn *pgx.Conn
	once sync.Once
)

// ConnectPostgres establishes a connection to the PostgreSQL database
func ConnectPostgres(dsn string) *pgx.Conn {
	once.Do(func() {
		var err error
		log.Printf("Connecting to PostgreSQL with DSN: %v", dsn) // Debug log for DSN
		conn, err = pgx.Connect(context.Background(), dsn)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		log.Println("Successfully connected to PostgreSQL!")
	})
	return conn
}

// GetConnection returns the established database connection
func GetConnection() *pgx.Conn {
	if conn == nil {
		log.Fatal("Database connection is not established yet")
	}
	return conn
}

// Close the database connection
func Close() {
	if conn != nil {
		err := conn.Close(context.Background())
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}
}
