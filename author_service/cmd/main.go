package main

import (
	"author_service/internal/constants"
	"author_service/internal/grpc_server"
	"author_service/internal/repository"
	"author_service/internal/service"
	"author_service/pkg/logger"
	"author_service/pkg/rabbitmq"
	"log"
	"net"
	"os"

	protoAuthor "author_service/proto/author_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	DSN := os.Getenv("DSN")

	db, err := sqlx.Open("postgres", DSN)
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

	// Repository and Service Layer
	authorRepo := repository.NewAuthorRepository(db)
	authorService := service.NewAuthorService(authorRepo)

	// gRPC Server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	authorServer := grpc_server.NewAuthorGRPCServer(authorService, logger)
	protoAuthor.RegisterAuthorServiceServer(grpcServer, authorServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
