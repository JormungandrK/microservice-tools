package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

// AMQPChannel wraps a amqp.Channel
type AMQPChannel struct {
	*amqp.Channel
}

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
