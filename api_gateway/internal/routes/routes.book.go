package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middlewares"
	"api_gateway/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type bookRoutes struct {
	router         fiber.Router
	authMiddleware middlewares.AuthMiddleware
	handler        handlers.BookHandler
}

func NewBookRoute(router fiber.Router, authMiddleware middlewares.AuthMiddleware, client clients.BookClient, authorClient clients.AuthorClient, categoryClient clients.CategoryClient, logger *logger.Logger) *bookRoutes {
	handler := handlers.NewBookHandler(client, authorClient, categoryClient, logger)

	return &bookRoutes{
		router:         router,
		authMiddleware: authMiddleware,
		handler:        handler,
	}
}

func (r *bookRoutes) Routes() {
	route := r.router.Group("/books")

	// Public routes (authentication required)
	route.Use(r.authMiddleware.Authenticate())
	route.Get("", r.handler.GetAllBooksHandler)
	route.Get("/:id", r.handler.GetBookByIdHandler)
	route.Get("/author/:authorId", r.handler.GetBooksByAuthorIdHandler)
	route.Get("/category/:categoryId", r.handler.GetBooksByCategoryIdHandler)

	// Admin routes (authentication and authorization required)
	adminOnly := r.authMiddleware.HasAuthority([]string{"admin"})
	route.Post("", adminOnly, r.handler.CreateBookHandler)
	route.Put("/:id", adminOnly, r.handler.UpdateBookByIdHandler)
	route.Delete("/:id", adminOnly, r.handler.DeleteBookByIdHandler)
}
