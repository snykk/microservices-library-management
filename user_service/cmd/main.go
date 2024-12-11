package main

import (
	"log"
	"net"
	"os"
	"user_service/internal/constants"
	"user_service/internal/grpc_server"
	"user_service/internal/repository"
	"user_service/internal/service"
	"user_service/pkg/logger"
	"user_service/pkg/rabbitmq"

	protoUser "user_service/proto/user_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	DSN := os.Getenv("DSN")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

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
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// gRPC Server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	userServer := grpc_server.NewUserGRPCServer(userService, logger)
	protoUser.RegisterUserServiceServer(grpcServer, userServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
