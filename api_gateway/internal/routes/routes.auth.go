package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"
	"api_gateway/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type authRoutes struct {
	handler handlers.AuthHandler
	router  fiber.Router
}

func NewAuthRoute(router fiber.Router, client clients.AuthClient, logger *logger.Logger) *authRoutes {
	handler := handlers.NewAuthHandler(client, logger)

	return &authRoutes{
		router:  router,
		handler: handler,
	}
}

func (r *authRoutes) Routes() {
	route := r.router.Group("/auth")
	route.Post("/register", r.handler.RegisterHandler)
	route.Post("/send-otp", r.handler.SendOtpHandler)
	route.Post("/verify-email", r.handler.VerifyEmailHandler)
	route.Post("/login", r.handler.LoginHandler)
	route.Post("/validate-token", r.handler.ValidateTokenHandler)

}
