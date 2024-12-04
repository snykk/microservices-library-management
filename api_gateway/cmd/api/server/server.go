package server

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/routes"
	"api_gateway/pkg/logger"
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
	"github.com/sirupsen/logrus"
)

type App struct {
	HttpServer *fiber.App
}

func NewApp() (*App, error) {
	// setup fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// middlewares
	app.Use(cors.New())
	app.Use(loggerFiber.New())

	// client gRPC
	authClient, err := clients.NewAuthClient()
	if err != nil {
		return nil, err
	}

	// routes
	router := app.Group("/api")
	routes.NewAuthRoute(router, authClient).Routes()

	return &App{
		HttpServer: app,
	}, nil
}

func (a *App) Run() error {
	// Graceful shutdown
	go func() {
		logger.InfoF("success to listen and serve on :%d", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer}, 80)
		if err := a.HttpServer.Listen(fmt.Sprintf(":%d", 80)); err != nil && err != fiber.ErrServiceUnavailable {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// make blocking channel and wait for signal
	<-quit
	log.Println("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.HttpServer.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("error when shutting down server: %v", err)
	}

	log.Println("Server exited properly.")
	return nil
}
