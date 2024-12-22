package middlewares

import (
	"api_gateway/configs"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// ThrottleMiddleware provides rate limiting based on the client's IP address
func ThrottleMiddleware(logger *logger.Logger) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        configs.AppConfig.MaxRequestPerMinute, // Maximum 100 requests
		Expiration: 1 * time.Minute,                       // Time window of 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use the client's IP address as the unique key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Retrieve requestID from context
			requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
			if !ok || requestID == "" {
				requestID = "unknown"
			}

			logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelWarn, "You have exceeded the request limit. Please try again later.", map[string]interface{}{
				"method": c.Method(),
				"url":    c.OriginalURL(),
			}, errors.New("too many requests"))

			// Respond with a 429 Too Many Requests error when the limit is exceeded
			return c.Status(fiber.StatusTooManyRequests).JSON(datatransfers.ResponseError("You have exceeded the request limit. Please try again later.", "Too many requests"))
		},
	})
}
