package main

import (
	"book_service/internal/clients"
	"book_service/internal/constants"
	"book_service/internal/grpc_server"
	"book_service/internal/repository"
	"book_service/internal/service"
	"book_service/pkg/logger"
	"book_service/pkg/rabbitmq"
	"log"
	"net"
	"os"

	protoBook "book_service/proto/book_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	dsn := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

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

	// logger
	logger := logger.NewLoggerSingleWorker(rabbitMQPublisher, 100)
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
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	bookServer := grpc_server.NewBookGRPCServer(bookService, logger)
	protoBook.RegisterBookServiceServer(grpcServer, bookServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
