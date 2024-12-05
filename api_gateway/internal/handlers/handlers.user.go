package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	client clients.UserClient
}

func NewUserHandler(client clients.UserClient) UserHandler {
	return UserHandler{
		client: client,
	}
}

func (b *UserHandler) GetAllUsersHandler(c *fiber.Ctx) error {
	resp, err := b.client.ListUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get list users", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("User data fetched successfully", resp))
}

func (b *UserHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	resp, err := b.client.GetUserById(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get list users", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("User data fetched successfully", resp))
}
