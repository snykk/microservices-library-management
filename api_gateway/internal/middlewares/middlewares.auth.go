package middlewares

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
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
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", errors.New("authorization header is required")))
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", errors.New("invalid Authorization header format")))
		}
		token := tokenParts[1]

		dto := datatransfers.ValidateTokenRequest{
			Token: token,
		}
		res, err := m.authClient.ValidateToken(context.Background(), dto) // Call auth_service to validate token
		if err != nil || !res.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Middleware authentication failed", err))
		}

		// Pass role and userID to context
		c.Locals("role", res.Role)
		c.Locals("userID", res.UserID)

		return c.Next()
	}
}

// Middleware to authorize certain roles
func (m *AuthMiddleware) HasAuthority(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next() // Access granted
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(datatransfers.ResponseError("Middleware authorization failed", errors.New("access denied")))
	}
}
