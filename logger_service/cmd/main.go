package main

import (
	"context"
	"log"
	"logger_service/configs"
	"logger_service/internal/consumer"
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

	// Start consuming logs
	err = consumer.StartConsuming(ch, mongoClient)
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
}
