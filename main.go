package main

import (
	"log"
	"github.com/adramelech-123/fiber-api/database"
	"github.com/adramelech-123/fiber-api/routes"
	"github.com/gofiber/fiber/v2"
)

// Ctx represents the Context which hold the HTTP request and response.
// It has methods for the request query string, parameters, body, HTTP headers and so on.
func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to this awesome Go based API")
}

func setupRoutes(app*fiber.App) {
	//Welcome endpoint
	app.Get("/api", welcome)

	// User endpoints
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetUsers)
	app.Get("/api/users/:id", routes.GetUser)
	app.Put("/api/users/:id", routes.UpdateUser)
}

func main() {
	database.ConnectDb()
	app := fiber.New()

	setupRoutes(app) 
	log.Fatal(app.Listen(":3000"))
}