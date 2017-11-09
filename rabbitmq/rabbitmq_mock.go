package rabbitmq

import (
	"github.com/streadway/amqp"
)

// MockAMQPChannel is a mock from amqp.Channel
type MockAMQPChannel struct{}

// Send mocks Send method
func (channel *MockAMQPChannel) Send(name string, body []byte) error {
	return nil
}

// Receive mocks Receive method
func (channel *MockAMQPChannel) Receive(name string) (<-chan amqp.Delivery, error) {
	return nil, nil
}
