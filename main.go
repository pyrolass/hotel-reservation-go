package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pyrolass/hotel-reservation-go/handlers"
)

func main() {

	app := fiber.New()

	api := app.Group("api/v1")

	port := ":3000"

	api.Get(
		"/ping",
		func(c *fiber.Ctx) error {
			return c.JSON(
				map[string]any{
					"message": "pong",
				},
			)
		},
	)

	api.Get("/user", handlers.HandleGetUsers)

	api.Get("/user/:id", handlers.HandleGetUser)

	app.Listen(port)
}
