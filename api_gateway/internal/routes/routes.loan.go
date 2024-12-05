package routes

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type loanRoutes struct {
	router         fiber.Router
	authMiddleware middlewares.AuthMiddleware
	handler        handlers.LoanHandler
}

func NewLoanRoute(router fiber.Router, authMiddleware middlewares.AuthMiddleware, client clients.LoanClient) *loanRoutes {
	handler := handlers.NewLoanHandler(client)

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
	route.Get("/:id", r.handler.GetLoanHandler)
	route.Get("", r.handler.ListLoansHandler)

	// Admin routes (authentication and authorization required)
	adminOnly := r.authMiddleware.HasAuthority([]string{"admin"})
	route.Patch("/:id/status", adminOnly, r.handler.UpdateLoanStatusHandler)
}
