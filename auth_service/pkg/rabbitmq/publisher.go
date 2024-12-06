package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel   *amqp.Channel
	exchanges map[string]string
}

// NewPublisher initializes a new publisher with the provided RabbitMQ connection.
func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{
		channel:   ch,
		exchanges: make(map[string]string),
	}, nil
}

// DeclareExchange declares an exchange with the specified name and type.
func (p *Publisher) DeclareExchange(name, exchangeType string) error {
	err := p.channel.ExchangeDeclare(
		name,         // Exchange name
		exchangeType, // Exchange type
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return err
	}
	p.exchanges[name] = exchangeType
	log.Printf("Exchange declared: %s (%s)", name, exchangeType)
	return nil
}

// Publish sends a message to the specified exchange with the given routing key.
func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
	if _, exists := p.exchanges[exchange]; !exists {
		return fmt.Errorf("exchange %s is not declared", exchange)
	}

	err := p.channel.Publish(
		exchange,   // Exchange
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to exchange %s with routing key %s: %v", exchange, routingKey, err)
		return err
	}
	log.Printf("Message published to exchange %s with routing key %s", exchange, routingKey)
	return nil
}

// Close closes the publisher channel.
func (p *Publisher) Close() {
	p.channel.Close()
}
