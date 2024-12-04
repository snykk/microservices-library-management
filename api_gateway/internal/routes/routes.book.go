package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

type bookRoutes struct {
	handler handlers.BookHandler
	router  fiber.Router
}

func NewBookRoute(router fiber.Router, client clients.BookClient) *bookRoutes {
	handler := handlers.NewBookHandler(client)

	return &bookRoutes{
		router:  router,
		handler: handler,
	}
}

func (r *bookRoutes) Routes() {
	route := r.router.Group("/books")
	route.Post("", r.handler.CreateBookHandler)
	route.Get("", r.handler.GetAllBooksHandler)
	route.Get("/:id", r.handler.GetBookByIdHandler)
	route.Put("/:id", r.handler.UpdateBookByIdHandler)
	route.Delete("/:id", r.handler.DeleteBookByIdHandler)

}
