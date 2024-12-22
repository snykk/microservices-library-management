package routes

import (
	"api_gateway/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

type commonRoutes struct {
	router  fiber.Router
	handler handlers.CommonHandler
}

func NewCommonRoute(router fiber.Router) *commonRoutes {
	handler := handlers.NewCommonHandler()
	return &commonRoutes{
		router:  router,
		handler: handler,
	}
}

func (r *commonRoutes) Routes() {
	r.router.Get("", r.handler.Root)
	r.router.Get("/healthy", r.handler.Healthy)
}
