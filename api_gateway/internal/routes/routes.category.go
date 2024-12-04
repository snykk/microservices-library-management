package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

type categoryRoutes struct {
	handler handlers.CategoryHandler
	router  fiber.Router
}

func NewCategoryRoute(router fiber.Router, client clients.CategoryClient) *categoryRoutes {
	handler := handlers.NewCategoryHandler(client)

	return &categoryRoutes{
		router:  router,
		handler: handler,
	}
}

func (r *categoryRoutes) Routes() {
	route := r.router.Group("/categories")
	route.Post("", r.handler.CreateCategoryHandler)
	route.Get("", r.handler.GetAllCategoriesHandler)
	route.Get("/:id", r.handler.GetCategoryByIdHandler)
	route.Put("/:id", r.handler.UpdateCategoryByIdHandler)
	route.Delete("/:id", r.handler.DeleteCategoryByIdHandler)

}
