package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type authorRoutes struct {
	router         fiber.Router
	authMiddleware middlewares.AuthMiddleware
	handler        handlers.AuthorHandler
}

func NewAuthorRoute(router fiber.Router, authMiddleware middlewares.AuthMiddleware, client clients.AuthorClient) *authorRoutes {
	handler := handlers.NewAuthorHandler(client)

	return &authorRoutes{
		router:         router,
		authMiddleware: authMiddleware,
		handler:        handler,
	}
}

func (r *authorRoutes) Routes() {
	route := r.router.Group("/authors")

	// Public routes (authentication required)
	route.Use(r.authMiddleware.Authenticate())
	route.Get("", r.handler.GetAllAuthorsHandler)
	route.Get("/:id", r.handler.GetAuthorByIdHandler)

	// Admin routes (authentication and authorization required)
	adminOnly := r.authMiddleware.HasAuthority([]string{"admin"})
	route.Post("", adminOnly, r.handler.CreateAuthorHandler)
	route.Put("/:id", adminOnly, r.handler.UpdateAuthorByIdHandler)
	route.Delete("/:id", adminOnly, r.handler.DeleteAuthorByIdHandler)

}
