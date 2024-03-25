package handlers

import (
	"context"

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

	ctx := context.Background()

	user, err := h.userStore.GetUserById(ctx, userId)

	if err != nil {
		return err
	}

	return c.JSON(
		user,
	)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {

	ctx := context.Background()

	users, err := h.userStore.GetAllUsers(ctx)

	if err != nil {
		return err
	}

	return c.JSON(
		users,
	)
}
