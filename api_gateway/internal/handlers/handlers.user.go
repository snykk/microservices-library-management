package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	client clients.UserClient
	logger *logger.Logger
}

func NewUserHandler(client clients.UserClient, logger *logger.Logger) UserHandler {
	return UserHandler{
		client: client,
		logger: logger,
	}
}

func (b *UserHandler) GetAllUsersHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	// Call client to get all users
	resp, err := b.client.ListUsers(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID))
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get list users", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get list users", err))
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched all users successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("User data fetched successfully", resp))
}

func (b *UserHandler) GetMe(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	// Retrieve userID from locals (user's session or authentication context)
	userID := c.Locals("userID").(string)

	// Call client to get user by ID
	resp, err := b.client.GetUserById(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), userID)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get user data", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get user data", err))
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched user data successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("User data fetched successfully", resp))
}
