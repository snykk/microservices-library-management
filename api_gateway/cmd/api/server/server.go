package server

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/middlewares"
	"api_gateway/internal/routes"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/rabbitmq"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerFiber "github.com/gofiber/fiber/v2/middleware/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	HttpServer *fiber.App
}

func NewApp() (*App, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return nil, err
	}
	// defer conn.Close()
	log.Println("Success connect to rabbitMQ")

	// Initialize publisher
	rabbitMQPublisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Println("Failed to initialize RabbitMQPublisher:", err)
		return nil, err
	}
	// defer rabbitMQPublisher.Close()

	// Declare exchanges
	err = rabbitMQPublisher.DeclareExchange("log_exchange", "direct")
	if err != nil {
		log.Println("Failed to declare exchange:", err)
		return nil, err
	}

	// logger
	logger := logger.NewLoggerSingleWorker(rabbitMQPublisher, 100)

	// Client gRPC
	authClient, err := clients.NewAuthClient(logger)
	if err != nil {
		log.Println("Failed to create AuthClient:", err)
		return nil, err
	}
	bookClient, err := clients.NewBookClient(logger)
	if err != nil {
		log.Println("Failed to create BookClient:", err)
		return nil, err
	}
	categoryClient, err := clients.NewCategoryClient(logger)
	if err != nil {
		log.Println("Failed to create CategoryClient:", err)
		return nil, err
	}
	authorClient, err := clients.NewAuthorClient(logger)
	if err != nil {
		log.Println("Failed to create AuthorClient:", err)
		return nil, err
	}
	userClient, err := clients.NewUserClient(logger)
	if err != nil {
		log.Println("Failed to create UserClient:", err)
		return nil, err
	}
	loanClient, err := clients.NewLoanClient(logger)
	if err != nil {
		log.Println("Failed to create LoanClient:", err)
		return nil, err
	}

	// Fiber middlewares
	app.Use(middlewares.RequestIDMiddleware(logger)) // Add the RequestID middleware
	app.Use(cors.New())
	app.Use(loggerFiber.New())

	// Authentication middleware
	authMiddleware := middlewares.NewAuthMiddleware(authClient, logger)

	// Routes
	router := app.Group("/api")
	routes.NewAuthRoute(router, authClient, logger).Routes()
	routes.NewBookRoute(router, authMiddleware, bookClient, authorClient, categoryClient, logger).Routes()
	routes.NewCategoryRoute(router, authMiddleware, categoryClient, bookClient, logger).Routes()
	routes.NewAuthorRoute(router, authMiddleware, authorClient, bookClient, logger).Routes()
	routes.NewUserRoute(router, authMiddleware, userClient, logger).Routes()
	routes.NewLoanRoute(router, authMiddleware, loanClient, logger).Routes()

	log.Println("Fiber app initialized successfully")

	return &App{
		HttpServer: app,
	}, nil
}

func (a *App) Run() error {
	// Start server in a goroutine
	go func() {
		address := ":80"
		log.Println("Server is starting on", address)
		if err := a.HttpServer.Listen(address); err != nil && err != fiber.ErrServiceUnavailable {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-quit
	log.Println("Shutdown signal reveived:", sig.String())

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.HttpServer.ShutdownWithContext(ctx); err != nil {
		log.Println("Error during server shutdown:", err)
		return fmt.Errorf("error when shutting down server: %v", err)
	}

	log.Println("server exited properly")
	return nil
}
