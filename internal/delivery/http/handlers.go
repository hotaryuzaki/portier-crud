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

	// USER routes
	app.Get("/users", getUsers)
	app.Get("/users/:id", getUsersById)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	// KEYS routes
	app.Get("/keys", getKeys)
	app.Get("/keys/:id", getKeysById)
	app.Post("/keys", createKey)
	app.Put("/keys/:id", updateKey)
	app.Delete("/keys/:id", deleteKey)
}

/*** USERS HANDLERS ***/

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

func getUsersById(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl http://localhost:4000/users/1

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error getting id: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	user, err := service.GetUserByID(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(user)
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

/*** KEYS HANDLERS ***/

func getKeys(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl http://localhost:4000/keys

	keys, err := service.GetAllKeys()
	if err != nil {
		log.Printf("Error getting keys: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(keys)
}

func getKeysById(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl http://localhost:4000/keys/1

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error getting id: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	user, err := service.GetKeysByID(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func createKey(c *fiber.Ctx) error {
	// REQUEST EXAMPLE// curl -X POST http://localhost:4000/keys ^
	// -H "Content-Type: application/json" ^
	// -d "{\"name\": \"TEST Key\", \"is_active\": true, \"created_by\": 1}"

	var key service.Key
	if err := c.BodyParser(&key); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	createdKey, err := service.CreateKey(key)
	if err != nil {
		log.Printf("Error creating key: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusCreated).JSON(createdKey)
}

func updateKey(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X PUT http://localhost:4000/keys/1 ^
	// -H "Content-Type: application/json" ^
	// -d "{\"name\": \"Updated Key Name\", \"is_active\": false, \"created_by\": 1}"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	var key service.Key
	if err := c.BodyParser(&key); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	updatedKey, err := service.UpdateKey(id, key)
	if err != nil {
		log.Printf("Error updating key: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(updatedKey)
}

func deleteKey(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X DELETE http://localhost:4000/keys/3 ^
	// -H "Content-Type: application/json"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	err = service.DeleteKey(id)
	if err != nil {
		log.Printf("Error deleting key: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}
