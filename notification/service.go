package notification

import (
	"encoding/json"

	commonCfg "github.com/Microkubes/microservice-tools/config"
	"github.com/Microkubes/microservice-tools/rabbitmq"
	"github.com/streadway/amqp"
)

// EventPayload describes the event object for handling dynamic services
type EventPayload struct {
	// ObjectID is the ID of the object( Organization ID or schema ID or resource ID etc)
	ObjectID string `json:"objectID"`
	// ObjectType is the type of the object("organization", "schema" etc)
	ObjectType string `json:"objectType"`
	// Event is one of "CREATE", "UPDATE" or "DELETE
	Event string `json:"event"`
	// Data holds the additional event data
	Data map[string]interface{} `json:"data"`
}

// ObjectHandler receives message from message queue
type ObjectHandler func(amqpChan <-chan amqp.Delivery)

// AMQPService holds info for AMQP based implementation
type AMQPService struct {
	Channel   *rabbitmq.AMQPChannel
	QueueName string
	handlers  []ObjectHandler
}

// Service defines the interface for sending object to the queue
type Service interface {
	// SendObject sends an event obect to the queue to notify for new object event
	SendObject(e *EventPayload) error
	// ReceiveObjects receives platform object payload
	ReceiveObjects() error
	// AddObjectHandler adds new platform object handler
	AddObjectHandler(handler ObjectHandler)
}

// NewService creates notification service for sending and receiving messages on the queue
func NewService(mqConfig *commonCfg.MQConfig) (*AMQPService, func(), error) {
	conn, ch, err := rabbitmq.Dial(mqConfig.Username, mqConfig.Password, mqConfig.Host, mqConfig.Port)
	if err != nil {
		return nil, nil, err
	}

	channelWrapper := rabbitmq.AMQPChannel{ch}
	AMQPService := AMQPService{
		Channel: &channelWrapper,
	}

	return &AMQPService, func() {
		conn.Close()
		ch.Close()
	}, nil
}

// SetQueueName set the queue name for the notification service
func (a *AMQPService) SetQueueName(queueName string) {
	a.QueueName = queueName
}

// GetQueueName return the queue name of notification service
func (a *AMQPService) GetQueueName(queueName string) string {
	return a.QueueName
}

// SendObject sends an event object for handling platform objects to the queue
func (a *AMQPService) SendObject(e *EventPayload) error {
	paylaod, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return a.Channel.Send(a.QueueName, paylaod)
}

// ReceiveObjects receives platorm object payload
// It calls notification service handlers
func (a *AMQPService) ReceiveObjects() error {
	ampqChan, err := a.Channel.Receive(a.QueueName)
	if err != nil {
		return err
	}

	for _, handler := range a.handlers {
		go func(h ObjectHandler) {
			h(ampqChan)
		}(handler)
	}

	return nil
}

// AddObjectHandler adds new Platform Object Handler. Handlers are called when ReceiveObjects is execute
func (a *AMQPService) AddObjectHandler(handler ObjectHandler) {
	a.handlers = append(a.handlers, handler)
}
