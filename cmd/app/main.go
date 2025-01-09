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
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func requestLogger(c *fiber.Ctx) error {
	log.Printf("Received request: %s %s", c.Method(), c.Path())
	return c.Next()
}

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

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins for testing purposes
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Use the request logger middleware
	app.Use(requestLogger)

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
