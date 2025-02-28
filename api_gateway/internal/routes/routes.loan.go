package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middlewares"
	"api_gateway/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type loanRoutes struct {
	router         fiber.Router
	authMiddleware middlewares.AuthMiddleware
	handler        handlers.LoanHandler
}

func NewLoanRoute(router fiber.Router, authMiddleware middlewares.AuthMiddleware, client clients.LoanClient, logger *logger.Logger) *loanRoutes {
	handler := handlers.NewLoanHandler(client, logger)

	return &loanRoutes{
		router:         router,
		authMiddleware: authMiddleware,
		handler:        handler,
	}
}

func (r *loanRoutes) Routes() {
	route := r.router.Group("/loans")

	// Public routes (authentication required)
	route.Use(r.authMiddleware.Authenticate())
	route.Post("", r.handler.CreateLoanHandler)
	route.Post("/:id/return", r.handler.ReturnLoanHandler)
	route.Get("", r.handler.ListUserLoansHandler)

	// Admin routes (authentication and authorization required)
	adminOnly := r.authMiddleware.HasAuthority([]string{"admin"})
	route.Patch("/:id/status", adminOnly, r.handler.UpdateLoanStatusHandler)
	route.Get("/all", adminOnly, r.handler.ListLoansHandler)

	// avoid wildcard effect on `/all` endpoint
	route.Get("/:id", r.handler.GetLoanHandler)
}
