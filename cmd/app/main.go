package main

import (
	"log"
	"portier/internal/config"
	"portier/internal/delivery/http"
	"portier/pkg/storage"

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

	// Start the server
	log.Fatal(app.Listen(cfg.ServerPort))
}
