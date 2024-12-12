package main

import (
	"auth_service/configs"
	"auth_service/internal/constants"
	"auth_service/internal/grpc_server"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"auth_service/pkg/jwt"
	loggerPackage "auth_service/pkg/logger"
	"auth_service/pkg/mailer"
	"auth_service/pkg/rabbitmq"
	"auth_service/pkg/redis"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	protoAuth "auth_service/proto/auth_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

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
	// Connect to RabbitMQ
	conn, err := amqp.Dial(configs.AppConfig.RabbitMQURL)
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
	var logger *loggerPackage.Logger
	if configs.AppConfig.LoggerWorkerType == constants.LoggerWorkerTypeSingle {
		logger = loggerPackage.NewLoggerSingleWorker(rabbitMQPublisher, configs.AppConfig.LoggerWorkerBufferSize)
	} else if configs.AppConfig.LoggerWorkerType == constants.LoggerWorkerTypeMultiple {
		logger = loggerPackage.NewLoggerMultipleWorker(rabbitMQPublisher, configs.AppConfig.LoggerWorkerNum, configs.AppConfig.LoggerWorkerBufferSize)
	}
	defer logger.Close()

	// Email service env
	emailSenderBytes, err := ioutil.ReadFile(configs.AppConfig.EmailSenderContainerFile)
	if err != nil {
		log.Fatalf("Error reading email sender secret: %v", err)
	}

	emailPasswordBytes, err := ioutil.ReadFile(configs.AppConfig.EmailPasswordContainerFile)
	if err != nil {
		log.Fatalf("Error reading email password secret: %v", err)
	}

	db, err := sql.Open("pgx", configs.AppConfig.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Redis cache
	redisCache := redis.NewRedisCache(fmt.Sprintf("%s:%s", configs.AppConfig.RedisHost, configs.AppConfig.RedisPort), configs.AppConfig.RedisDB, configs.AppConfig.RedisPassword, configs.AppConfig.RedisDefaultExp)

	// JWT Service
	jwtService := jwt.NewJWTService(configs.AppConfig.JwtSecret, configs.AppConfig.JwtIssuer, configs.AppConfig.JwtExpAccessToken, configs.AppConfig.JwtExpRefreshToken)

	// Mailer Service
	mailerService := mailer.NewOTPMailer(string(emailSenderBytes), string(emailPasswordBytes))

	// Repository and Service Layer
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtService, mailerService, rabbitMQPublisher)

	// gRPC Server
	address := fmt.Sprintf(":%s", configs.AppConfig.GrpcPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	grpcServer := grpc.NewServer()
	authServer := grpc_server.NewAuthServer(authService, redisCache, logger)
	protoAuth.RegisterAuthServiceServer(grpcServer, authServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
