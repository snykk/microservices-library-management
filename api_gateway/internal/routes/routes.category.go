package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type categoryRoutes struct {
	router         fiber.Router
	authMiddleware middlewares.AuthMiddleware
	handler        handlers.CategoryHandler
}

func NewCategoryRoute(router fiber.Router, authMiddleware middlewares.AuthMiddleware, client clients.CategoryClient) *categoryRoutes {
	handler := handlers.NewCategoryHandler(client)

	return &categoryRoutes{
		router:         router,
		authMiddleware: authMiddleware,
		handler:        handler,
	}
}

func (r *categoryRoutes) Routes() {
	route := r.router.Group("/categories")

	// Public routes (authentication required)
	route.Use(r.authMiddleware.Authenticate())
	route.Get("", r.handler.GetAllCategoriesHandler)
	route.Get("/:id", r.handler.GetCategoryByIdHandler)

	// Admin routes (authentication and authorization required)
	adminOnly := r.authMiddleware.HasAuthority([]string{"admin"})
	route.Post("", adminOnly, r.handler.CreateCategoryHandler)
	route.Put("/:id", adminOnly, r.handler.UpdateCategoryByIdHandler)
	route.Delete("/:id", adminOnly, r.handler.DeleteCategoryByIdHandler)

}
