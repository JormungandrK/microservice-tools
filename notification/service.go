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
	// ErrorMessage is the actual error if it happens.
	// Event payload in thi case is published on error queue
	ErrorMessage string `json:"errorMessage"`
}

// ObjectHandler receives message from message queue
type ObjectHandler func(amqpChan <-chan amqp.Delivery)

// AMQPService holds info for AMQP based implementation
type AMQPService struct {
	Channel      *rabbitmq.AMQPChannel
	ExchnageName string
	ExchangeType string
	handlers     []ObjectHandler
}

// Service defines the interface for sending object to the queue
type Service interface {
	// SendObject sends an event obect to the queue to notify for new object event
	SendObject(e *EventPayload) error
	// ReceiveObjects receives platform object payload
	ReceiveObjects() error
	// AddObjectHandler adds new platform object handler
	AddObjectHandler(handler ObjectHandler)
	// SetExchangeName set the exchange name for the notification service
	SetExchangeName(exchnageName string)
	// GetExchangeName return the exchange name of notification service
	GetExchangeName(queueName string) string
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

// SetExchangeName set the exchange name for the notification service
func (a *AMQPService) SetExchangeName(exchnageName string) {
	a.ExchnageName = exchnageName
}

// GetExchangeName return the exchange name of notification service
func (a *AMQPService) GetExchangeName(exchnageName string) string {
	return a.ExchnageName
}

// SetExchnageType sets the exchnage type
func (a *AMQPService) SetExchnageType(exchangeType string) {
	a.ExchangeType = exchangeType
}

// GetExchangeType returns the exchange type
func (a *AMQPService) GetExchangeType(exchnageType string) string {
	return a.ExchangeType
}

// SendObject sends an event object for handling platform objects to the queue
func (a *AMQPService) SendObject(e *EventPayload) error {
	paylaod, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return a.Channel.SendToExchange(a.ExchnageName, a.ExchangeType, paylaod)
}

// ReceiveObjects receives platorm object payload. It calls notification service handlers
func (a *AMQPService) ReceiveObjects() error {
	ampqChan, err := a.Channel.ReceiveOnExchange(a.ExchnageName, a.ExchangeType)
	if err != nil {
		return err
	}

	for _, handler := range a.handlers {
		go handler(ampqChan)
	}

	return nil
}

// AddObjectHandler adds new Platform Object Handler. Handlers are called when ReceiveObjects is execute
func (a *AMQPService) AddObjectHandler(handler ObjectHandler) {
	a.handlers = append(a.handlers, handler)
}
