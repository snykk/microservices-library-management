package consumer

import (
	"context"
	"encoding/json"
	"log"
	"logger_service/configs"
	"logger_service/internal/constants"
	"logger_service/internal/logger"
	"logger_service/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func StartConsuming(ch *amqp.Channel, mongoClient *mongo.Client) error {
	// Declare exchange (direct exchange)
	err := ch.ExchangeDeclare(
		constants.LogExchange, // Exchange name
		"direct",              // Exchange type
		true,                  // Durable
		false,                 // Auto-deleted
		false,                 // Internal
		false,                 // No-wait
		nil,                   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Declare queue
	queueName := constants.LogQueue
	_, err = ch.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		queueName,      // Queue name
		queueName,      // Routing key (same as the queue name for direct exchange)
		"log_exchange", // Exchange name
		false,          // No-wait
		nil,            // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	// Start consuming logs from the queue
	go consumeQueue(ch, queueName, mongoClient)

	// Wait for consumer process
	log.Println("Waiting for log messages...")
	select {}
}

func consumeQueue(ch *amqp.Channel, queueName string, mongoClient *mongo.Client) {
	msgs, err := ch.Consume(
		queueName, // Queue
		"",        // Consumer
		false,     // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to start consuming from queue %s: %v", queueName, err)
	}

	// MongoDB collection for logs
	logCollection := mongoClient.Database(configs.AppConfig.MongoDB).Collection(configs.AppConfig.MongoCollection)

	for d := range msgs {
		var logMsg model.LogMessage
		if err := json.Unmarshal(d.Body, &logMsg); err != nil {
			log.Printf("Failed to parse log message: %v", err)
			d.Nack(false, true)
			continue
		}

		// Save log to MongoDB
		err := saveLogToMongoDB(logCollection, logMsg)
		if err != nil {
			log.Printf("Failed to save log to MongoDB: %v", err)
			d.Nack(false, true)
			continue
		}

		// Write log based on its level (e.g., Info, Error)
		err = writeLogToFile(logger.Log, logMsg)
		if err != nil {
			log.Printf("Failed to log to file: %v", err)
			d.Nack(false, true) // Negative acknowledgment with requeue
			continue
		}

		// Acknowledge the message after successful processing
		log.Printf("Successfully processed log message: %v", logMsg)
		d.Ack(false)
	}
}

func saveLogToMongoDB(collection *mongo.Collection, logMsg model.LogMessage) error {
	_, err := collection.InsertOne(context.TODO(), logMsg)
	return err
}

func writeLogToFile(logger *zap.Logger, logMsg model.LogMessage) error {
	fields := []zap.Field{
		zap.Time(constants.LogFieldTimeStamp, logMsg.Timestamp),
		zap.String(constants.LogFieldService, logMsg.Service),
		zap.String(constants.LogFieldXCorrelationID, logMsg.XCorrelationID),
		zap.String(constants.LogFieldCaller, logMsg.Caller),
		zap.Any(constants.LogFieldExtra, logMsg.Extra),
	}

	if logMsg.Error != nil {
		fields = append(fields, zap.String("error", *logMsg.Error))
	}

	// Tulis log berdasarkan level
	switch logMsg.Level {
	case constants.LogLevelInfo:
		logger.Info(logMsg.Message, fields...)
	case constants.LogLevelDebug:
		logger.Debug(logMsg.Message, fields...)
	case constants.LogLevelWarn:
		logger.Warn(logMsg.Message, fields...)
	case constants.LogLevelError:
		logger.Error(logMsg.Message, fields...)
	case constants.LogLevelPanic:
		logger.Panic(logMsg.Message, fields...)
	case constants.LogLevelFatal:
		logger.Fatal(logMsg.Message, fields...)
	default:
		// Default ke Info jika level tidak dikenali
		logger.Info(logMsg.Message, fields...)
	}

	return nil
}
