package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"

	"go.uber.org/zap"

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
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Call client to get all users
	resp, err := b.client.ListUsers(c.Context())
	if err != nil {
		logger.Log.Error("Failed to get list users",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get list users", err))
	}

	// Log the success response
	logger.Log.Info("Fetched all users successfully",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("User data fetched successfully", resp))
}

func (b *UserHandler) GetMe(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Retrieve userID from locals (user's session or authentication context)
	userID := c.Locals("userID").(string)

	// Call client to get user by ID
	resp, err := b.client.GetUserById(c.Context(), userID)
	if err != nil {
		logger.Log.Error("Failed to get user data",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
			zap.String("user_id", userID),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get user data", err))
	}

	// Log the success response
	logger.Log.Info("Fetched user data successfully",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
		zap.String("user_id", userID),
	)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("User data fetched successfully", resp))
}
