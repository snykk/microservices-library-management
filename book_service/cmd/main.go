package main

import (
	"book_service/internal/clients"
	"book_service/internal/grpc_server"
	"book_service/internal/repository"
	"book_service/internal/service"
	"log"
	"net"
	"os"

	protoBook "book_service/proto/book_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	dsn := os.Getenv("DSN")

	// Initialize postgres connection
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

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
	bookServer := grpc_server.NewBookGRPCServer(bookService)
	protoBook.RegisterBookServiceServer(grpcServer, bookServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
