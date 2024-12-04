package main

import (
	"api_gateway/cmd/api/server"
	"api_gateway/internal/constants"
	"api_gateway/pkg/logger"
	"log"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// func init() {
// 	// if err := config.InitializeAppConfig(); err != nil {
// 	// 	logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
// 	// }
// 	// logger.Info("configuration loaded", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
// }

func init() {
	// Load Asia/Jakarta time zone globally
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	// Set the local timezone globally
	time.Local = loc
}

func main() {
	numCPU := runtime.NumCPU()
	logger.InfoF("The project is running on %d CPU(s)", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig}, numCPU)

	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(numCPU / 2)
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
	}
	if err := app.Run(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
	}
}
