package main

import (
	"io/ioutil"
	"log"
	"mailer_service/configs"
	"mailer_service/internal/consumer"
	"mailer_service/internal/mailer"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	// Load app config
	if err := configs.InitializeAppConfig(); err != nil {
		log.Fatal("Failed to load app config", err)
	}
	log.Println("App configuration loaded")
}

func main() {
	// Email service env
	emailSenderBytes, err := ioutil.ReadFile(configs.AppConfig.EmailSenderContainerFile)
	if err != nil {
		log.Fatalf("Error reading email sender secret: %v", err)
	}

	emailPasswordBytes, err := ioutil.ReadFile(configs.AppConfig.EmailPasswordContainerFile)
	if err != nil {
		log.Fatalf("Error reading email password secret: %v", err)
	}

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(configs.AppConfig.RabbitMQURL)
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
