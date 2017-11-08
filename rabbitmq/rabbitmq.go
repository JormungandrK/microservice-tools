package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Dial returns a new Connection over TCP using PlainAuth
func Dial(username, password, host, port string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			username,
			password,
			host,
			port,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	return conn, ch, err
}

// Publish sends body to an exchange on the server.
func Publish(name string, body []byte, channel *amqp.Channel) error {
	q, err := channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	return err
}

// Consume immediately starts delivering queued messages.
func Consume(name string, channel *amqp.Channel) (<-chan amqp.Delivery, error) {
	q, err := channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil
	}

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, nil
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return msgs, err
}
