package middlewares

import (
	"api_gateway/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDMiddleware generates and attaches a unique request ID to each incoming request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a new request ID
		requestID := uuid.New().String()

		// Attach the request ID to the context
		c.Locals("requestID", requestID)

		logger.Log.Info("Request received",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)

		// Call the next handler
		return c.Next()
	}
}
