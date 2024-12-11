package logger

import (
	"author_service/internal/constants"
	"author_service/internal/models"
	"author_service/pkg/rabbitmq"
	"log"
	"sync"
	"time"
)

type Logger struct {
	logChannel chan models.LogMessage
	wg         sync.WaitGroup
}

// NewLoggerSingle initializes the Logger with a RabbitMQ publisher and starts a single worker goroutine for processing logs.
func NewLoggerSingleWorker(publisher *rabbitmq.Publisher, bufferSize int) *Logger {
	logChannel := make(chan models.LogMessage, bufferSize)
	logger := &Logger{
		logChannel: logChannel,
	}

	// Worker goroutine for processing logs (single goroutine)
	go func() {
		for logMsg := range logChannel {
			log.Printf("Processing log: %v\n", logMsg)
			if err := publisher.Publish(constants.LogExchange, constants.LogQueue, logMsg); err != nil {
				log.Printf("[%s] Failed to publish log: %v\n", logMsg.Caller, err)
			} else {
				log.Printf("[%s] Successfully published log\n", logMsg.Caller)
			}
		}
		log.Println("Worker: Exiting...")
	}()

	return logger
}

// NewLoggerMultiple initializes the Logger with a RabbitMQ publisher and starts multiple worker goroutines for processing logs.
func NewLoggerMultipleWorker(publisher *rabbitmq.Publisher, numWorkers, bufferSize int) *Logger {
	logChannel := make(chan models.LogMessage, bufferSize)
	logger := &Logger{
		logChannel: logChannel,
	}

	// Start multiple workers to process logs concurrently
	for i := 0; i < numWorkers; i++ {
		logger.wg.Add(1) // Increment WaitGroup counter
		go func(workerID int) {
			defer logger.wg.Done() // Decrement WaitGroup counter when the goroutine finishes

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

// LogMessage sends a log message to the channel for processing.
func (l *Logger) LogMessage(caller, requestID, level, message string, extra map[string]interface{}, err error) {
	// Ensure that the logger is not nil
	if l == nil {
		log.Printf("[%s] Logger is nil, unable to log message", caller)
		return
	}

	logMsg := models.LogMessage{
		Timestamp:      time.Now(),
		Service:        constants.LogServiceAuth,
		Level:          level,
		XCorrelationID: requestID,
		Caller:         caller,
		Message:        message,
		Extra:          extra,
	}

	if err != nil {
		logMsg.Error = err.Error()
	}

	// Send the log to the channel
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

	// Close the channel to signal workers to stop reading
	close(l.logChannel)

	// Wait for all goroutines to finish
	l.wg.Wait()
	log.Println("All log workers have exited.")
}
