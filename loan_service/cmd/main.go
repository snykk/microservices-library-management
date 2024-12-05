package main

import (
	"loan_service/internal/grpc_server"
	"loan_service/internal/repository"
	"loan_service/internal/service"
	"log"
	"net"
	"os"

	protoLoan "loan_service/proto/loan_service"

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
	loanRepo := repository.NewLoanRepository(db)
	loanService := service.NewLoanService(loanRepo)

	// gRPC Server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	loanServer := grpc_server.NewLoanGRPCServer(loanService)
	protoLoan.RegisterLoanServiceServer(grpcServer, loanServer)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
