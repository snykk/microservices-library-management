package consumer

import (
	"encoding/json"
	"log"
	"mailer_service/internal/mailer"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OTPMessage struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type LoanNotificationMessage struct {
	Email string `json:"email"`
	Book  string `json:"book"`
	Due   string `json:"due"`
}

type ReturnNotificationMessage struct {
	Email string `json:"email"`
	Book  string `json:"book"`
}

func StartConsuming(ch *amqp.Channel, mailerService mailer.MailerService) error {
	// Declare exchange (e.g., direct exchange)
	err := ch.ExchangeDeclare(
		"email_exchange", // Exchange name
		"direct",         // Exchange type
		true,             // Durable
		false,            // Auto-deleted
		false,            // Internal
		false,            // No-wait
		nil,              // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Declare queues and bind them to the exchange
	queues := []string{"otp_code", "loan_notification", "return_notification"}
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
			queueName,        // Queue name
			queueName,        // Routing key (same as the queue name for direct exchange)
			"email_exchange", // Exchange name
			false,            // No-wait
			nil,              // Arguments
		)
		if err != nil {
			log.Fatalf("Failed to bind queue %s to exchange: %v", queueName, err)
		}
	}

	// Consume all queues
	for _, queueName := range queues {
		go consumeQueue(ch, queueName, mailerService)
	}

	// Wait for consumer processes
	log.Println("Waiting for messages...")
	select {}
}

func consumeQueue(ch *amqp.Channel, queueName string, mailerService mailer.MailerService) {
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

	for d := range msgs {
		switch queueName {
		case "otp_code":
			var message OTPMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to parse OTP message: %v", err)
				d.Nack(false, true)
				continue
			}
			err = mailerService.SendOTP(message.Email, message.OTP)

		case "loan_notification":
			var message LoanNotificationMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to parse loan notification message: %v", err)
				d.Nack(false, true)
				continue
			}
			err = mailerService.SendLoanNotification(message.Email, message.Book, message.Due)

		case "return_notification":
			var message ReturnNotificationMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to parse return notification message: %v", err)
				d.Nack(false, true)
				continue
			}
			err = mailerService.SendReturnNotification(message.Email, message.Book)
		}

		if err != nil {
			log.Printf("Failed to send email for queue %s: %v", queueName, err)
			d.Nack(false, true)
		} else {
			log.Printf("Successfully processed message from queue %s", queueName)
			d.Ack(false)
		}
	}
}
