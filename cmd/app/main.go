package main

import (
	"log"
	"os"
	"os/signal"
	"portier/internal/config"
	"portier/internal/delivery/http"
	"portier/pkg/db"
	"portier/pkg/storage"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Setup Fiber app
	app := fiber.New()

	// Setup PostgreSQL middleware
	storage.SetupPostgresMiddleware(app, cfg)

	// Register routes
	http.RegisterRoutes(app)

	// Start the server in a goroutine
	go func() {
		if err := app.Listen(cfg.ServerPort); err != nil {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Graceful shutdown logic: handle interrupt signal (e.g., CTRL+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received

	// Perform cleanup tasks before shutting down
	log.Println("Shutting down gracefully...")
	db.Close() // Close the PostgreSQL connection
	log.Println("Database connection closed.")
	log.Println("Server stopped.")
}
