package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pyrolass/hotel-reservation-go/db"
	"github.com/pyrolass/hotel-reservation-go/entities"
	"github.com/pyrolass/hotel-reservation-go/middleware"
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

func (h *UserHandler) HandleUserLogin(c *fiber.Ctx) error {
	var loginParams entities.LoginParams

	if err := c.BodyParser(&loginParams); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), loginParams.Email)

	if err != nil {
		return err
	}

	if !entities.CheckPasswordHash(loginParams.Password, user.EncryptedPassword) {
		return c.Status(401).JSON(
			map[string]string{
				"message": "Wrong Creds",
			},
		)
	}

	token, err := middleware.GenerateToken(*user)

	if err != nil {
		return err
	}

	return c.JSON(
		map[string]any{
			"message": "Success",
			"token":   token,
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

	token, err := middleware.GenerateToken(*insertedUser)

	if err != nil {
		return err
	}

	return c.JSON(
		map[string]any{
			"token": token,
			"user":  insertedUser,
		},
	)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var params entities.UpdateUserParams

	userId := c.Params("id")

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	err := h.userStore.UpdateUser(c.Context(), userId, params)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(
		map[string]string{
			"message": "User updated",
		},
	)
}
