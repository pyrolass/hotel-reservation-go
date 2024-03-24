package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pyrolass/hotel-reservation-go/entities"
)

func HandleGetUsers(c *fiber.Ctx) error {

	user := entities.User{
		FirstName: "John",
		LastName:  "Doe",
	}

	return c.JSON(
		user,
	)
}

func HandleGetUser(c *fiber.Ctx) error {

	userId := c.Params("id")

	return c.JSON(
		map[string]any{
			"message": "user" + userId,
		},
	)
}
