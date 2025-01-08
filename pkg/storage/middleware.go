package storage

import (
	"portier/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/postgres/v3"
)

func SetupPostgresMiddleware(app *fiber.App, cfg config.Config) {
	storage := postgres.New(postgres.Config{
		ConnectionURI: cfg.PostgresDSN,
		Table:         "cache",
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("storage", storage)
		return c.Next()
	})
}
