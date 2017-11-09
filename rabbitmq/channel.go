package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Channel sends and receives the messages from the queus
type Channel interface {
	// Send sends body to an exchange on the server.
	Send(name string, body []byte) error
	// Receive immediately starts delivering queued messages.
	Receive(name string) (<-chan amqp.Delivery, error)
}
