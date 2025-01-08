package http

import (
	"log"
	"portier/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Add a route for the root ("/") that returns "Hello, world"
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world")
	})

	// USER
	app.Get("/users", getUsers)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)
}

func getUsers(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl http://localhost:4000/users

	users, err := service.GetAllUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func createUser(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X POST http://localhost:4000/users ^
	// -H "Content-Type: application/json" ^
	// -d "{\"username\": \"johndoe\", \"email\": \"johndoe@example.com\", \"password\": \"securepassword123\", \"name\": \"John Doe\", \"gender\": true, \"id_number\": \"123456789\", \"user_image\": \"http://example.com/image.jpg\", \"tenant_id\": 1}"

	var user service.User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	createdUser, err := service.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func updateUser(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X PUT http://localhost:4000/users/1 ^
	// -H "Content-Type: application/json" ^
	// -d "{\"username\": \"AMRI\", \"email\": \"amri@example.com\", \"password\": \"12345\", \"name\": \"AMRI\", \"gender\": true, \"id_number\": \"317172727272\", \"user_image\": \"http://example.com/image.jpg\", \"tenant_id\": 1}"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	var user service.User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	updatedUser, err := service.UpdateUser(id, user)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Explicitly returning status code for successful update
	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func deleteUser(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X DELETE http://localhost:4000/users/3 ^
	// -H "Content-Type: application/json"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	err = service.DeleteUser(id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}
