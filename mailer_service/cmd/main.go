package main

import (
	"context"
	"io/ioutil"
	"log"
	"mailer_service/configs"
	"mailer_service/internal/constants"
	"mailer_service/internal/consumer"
	"mailer_service/internal/mailer"
	loggerPackage "mailer_service/pkg/logger"
	"mailer_service/pkg/rabbitmq"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	// Load app config
	if err := configs.InitializeAppConfig(); err != nil {
		log.Fatal("Failed to load app config", err)
	}
	log.Println("App configuration loaded")
}

func main() {
	// Email service env
	emailSenderBytes, err := ioutil.ReadFile(configs.AppConfig.EmailSenderContainerFile)
	if err != nil {
		log.Fatalf("Error reading email sender secret: %v", err)
	}

	emailPasswordBytes, err := ioutil.ReadFile(configs.AppConfig.EmailPasswordContainerFile)
	if err != nil {
		log.Fatalf("Error reading email password secret: %v", err)
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(configs.AppConfig.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Initialize rabbitMQPublisher
	rabbitMQPublisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQPublisher: %v", err)
	}
	defer rabbitMQPublisher.Close()

	// Declare exchanges
	err = rabbitMQPublisher.DeclareExchange(constants.LogExchange, constants.ExchangeTypeDirect)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Mailer Service
	mailerService, err := mailer.NewMailerService(string(emailSenderBytes), string(emailPasswordBytes))
	if err != nil {
		log.Fatalf("Failed to create mailer service %v", err)
	}

	// Logger
	var logger *loggerPackage.Logger
	if configs.AppConfig.LoggerWorkerType == constants.LoggerWorkerTypeSingle {
		logger = loggerPackage.NewLoggerSingleWorker(rabbitMQPublisher, configs.AppConfig.LoggerWorkerBufferSize)
	} else if configs.AppConfig.LoggerWorkerType == constants.LoggerWorkerTypeMultiple {
		logger = loggerPackage.NewLoggerMultipleWorker(rabbitMQPublisher, configs.AppConfig.LoggerWorkerNum, configs.AppConfig.LoggerWorkerBufferSize)
	}
	defer logger.Close()

	// Set up context with cancel for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Handle OS signals for shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start consumer in a goroutine
	go func() {
		if err := consumer.StartConsuming(ctx, ch, mailerService, logger); err != nil {
			log.Fatalf("Failed to start consuming: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-signalChan
	log.Println("Shutdown signal received, cleaning up resources...")

	// Cancel context to signal goroutines to stop
	cancel()

	// Allow time for graceful shutdown
	<-time.After(5 * time.Second)
	log.Println("Service shutdown completed gracefully.")
}
