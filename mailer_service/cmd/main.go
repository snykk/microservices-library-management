package main

import (
	"io/ioutil"
	"log"
	"mailer_service/internal/consumer"
	"mailer_service/internal/mailer"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	// Email service env
	emailSenderBytes, err := ioutil.ReadFile(os.Getenv("EMAIL_SENDER_CONTAINER_FILE"))
	if err != nil {
		log.Fatalf("Error reading email sender secret: %v", err)
	}

	emailPasswordBytes, err := ioutil.ReadFile(os.Getenv("EMAIL_PASSWORD_CONTAINER_FILE"))
	if err != nil {
		log.Fatalf("Error reading email password secret: %v", err)
	}

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Mailer Service
	mailerService, err := mailer.NewMailerService(string(emailSenderBytes), string(emailPasswordBytes))
	if err != nil {
		log.Fatalf("Failed to create mailer service %v", err)
	}

	err = consumer.StartConsuming(ch, mailerService)
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
}
