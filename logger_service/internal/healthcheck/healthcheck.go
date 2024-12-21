package healthcheck

import (
	"context"
	"fmt"
	"log"
	"logger_service/configs"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

// PerformHealthCheck checks the status of MongoDB and RabbitMQ connections.
func PerformHealthCheck(rabbitConn *amqp.Connection, mongoClient *mongo.Client) bool {
	// Check MongoDB connection
	err := mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Println("MongoDB health check failed:", err)
		return false
	}

	// Check RabbitMQ connection
	_, err = rabbitConn.Channel()
	if err != nil {
		log.Println("RabbitMQ health check failed:", err)
		return false
	}

	return true
}

// StartHealthCheckServer starts the HTTP server for health check with graceful shutdown handling.
func StartHealthCheckServer(shutdownContext context.Context, rabbitConn *amqp.Connection, mongoClient *mongo.Client) {
	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Handle the health check route
	mux.HandleFunc("/healthy", func(w http.ResponseWriter, r *http.Request) {
		healthy := PerformHealthCheck(rabbitConn, mongoClient)
		if !healthy {
			http.Error(w, "Service is unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Service is healthy!!!")
	})

	// Initialize HTTP server with the custom ServeMux
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.AppConfig.HEALTH_API_PORT),
		Handler: mux,
	}

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Health check server running on :%s...\n", configs.AppConfig.HEALTH_API_PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Health check server failed: %v", err)
		}
	}()

	// Wait for the shutdown signal from the main application
	<-shutdownContext.Done()

	// Graceful shutdown: wait for ongoing requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown health check server: %v", err)
	} else {
		log.Println("Health check server gracefully stopped.")
	}
}
