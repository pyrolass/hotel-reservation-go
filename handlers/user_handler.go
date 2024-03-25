package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pyrolass/hotel-reservation-go/db"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {

	userId := c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), userId)

	if err != nil {
		return err
	}

	return c.JSON(
		user,
	)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {

	users, err := h.userStore.GetAllUsers(c.Context())

	if err != nil {
		return err
	}

	return c.JSON(
		users,
	)
}
