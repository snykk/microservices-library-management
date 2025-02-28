package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middlewares"
	"api_gateway/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type userRoute struct {
	router         fiber.Router
	authMiddleware middlewares.AuthMiddleware
	handler        handlers.UserHandler
}

func NewUserRoute(router fiber.Router, authMiddleware middlewares.AuthMiddleware, client clients.UserClient, logger *logger.Logger) *userRoute {
	handler := handlers.NewUserHandler(client, logger)

	return &userRoute{
		router:         router,
		authMiddleware: authMiddleware,
		handler:        handler,
	}
}

func (r *userRoute) Routes() {
	route := r.router.Group("/users")

	// Public routes (authentication required)
	route.Use(r.authMiddleware.Authenticate())
	route.Get("/me", r.handler.GetMe)

	// Admin routes (authentication and authorization required)
	adminOnly := r.authMiddleware.HasAuthority([]string{"admin"})
	route.Get("", adminOnly, r.handler.GetAllUsersHandler)
}
