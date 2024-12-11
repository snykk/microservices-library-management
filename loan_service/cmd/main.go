package main

import (
	"loan_service/internal/clients"
	"loan_service/internal/constants"
	"loan_service/internal/grpc_server"
	"loan_service/internal/repository"
	"loan_service/internal/service"
	"loan_service/pkg/logger"
	"loan_service/pkg/rabbitmq"
	"log"
	"net"
	"os"

	protoLoan "loan_service/proto/loan_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	dsn := os.Getenv("DSN")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Client
	bookClient, err := clients.NewBookClient()
	if err != nil {
		log.Fatalf("Failed to establish book client connection %v", err)
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
		log.Fatalf("Failed to initialize rabbitMQPublisher: %v", err)
	}
	defer rabbitMQPublisher.Close()

	// Declare exchanges
	err = rabbitMQPublisher.DeclareExchange(constants.EmailExchange, constants.ExchangeTypeDirect)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	err = rabbitMQPublisher.DeclareExchange(constants.LogExchange, constants.ExchangeTypeDirect)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// logger
	logger := logger.NewLoggerSingleWorker(rabbitMQPublisher, 100)
	defer logger.Close()

	// Repository and Service Layer
	loanRepo := repository.NewLoanRepository(db)
	loanService := service.NewLoanService(loanRepo, bookClient, rabbitMQPublisher)

	// gRPC Server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	loanServer := grpc_server.NewLoanGRPCServer(loanService, logger)
	protoLoan.RegisterLoanServiceServer(grpcServer, loanServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
