package main

import (
	"auth_service/internal/constants"
	"auth_service/internal/grpc_server"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"auth_service/pkg/jwt"
	"auth_service/pkg/logger"
	"auth_service/pkg/mailer"
	"auth_service/pkg/rabbitmq"
	"auth_service/pkg/redis"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	protoAuth "auth_service/proto/auth_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	log.Println("Initializing application...")
}

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	grpcPort := os.Getenv("GRPC_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	dsn := os.Getenv("DSN")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

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

	// Email service env
	emailSenderBytes, err := ioutil.ReadFile(os.Getenv("EMAIL_SENDER_CONTAINER_FILE"))
	if err != nil {
		log.Fatalf("Error reading email sender secret: %v", err)
	}

	emailPasswordBytes, err := ioutil.ReadFile(os.Getenv("EMAIL_PASSWORD_CONTAINER_FILE"))
	if err != nil {
		log.Fatalf("Error reading email password secret: %v", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Redis cache
	redisCache := redis.NewRedisCache(fmt.Sprintf("%s:%s", redisHost, redisPort), 0, redisPassword, 10*time.Minute)

	// JWT Service
	jwtService := jwt.NewJWTService(jwtSecret, "auth_service", 15, 60)

	// Mailer Service
	mailerService := mailer.NewOTPMailer(string(emailSenderBytes), string(emailPasswordBytes))

	// Repository and Service Layer
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtService, mailerService, rabbitMQPublisher)

	// gRPC Server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	authServer := grpc_server.NewAuthServer(authService, redisCache, logger)
	protoAuth.RegisterAuthServiceServer(grpcServer, authServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
