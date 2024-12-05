package main

import (
	"category_service/internal/grpc_server"
	"category_service/internal/repository"
	"category_service/internal/service"
	"log"
	"net"
	"os"

	protoCategory "category_service/proto/category_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
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

	// Repository and Service Layer
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)

	// gRPC Server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	categoryServer := grpc_server.NewCategoryGRPCServer(categoryService)
	protoCategory.RegisterCategoryServiceServer(grpcServer, categoryServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
