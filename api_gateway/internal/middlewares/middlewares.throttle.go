package middlewares

import (
	"api_gateway/configs"
	"api_gateway/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// ThrottleMiddleware provides rate limiting based on the client's IP address
func ThrottleMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        configs.AppConfig.MaxRequestPerMinute, // Maximum 100 requests
		Expiration: 1 * time.Minute,                       // Time window of 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use the client's IP address as the unique key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Respond with a 429 Too Many Requests error when the limit is exceeded
			return c.Status(fiber.StatusTooManyRequests).JSON(utils.ResponseError("You have exceeded the request limit. Please try again later.", "Too many requests"))
		},
	})
}
