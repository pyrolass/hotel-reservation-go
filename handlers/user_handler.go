package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/entities"
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

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {

	return nil
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {

	userId := c.Params("id")

	err := h.userStore.DeleteUser(c.Context(), userId)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(
		map[string]string{
			"message": "User deleted",
		},
	)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {

	var params entities.CreateUserParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	user, err := entities.NewUserFromParams(params)

	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.CreateUser(c.Context(), user)

	if err != nil {
		return err
	}

	return c.JSON(
		insertedUser,
	)
}
