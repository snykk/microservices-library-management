package middlewares

import (
	"api_gateway/internal/constants"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates and attaches a unique request ID to each incoming request
func RequestIDMiddleware(logger *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a new request ID
		requestID := uuid.New().String()

		// Attach the request ID to the context
		c.Locals(constants.ContextRequestIDKey, requestID)

		logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Request received1", map[string]interface{}{
			"method": c.Method(),
			"url":    c.OriginalURL(),
		}, nil)

		// Call the next handler
		return c.Next()
	}
}
