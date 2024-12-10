package middlewares

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	authClient clients.AuthClient
}

func NewAuthMiddleware(authClient clients.AuthClient) AuthMiddleware {
	return AuthMiddleware{authClient: authClient}
}

// Middleware for authentication (Validate Token)
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve requestID from context
		requestID, ok := c.Locals("requestID").(string)
		if !ok || requestID == "" {
			requestID = "unknown"
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			err := errors.New("authorization header is required")
			logger.Log.Error("Authentication failed",
				zap.String("request_id", requestID),
				zap.String("method", c.Method()),
				zap.String("url", c.OriginalURL()),
				zap.Error(err),
				zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", err))
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			err := errors.New("invalid Authorization header format")
			logger.Log.Error("Authentication failed",
				zap.String("request_id", requestID),
				zap.String("method", c.Method()),
				zap.String("url", c.OriginalURL()),
				zap.Error(err),
				zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", err))
		}
		token := tokenParts[1]

		dto := datatransfers.ValidateTokenRequest{
			Token: token,
		}
		res, err := m.authClient.ValidateToken(context.Background(), dto) // Call auth_service to validate token
		if err != nil || !res.Valid {
			logger.Log.Error("Authentication failed",
				zap.String("request_id", requestID),
				zap.String("method", c.Method()),
				zap.String("url", c.OriginalURL()),
				zap.Error(err),
				zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", err))
		}

		// Pass role and userID to context for further usage
		c.Locals("role", res.Role)
		c.Locals("userID", res.UserID)
		c.Locals("email", res.Email)

		// Log successful authentication
		logger.Log.Info("Authentication successful",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
			zap.String("role", res.Role),
			zap.String("user_id", res.UserID),
			zap.String("email", res.Email),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
		)

		return c.Next()
	}
}

// Middleware to authorize certain roles
func (m *AuthMiddleware) HasAuthority(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get request ID from context to track specific request
		requestID := c.Locals("requestID").(string)
		role := c.Locals("role").(string)

		// Log role check attempt
		logger.Log.Info("Checking user role for access",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
			zap.String("role", role),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
		)

		// Check if the user's role is authorized
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				// Log authorized access
				logger.Log.Info("Access granted based on role",
					zap.String("request_id", requestID),
					zap.String("method", c.Method()),
					zap.String("url", c.OriginalURL()),
					zap.String("role", role),
					zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
				)
				return c.Next() // Access granted
			}
		}

		// Log failed authorization attempt
		err := errors.New("access denied")
		logger.Log.Warn("Authorization failed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
			zap.String("role", role),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
		)

		return c.Status(fiber.StatusForbidden).JSON(datatransfers.ResponseError("Middleware authorization failed", err))
	}
}
