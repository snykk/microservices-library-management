package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"context"
	"strconv"

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

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	extra := map[string]interface{}{
		"method":    c.Method(),
		"url":       c.OriginalURL(),
		"page":      page,
		"page_size": pageSize,
	}

	// Call client to get all users
	users, totalItems, totalPages, err := b.client.ListUsers(
		context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID),
		page,
		pageSize,
	)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get list users", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get list users", err))
	}

	extra["users_count"] = len(users)
	extra["total_items"] = totalItems
	extra["total_pages"] = totalPages
	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched all users successfully", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("User data fetched successfully", map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"currentPage": page,
			"page_size":   pageSize,
			"totalItems":  totalItems,
			"totalPages":  totalPages,
		},
	}))
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
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get user data", err))
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched user data successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("User data fetched successfully", resp))
}
