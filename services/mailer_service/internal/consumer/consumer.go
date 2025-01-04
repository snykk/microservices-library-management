package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mailer_service/internal/constants"
	"mailer_service/internal/mailer"
	"mailer_service/pkg/logger"
	"mailer_service/pkg/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OTPMessage struct {
	RequestID string `json:"X-Correlation-ID"` // for logging purpose
	Email     string `json:"email"`
	OTP       string `json:"otp"`
}

type LoanNotificationMessage struct {
	RequestID string    `json:"X-Correlation-ID"` // for logging purpose
	Email     string    `json:"email"`
	Book      string    `json:"book"`
	Due       time.Time `json:"due"`
}

type ReturnNotificationMessage struct {
	RequestID string `json:"X-Correlation-ID"` // for logging purpose
	Email     string `json:"email"`
	Book      string `json:"book"`
}

func StartConsuming(ctx context.Context, ch *amqp.Channel, mailerService mailer.MailerService, logger *logger.Logger) error {
	// Declare exchange (e.g., direct exchange)
	err := ch.ExchangeDeclare(
		constants.EmailExchange, // Exchange name
		"direct",                // Exchange type
		true,                    // Durable
		false,                   // Auto-deleted
		false,                   // Internal
		false,                   // No-wait
		nil,                     // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Declare queues and bind them to the exchange
	queues := []string{constants.OTPQueue, constants.LoanNotificationQueue, constants.ReturnNotificationQueue}
	for _, queueName := range queues {
		// Declare queue
		_, err := ch.QueueDeclare(
			queueName,
			true,  // Durable
			false, // Delete when unused
			false, // Exclusive
			false, // No-wait
			nil,   // Arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare queue %s: %v", queueName, err)
		}

		// Bind queue to the exchange
		err = ch.QueueBind(
			queueName,               // Queue name
			queueName,               // Routing key (same as the queue name for direct exchange)
			constants.EmailExchange, // Exchange name
			false,                   // No-wait
			nil,                     // Arguments
		)
		if err != nil {
			log.Fatalf("Failed to bind queue %s to exchange: %v", queueName, err)
		}
	}

	// Consume all queues in goroutines
	for _, queueName := range queues {
		go consumeQueue(ctx, ch, queueName, mailerService, logger)
	}

	// Wait for context cancellation
	log.Println("Waiting for messages...")
	<-ctx.Done()
	log.Println("Context canceled, stopping consumer...")
	return nil
}

func consumeQueue(ctx context.Context, ch *amqp.Channel, queueName string, mailerService mailer.MailerService, logger *logger.Logger) {
	msgs, err := ch.Consume(
		queueName, // Queue
		"",        // Consumer
		false,     // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		log.Fatalf("Failed to start consuming from queue %s: %v", queueName, err)
	}

	for {
		select {
		case <-ctx.Done(): // Stop consuming when the context is canceled
			log.Println("Graceful shutdown: stopping message consumption")
			return
		case d := <-msgs:
			var requestID string = "unknown"
			var extra map[string]interface{}
			switch queueName {
			case constants.OTPQueue:
				var message OTPMessage
				if err := json.Unmarshal(d.Body, &message); err != nil {
					log.Printf("Failed to parse OTP message: %v", err)
					d.Nack(false, true)
					log.Println("Retry later...")
					logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to parse OTP message: %v", err), extra, err)
					continue
				}
				requestID = message.RequestID
				extra = map[string]interface{}{
					"email": message.Email,
					"otp":   message.OTP,
				}
				err = mailerService.SendOTP(message.Email, message.OTP)

			case constants.LoanNotificationQueue:
				var message LoanNotificationMessage
				if err := json.Unmarshal(d.Body, &message); err != nil {
					log.Printf("Failed to parse loan notification message: %v", err)
					d.Nack(false, true)
					log.Println("Retry later...")
					logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to parse loan notification message: %v", err), extra, err)
					continue
				}
				requestID = message.RequestID
				extra = map[string]interface{}{
					"email":      message.Email,
					"book_title": message.Book,
					"due_date":   message.Due,
				}
				err = mailerService.SendLoanNotification(message.Email, message.Book, message.Due)

			case constants.ReturnNotificationQueue:
				var message ReturnNotificationMessage
				if err := json.Unmarshal(d.Body, &message); err != nil {
					log.Printf("Failed to parse return notification message: %v", err)
					d.Nack(false, true)
					log.Println("Retry later...")
					logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to parse return notification message: %v", err), extra, err)
					continue
				}
				requestID = message.RequestID
				extra = map[string]interface{}{
					"email":      message.Email,
					"book_title": message.Book,
				}
				err = mailerService.SendReturnNotification(message.Email, message.Book)
			}

			if err != nil {
				logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to send email for queue %s: %v", queueName, err), extra, err)
				log.Printf("Failed to send email for queue %s: %v", queueName, err)
				d.Nack(false, true)
				log.Println("Retry later...")
			} else {
				logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, fmt.Sprintf("Successfully processed message from exchange:%s queue:%s", constants.EmailExchange, queueName), extra, nil)
				log.Printf("Successfully processed message from queue %s", queueName)
				d.Ack(false)
			}

		}
	}
}
