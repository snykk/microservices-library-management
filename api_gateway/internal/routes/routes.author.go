package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

type authorRoutes struct {
	handler handlers.AuthorHandler
	router  fiber.Router
}

func NewAuthorRoute(router fiber.Router, client clients.AuthorClient) *authorRoutes {
	handler := handlers.NewAuthorHandler(client)

	return &authorRoutes{
		router:  router,
		handler: handler,
	}
}

func (r *authorRoutes) Routes() {
	route := r.router.Group("/authors")
	route.Post("", r.handler.CreateAuthorHandler)
	route.Get("", r.handler.GetAllAuthorsHandler)
	route.Get("/:id", r.handler.GetAuthorByIdHandler)
	route.Put("/:id", r.handler.UpdateAuthorByIdHandler)
	route.Delete("/:id", r.handler.DeleteAuthorByIdHandler)

}
