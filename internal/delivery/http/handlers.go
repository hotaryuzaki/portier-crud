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

	// COPIES routes
	app.Get("/copies", getCopies)
	app.Get("/copies/:id", getCopiesById)
	app.Post("/copies", createCopy)
	app.Put("/copies/:id", updateCopy)
	app.Delete("/copies/:id", deleteCopy)

	// TENANT routes
	app.Get("/tenants", getTenants)
	app.Get("/tenants/:id", getTenantById)
	app.Post("/tenants", createTenant)
	app.Put("/tenants/:id", updateTenant)
	app.Delete("/tenants/:id", deleteTenant)
}

/*** USERS HANDLERS ***/

func getUsers(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl "http://localhost:4000/users?limit=10&offset=0"

	// Parse limit and offset from query parameters
	limitStr := c.Query("limit", "10")  // Default limit is 10
	offsetStr := c.Query("offset", "0") // Default offset is 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'limit' parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'offset' parameter",
		})
	}

	// Call the service to get paginated users
	response, err := service.GetAllUsers(limit, offset)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Convert Gender boolean to string for the response
	for i := range response.Users {
		response.Users[i].GenderStr = response.Users[i].ConvertGenderToStr()
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
	// -d "{\"username\": \"johndoe\", \"email\": \"johndoe@example.com\", \"password\": \"securepassword123\", \"name\": \"John Doe\", \"gender\": \"1\", \"id_number\": \"123456789\", \"user_image\": \"http://example.com/image.jpg\", \"tenant_id\": 1}"

	var user service.User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	// Convert the GenderStr to a boolean
	if err := user.ConvertGender(); err != nil {
		log.Printf("Error converting gender: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid gender value")
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
	// curl "http://localhost:4000/keys?limit=10&offset=0"

	// Parse limit and offset from query parameters
	limitStr := c.Query("limit", "10")  // Default limit is 10
	offsetStr := c.Query("offset", "0") // Default offset is 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'limit' parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'offset' parameter",
		})
	}

	// Call the service to get paginated keys
	keys, err := service.GetAllKeys(limit, offset)
	if err != nil {
		log.Printf("Error getting keys: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch keys",
		})
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
	// REQUEST EXAMPLE
	// curl -X POST http://localhost:4000/keys ^
	// -H "Content-Type: application/json" ^
	// -d "{\"name\": \"TEST Key\"}"

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
	// -d "{\"name\": \"Updated Key Name\", \"is_active\": false}"

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

/*** COPIES HANDLERS ***/

func getCopies(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl "http://localhost:4000/copies?limit=10&offset=0"

	// Parse limit and offset from query parameters
	limitStr := c.Query("limit", "10")  // Default limit is 10
	offsetStr := c.Query("offset", "0") // Default offset is 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'limit' parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'offset' parameter",
		})
	}

	// Call the service to get paginated copies
	copies, err := service.GetAllCopies(limit, offset)
	if err != nil {
		log.Printf("Error getting copies: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch copies",
		})
	}

	return c.Status(fiber.StatusOK).JSON(copies)
}

func getCopiesById(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl http://localhost:4000/copies/1

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error getting id: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	copy, err := service.GetCopyByID(id)
	if err != nil {
		log.Printf("Error getting copy: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(copy)
}

func createCopy(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X POST http://localhost:4000/copies ^
	// -H "Content-Type: application/json" ^
	// -d "{\"name\": \"TEST Copy\"}"

	var copy service.Copy
	if err := c.BodyParser(&copy); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	createdCopy, err := service.CreateCopy(copy)
	if err != nil {
		log.Printf("Error creating copy: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusCreated).JSON(createdCopy)
}

func updateCopy(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X PUT http://localhost:4000/copies/1 ^
	// -H "Content-Type: application/json" ^
	// -d "{\"name\": \"Updated Copy Name\", \"is_active\": true}"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	var copy service.Copy
	if err := c.BodyParser(&copy); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	updatedCopy, err := service.UpdateCopy(id, copy)
	if err != nil {
		log.Printf("Error updating copy: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(updatedCopy)
}

func deleteCopy(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X DELETE http://localhost:4000/copies/3 ^
	// -H "Content-Type: application/json"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	err = service.DeleteCopy(id)
	if err != nil {
		log.Printf("Error deleting copy: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}

/*** TENANTS HANDLERS ***/

func getTenants(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl "http://localhost:4000/tenants?limit=10&offset=0"

	// Parse limit and offset from query parameters
	limitStr := c.Query("limit", "10")  // Default limit is 10
	offsetStr := c.Query("offset", "0") // Default offset is 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'limit' parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'offset' parameter",
		})
	}

	// Call the service to get paginated tenants
	tenants, err := service.GetAllTenants(limit, offset)
	if err != nil {
		log.Printf("Error getting tenants: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch copies",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tenants)
}

func getTenantById(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl http://localhost:4000/tenants/1

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	tenant, err := service.GetTenantByID(id)
	if err != nil {
		log.Printf("Error getting tenant: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(tenant)
}

func createTenant(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X POST http://localhost:4000/tenants ^
	// -H "Content-Type: application/json" ^
	// -d "{\"name\": \"PT ZIG ZAG\", \"address\": \"Jln banyak belok\", \"status\": \"Pending\"}"

	var tenant service.Tenant
	if err := c.BodyParser(&tenant); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	createdTenant, err := service.CreateTenant(tenant)
	if err != nil {
		log.Printf("Error creating tenant: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusCreated).JSON(createdTenant)
}

func updateTenant(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X PUT http://localhost:4000/tenants/1 ^
	// -H "Content-Type: application/json" ^
	// -d "{\"name\": \"Updated Name\", \"address\": \"Updated Address\", \"status\": \"active\"}"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	var tenant service.Tenant
	if err := c.BodyParser(&tenant); err != nil {
		log.Printf("Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	updatedTenant, err := service.UpdateTenant(id, tenant)
	if err != nil {
		log.Printf("Error updating tenant: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(updatedTenant)
}

func deleteTenant(c *fiber.Ctx) error {
	// REQUEST EXAMPLE
	// curl -X DELETE http://localhost:4000/tenants/3 ^
	// -H "Content-Type: application/json"

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	err = service.DeleteTenant(id)
	if err != nil {
		log.Printf("Error deleting tenant: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusNoContent).SendString("Tenant deleted successfully")
}
