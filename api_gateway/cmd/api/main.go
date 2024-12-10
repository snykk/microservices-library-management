package main

import (
	"api_gateway/cmd/api/server"
	"api_gateway/internal/constants"
	"api_gateway/pkg/logger"
	"log"
	"runtime"
	"time"

	"go.uber.org/zap"
)

func init() {
	// Load Asia/Jakarta time zone globally
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	time.Local = loc

	// Initialize logger with default configuration
	err = logger.Initialize(logger.LoggerConfig{
		// OutputPaths: []string{"stdout", "../../logs/api-gateway.log"}, // local
		OutputPaths: []string{"stdout", "/app/logs/api-gateway.log"}, // container
		MaxSize:     10,                                              // 10 MB
		MaxBackups:  5,
		MaxAge:      30, // 30 days
		Compress:    true,
		IsDev:       false,
		ServiceName: "api-gateway",
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
}

func main() {
	defer logger.Sync()

	// Log CPU usage
	numCPU := runtime.NumCPU()
	logger.Log.Info("Starting application",
		zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		zap.Int("CPU count", numCPU),
	)

	// Adjust GOMAXPROCS if applicable
	if numCPU > 2 {
		newProcs := numCPU / 2
		runtime.GOMAXPROCS(newProcs)
		logger.Log.Info("Adjusted GOMAXPROCS",
			zap.Int("new GOMAXPROCS", newProcs),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
	}

	// Initialize and run the application
	app, err := server.NewApp()
	if err != nil {
		logger.Log.Panic("Failed to create application instance",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
	}

	if err := app.Run(); err != nil {
		logger.Log.Fatal("Application run failed",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategorySetup),
		)
	}
}
