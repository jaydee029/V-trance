package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

//v1.1.0
type PubSub struct {
	conn *amqp.Connection
}
type PubSubinterface interface {
	PublishGob(exchange, key string, val any) error
	SubscribeGob(exchange, queueName, key string, simpleQueueType QueueType, acktype Acktype, handler func(any) Acktype) error
	DeclareAndBind(exchange, queueName, key string, simpleQueueType QueueType) (*amqp.Channel, amqp.Queue, error)
}

func New(conn *amqp.Connection) *PubSub {
	return &PubSub{conn: conn}
}
