package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new Fiber app
	app := fiber.New()

	// Define a simple GET route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Go Fiber!")
	})

	// Start the server on port 3000
	app.Listen(":3000")
}
