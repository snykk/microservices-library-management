package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"logger_service/configs"
	"logger_service/internal/consumer"
	"logger_service/internal/healthcheck"
	"logger_service/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// Load app config
	if err := configs.InitializeAppConfig(); err != nil {
		log.Fatal("Failed to load app config", err)
	}
	log.Println("App configuration loaded")
}

func main() {
	// Initialize logger
	err := logger.Initialize(logger.LoggerConfig{
		OutputPaths: []string{"stdout", configs.AppConfig.LogPath}, // container
		MaxSize:     10,                                            // 10 MB
		MaxBackups:  5,
		MaxAge:      30, // 30 days
		Compress:    true,
		IsDev:       false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(configs.AppConfig.MongoURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
		}
		log.Println("MongoDB connection closed.")
	}()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(configs.AppConfig.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Listen for interrupt signals to trigger graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	// Start health check server in a separate goroutine
	go healthcheck.StartHealthCheckServer(ctx, conn, mongoClient)

	// Start consuming logs in a separate goroutine
	go func() {
		if err := consumer.StartConsuming(ctx, ch, mongoClient); err != nil {
			log.Fatalf("Failed to start consuming: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigint
	log.Println("Received shutdown signal, initiating graceful shutdown...")

	// Cancel the context to stop consumers gracefully
	cancel()

	// Allow time for graceful shutdown
	<-time.After(5 * time.Second)
	log.Println("Graceful shutdown completed.")

}
