package middlewares

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authClient clients.AuthClient
	logger     *logger.Logger
}

func NewAuthMiddleware(authClient clients.AuthClient, logger *logger.Logger) AuthMiddleware {
	return AuthMiddleware{
		authClient: authClient,
		logger:     logger,
	}
}

// Middleware for authentication (Validate Token)
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve requestID from context
		requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
		if !ok || requestID == "" {
			requestID = "unknown"
		}

		extra := map[string]interface{}{
			"method": c.Method(),
			"url":    c.OriginalURL(),
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			err := errors.New("authorization header is required")
			// logger.Log.Error("Authentication failed",
			// 	zap.String("request_id", requestID),
			// 	zap.String("method", c.Method()),
			// 	zap.String("url", c.OriginalURL()),
			// 	zap.Error(err),
			// 	zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
			// )
			// if err := m.rabbitMQPublisher.Publish(constants.LogExchange, constants.LogQueue, models.LogMessage{
			// 	Timestamp:      time.Now(),
			// 	Service:        constants.LogServiceApiGateway,
			// 	Level:          constants.LogLevelError,
			// 	XCorrelationID: requestID,
			// 	Caller:         utils.GetLocation(),
			// 	Message:        "Authentication failed",
			// 	Error:          err.Error(),
			// 	Extra: map[string]interface{}{
			// 		"method": c.Method(),
			// 		"url":    c.OriginalURL(),
			// 	},
			// }); err != nil {
			// 	log.Printf("[%s] Failed to publish log to RabbitMQ: %v\n", utils.GetLocation(), err)
			// }
			m.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Authentication failed", extra, err)

			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", err))
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			err := errors.New("invalid Authorization header format")
			// logger.Log.Error("Authentication failed",
			// 	zap.String("request_id", requestID),
			// 	zap.String("method", c.Method()),
			// 	zap.String("url", c.OriginalURL()),
			// 	zap.Error(err),
			// 	zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
			// )

			// if err := m.rabbitMQPublisher.Publish(constants.LogExchange, constants.LogQueue, models.LogMessage{
			// 	Timestamp:      time.Now(),
			// 	Service:        constants.LogServiceApiGateway,
			// 	Level:          constants.LogLevelError,
			// 	XCorrelationID: requestID,
			// 	Caller:         utils.GetLocation(),
			// 	Message:        "Authentication failed",
			// 	Error:          err.Error(),
			// 	Extra: map[string]interface{}{
			// 		"method": c.Method(),
			// 		"url":    c.OriginalURL(),
			// 	},
			// }); err != nil {
			// 	log.Printf("[%s] Failed to publish log to RabbitMQ: %v\n", utils.GetLocation(), err)
			// }
			m.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Authentication failed", extra, err)
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", err))
		}
		token := tokenParts[1]

		dto := datatransfers.ValidateTokenRequest{
			Token: token,
		}
		res, err := m.authClient.ValidateToken(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), dto) // Call auth_service to validate token
		if err != nil || !res.Valid {
			// logger.Log.Error("Authentication failed",
			// 	zap.String("request_id", requestID),
			// 	zap.String("method", c.Method()),
			// 	zap.String("url", c.OriginalURL()),
			// 	zap.Error(err),
			// 	zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
			// )
			// if err := m.rabbitMQPublisher.Publish(constants.LogExchange, constants.LogQueue, models.LogMessage{
			// 	Timestamp:      time.Now(),
			// 	Service:        constants.LogServiceApiGateway,
			// 	Level:          constants.LogLevelError,
			// 	XCorrelationID: requestID,
			// 	Caller:         utils.GetLocation(),
			// 	Message:        "Authentication failed",
			// 	Error:          err.Error(),
			// 	Extra: map[string]interface{}{
			// 		"method": c.Method(),
			// 		"url":    c.OriginalURL(),
			// 	},
			// }); err != nil {
			// 	log.Printf("[%s] Failed to publish log to RabbitMQ: %v\n", utils.GetLocation(), err)
			// }
			m.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Authentication failed", extra, err)
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", err))
		}

		// Pass role and userID to context for further usage
		c.Locals("role", res.Role)
		c.Locals("userID", res.UserID)
		c.Locals("email", res.Email)

		// Log successful authentication
		// logger.Log.Info("Authentication successful",
		// 	zap.String("request_id", requestID),
		// 	zap.String("method", c.Method()),
		// 	zap.String("url", c.OriginalURL()),
		// 	zap.String("role", res.Role),
		// 	zap.String("user_id", res.UserID),
		// 	zap.String("email", res.Email),
		// 	zap.String(constants.LoggerCategory, constants.LoggerCategoryMiddleware),
		// )

		// if err := m.rabbitMQPublisher.Publish(constants.LogExchange, constants.LogQueue, models.LogMessage{
		// 	Timestamp:      time.Now(),
		// 	Service:        constants.LogServiceApiGateway,
		// 	Level:          constants.LogLevelInfo,
		// 	XCorrelationID: requestID,
		// 	Caller:         utils.GetLocation(),
		// 	Message:        "Authentication successful",
		// 	Extra: map[string]interface{}{
		// 		"method":  c.Method(),
		// 		"url":     c.OriginalURL(),
		// 		"role":    res.Role,
		// 		"user_id": res.UserID,
		// 		"email":   res.Email,
		// 	},
		// }); err != nil {
		// 	log.Printf("[%s] Failed to publish log to RabbitMQ: %v\n", utils.GetLocation(), err)
		// }

		extra["role"] = res.Role
		extra["user_id"] = res.UserID
		extra["email"] = res.Email

		m.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Authentication successful", extra, nil)

		return c.Next()
	}
}

// Middleware to authorize certain roles
func (m *AuthMiddleware) HasAuthority(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get request ID from context to track specific request
		requestID := c.Locals(constants.ContextRequestIDKey).(string)
		role := c.Locals("role").(string)

		extra := map[string]interface{}{
			"method": c.Method(),
			"url":    c.OriginalURL(),
			"role":   role,
		}

		// Log role check attempt
		m.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Checking user role for access", extra, nil)

		// Check if the user's role is authorized
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				// Log authorized access
				m.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Access granted based on role", extra, nil)

				return c.Next() // Access granted
			}
		}

		// Log failed authorization attempt
		err := errors.New("access denied")
		m.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelWarn, "Authorization failed", extra, err)

		return c.Status(fiber.StatusForbidden).JSON(datatransfers.ResponseError("Middleware authorization failed", err))
	}
}
