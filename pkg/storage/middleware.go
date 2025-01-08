package storage

import (
	"portier/internal/config"
	"portier/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/postgres/v3"
)

// SetupPostgresMiddleware initializes both the cache and database connections
func SetupPostgresMiddleware(app *fiber.App, cfg config.Config) {
	// Setup the PostgreSQL cache storage
	storage := postgres.New(postgres.Config{
		ConnectionURI: cfg.PostgresDSN,
		Table:         "cache",
	})

	// Initialize general PostgreSQL connection (for CRUD operations)
	db.ConnectPostgres(cfg.PostgresDSN)

	// Attach the cache storage to the context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("storage", storage) // Store the cache storage in the context
		return c.Next()
	})
}
