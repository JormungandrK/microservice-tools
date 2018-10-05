package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// SendToExchange sends body to an exchange.
func (c *AMQPChannel) SendToExchange(name string, excType string, body []byte) error {
	if err := c.ExchangeDeclare(
		name,    // name
		excType, // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	); err != nil {
		return fmt.Errorf("Failed to declare an exchange: %s", err.Error())
	}

	if err := c.Publish(
		name,  // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	); err != nil {
		return fmt.Errorf("Failed to publish a message: %s", err.Error())
	}

	log.Printf(" [x] Sent %s", body)

	return nil
}

// ReceiveOnExchange immediately starts delivering queued messages.
func (c *AMQPChannel) ReceiveOnExchange(name string, excType string) (<-chan amqp.Delivery, error) {
	if err := c.ExchangeDeclare(
		name,    // name
		excType, // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	); err != nil {
		return nil, fmt.Errorf("Failed to declare an exchange: %s", err.Error())
	}

	q, err := c.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %s", err.Error())
	}

	if err = c.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("Failed to bind a queue: %s", err.Error())
	}

	return c.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}
