package main

import (
	"book_service/configs"
	"book_service/internal/clients"
	"book_service/internal/constants"
	"book_service/internal/grpc_server"
	"book_service/internal/repository"
	"book_service/internal/service"
	loggerPackage "book_service/pkg/logger"
	"book_service/pkg/rabbitmq"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	protoBook "book_service/proto/book_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	// Database connection
	db, err := sqlx.Open("postgres", configs.AppConfig.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(configs.AppConfig.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Initialize RabbitMQ Publisher
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

	// Logger
	var logger *loggerPackage.Logger
	if configs.AppConfig.LoggerWorkerType == constants.LoggerWorkerTypeSingle {
		logger = loggerPackage.NewLoggerSingleWorker(rabbitMQPublisher, configs.AppConfig.LoggerWorkerBufferSize)
	} else if configs.AppConfig.LoggerWorkerType == constants.LoggerWorkerTypeMultiple {
		logger = loggerPackage.NewLoggerMultipleWorker(rabbitMQPublisher, configs.AppConfig.LoggerWorkerNum, configs.AppConfig.LoggerWorkerBufferSize)
	}
	defer logger.Close()

	// Clients
	authorClient, err := clients.NewAuthorClient()
	if err != nil {
		log.Fatalf("Failed to establish author client connection %v", err)
	}
	categoryClient, err := clients.NewCategoryClient()
	if err != nil {
		log.Fatalf("Failed to establish category client connection %v", err)
	}

	// Repository and Service Layer
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo, authorClient, categoryClient)

	// gRPC Server
	address := fmt.Sprintf(":%s", configs.AppConfig.GrpcPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	grpcServer := grpc.NewServer()
	bookServer := grpc_server.NewBookGRPCServer(bookService, logger)
	protoBook.RegisterBookServiceServer(grpcServer, bookServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("gRPC server is running on %s", address)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sigReceived := <-signalChannel
	log.Printf("Received signal: %v, initiating graceful shutdown...", sigReceived)

	// Log for starting cleanup
	log.Println("Starting cleanup tasks...")

	// Gracefully stop the gRPC server with a timeout
	gracefulShutdownTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	// Gracefully stop the gRPC server
	// grpcServer.Stop()
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")

	// Perform additional cleanup tasks ???
	// ...

	// Wait until all cleanup tasks are done
	<-ctx.Done() // Directly receive from the channel
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Timeout reached during graceful shutdown")
	}
	log.Println("Graceful shutdown completed successfully")
}
