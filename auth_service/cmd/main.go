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
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	protoAuth "auth_service/proto/auth_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	amqp "github.com/rabbitmq/amqp091-go"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
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

	// Initialize RabbitMQ Publisher
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

	// Logger
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
	healthCheckServer := grpc_server.NewHealthGRPCServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthCheckServer)

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
