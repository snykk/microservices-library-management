package server

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/middlewares"
	"api_gateway/internal/routes"
	"api_gateway/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerFiber "github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

type App struct {
	HttpServer *fiber.App
}

func NewApp() (*App, error) {
	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Client gRPC
	authClient, err := clients.NewAuthClient()
	if err != nil {
		logger.Log.Error("Failed to create AuthClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		return nil, err
	}
	bookClient, err := clients.NewBookClient()
	if err != nil {
		logger.Log.Error("Failed to create BookClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		return nil, err
	}
	categoryClient, err := clients.NewCategoryClient()
	if err != nil {
		logger.Log.Error("Failed to create CategoryClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		return nil, err
	}
	authorClient, err := clients.NewAuthorClient()
	if err != nil {
		logger.Log.Error("Failed to create AuthorClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		return nil, err
	}
	userClient, err := clients.NewUserClient()
	if err != nil {
		logger.Log.Error("Failed to create UserClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		return nil, err
	}
	loanClient, err := clients.NewLoanClient()
	if err != nil {
		logger.Log.Error("Failed to create LoanClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		return nil, err
	}

	// Fiber middlewares
	app.Use(middlewares.RequestIDMiddleware()) // Add the RequestID middleware
	app.Use(cors.New())
	app.Use(loggerFiber.New())

	// Authentication middleware
	authMiddleware := middlewares.NewAuthMiddleware(authClient)

	// Routes
	router := app.Group("/api")
	routes.NewAuthRoute(router, authClient).Routes()
	routes.NewBookRoute(router, authMiddleware, bookClient, authorClient, categoryClient).Routes()
	routes.NewCategoryRoute(router, authMiddleware, categoryClient, bookClient).Routes()
	routes.NewAuthorRoute(router, authMiddleware, authorClient, bookClient).Routes()
	routes.NewUserRoute(router, authMiddleware, userClient).Routes()
	routes.NewLoanRoute(router, authMiddleware, loanClient).Routes()

	logger.Log.Info("Fiber app initialized successfully",
		zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
	)

	return &App{
		HttpServer: app,
	}, nil
}

func (a *App) Run() error {
	// Start server in a goroutine
	go func() {
		address := ":80"
		logger.Log.Info("Server is starting",
			zap.String("address", address),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		if err := a.HttpServer.Listen(address); err != nil && err != fiber.ErrServiceUnavailable {
			logger.Log.Fatal("Failed to listen and serve",
				zap.Error(err),
				zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
			)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-quit
	logger.Log.Warn("Shutdown signal received",
		zap.String("signal", sig.String()),
		zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
	)

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.HttpServer.ShutdownWithContext(ctx); err != nil {
		logger.Log.Error("Error during server shutdown",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
		return fmt.Errorf("error when shutting down server: %v", err)
	}

	logger.Log.Info("Server exited properly",
		zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
	)
	return nil
}
