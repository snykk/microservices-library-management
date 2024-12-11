package logger

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/models"
	"api_gateway/pkg/rabbitmq"
	"log"
	"sync"
	"time"
)

type Logger struct {
	logChannel chan models.LogMessage
	wg         sync.WaitGroup
}

// NewLogger initializes the Logger with a RabbitMQ publisher and starts worker goroutines for processing logs.
func NewLogger(publisher *rabbitmq.Publisher, numWorkers int) *Logger {
	logChannel := make(chan models.LogMessage, 100)
	logger := &Logger{
		logChannel: logChannel,
	}

	// Menjalankan beberapa worker untuk memproses log secara paralel
	for i := 0; i < numWorkers; i++ {
		logger.wg.Add(1) // Menambah counter WaitGroup
		go func(workerID int) {
			defer logger.wg.Done() // Decrement counter ketika goroutine selesai
			for logMsg := range logChannel {
				log.Printf("Worker-%d: Processing log: %v\n", workerID, logMsg)
				if err := publisher.Publish(constants.LogExchange, constants.LogQueue, logMsg); err != nil {
					log.Printf("[%s] Worker-%d: Failed to publish log: %v\n", logMsg.Caller, workerID, err)
				} else {
					log.Printf("[%s] Worker-%d: Successfully published log\n", logMsg.Caller, workerID)
				}
			}
			log.Printf("Worker-%d: Exiting...\n", workerID)
		}(i)
	}

	return logger
}

// // NewLogger initializes the Logger with a RabbitMQ publisher and starts the worker goroutine.
// func NewLogger(publisher *rabbitmq.Publisher) *Logger {
// 	logChannel := make(chan models.LogMessage, 100)
// 	logger := &Logger{
// 		logChannel: logChannel,
// 	}

// 	// Worker goroutine for processing logs
// 	go func() {
// 		for logMsg := range logChannel {
// 			log.Println("ini channel")
// 			log.Println("ini message:", logMsg)
// 			if err := publisher.Publish(constants.LogExchange, constants.LogQueue, logMsg); err != nil {
// 				log.Printf("[%s] {%s} Failed to publish log to RabbitMQ: %v\n", logMsg.Caller, logMsg.Caller, err)
// 			}
// 			log.Printf("[%s] {%s} Success to publish log to RabbitMQ\n", logMsg.Caller, logMsg.Caller)
// 		}
// 	}()

// 	return logger
// }

// LogMessage sends a log message to the channel for processing.
func (l *Logger) LogMessage(caller, requestID, level, message string, extra map[string]interface{}, err error) {
	// Ensure that the logger is not nil
	if l == nil {
		log.Printf("[%s] Logger is nil, unable to log message", caller)
		return
	}

	logMsg := models.LogMessage{
		Timestamp:      time.Now(),
		Service:        constants.LogServiceApiGateway,
		Level:          level,
		XCorrelationID: requestID,
		Caller:         caller,
		Message:        message,
		Extra:          extra,
	}

	if err != nil {
		logMsg.Error = err.Error()
	}

	// Mengirim log ke channel
	select {
	case l.logChannel <- logMsg:
		log.Printf("[%s] Message sent to channel", caller)
	default:
		log.Printf("[%s] Failed to send message to channel, channel may be full or closed", caller)
	}
}

// Close ensures all workers are done and the logger is cleanly shut down.
func (l *Logger) Close() {
	if l == nil || l.logChannel == nil {
		return
	}

	// Tutup channel agar semua goroutine selesai membaca
	close(l.logChannel)

	// Tunggu semua goroutine selesai
	l.wg.Wait()
	log.Println("All log workers have exited.")
}
