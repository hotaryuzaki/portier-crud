package main

import (
	"log"
	"portier/internal/config"
	"portier/internal/delivery/http"
	"portier/pkg/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
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
