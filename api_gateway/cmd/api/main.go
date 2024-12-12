package main

import (
	"api_gateway/cmd/api/server"
	"api_gateway/configs"
	"api_gateway/pkg/logger"
	"log"
	"runtime"
	"time"
)

func init() {
	// Load app config
	if err := configs.InitializeAppConfig(); err != nil {
		log.Fatal("Failed to load app config", err)
	}
	log.Println("App configuration loaded")

	// Load Asia/Jakarta time zone globally
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	time.Local = loc

	// // Initialize logger with default configuration
	// err = logger.Initialize(logger.LoggerConfig{
	// 	// OutputPaths: []string{"stdout", "../../logs/api-gateway.log"}, // local
	// 	OutputPaths: []string{"stdout", "/app/logs/api-gateway.log"}, // container
	// 	MaxSize:     10,                                              // 10 MB
	// 	MaxBackups:  5,
	// 	MaxAge:      30, // 30 days
	// 	Compress:    true,
	// 	IsDev:       false,
	// 	ServiceName: "api-gateway",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to initialize logger: %v", err)
	// }
}

func main() {
	defer logger.Sync()

	// Log CPU usage
	numCPU := runtime.NumCPU()
	log.Println("Starting application...")
	log.Printf("Cpu count: %d\n", numCPU)

	// Adjust GOMAXPROCS if applicable
	if numCPU > 2 {
		newProcs := numCPU / 2
		runtime.GOMAXPROCS(newProcs)
		log.Printf("The project is running on %d CPU(s)", numCPU)
	}

	// Initialize and run the application
	app, err := server.NewApp()
	if err != nil {
		log.Fatalf("Failed to create application instance: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Application run failed: %v", err)
	}
}
