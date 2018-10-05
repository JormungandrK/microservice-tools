package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Channel sends and receives the messages from the queus
type Channel interface {
	// Send sends body to a message queue.
	Send(name string, body []byte) error
	// Receive immediately starts delivering queued messages.
	Receive(name string) (<-chan amqp.Delivery, error)
	// SendToExchange sends body to an exchange.
	SendToExchange(name string, excType string, body []byte) error
	// Receive immediately starts delivering queued messages.
	ReceiveOnExchange(name string, excType string) (<-chan amqp.Delivery, error)
}
