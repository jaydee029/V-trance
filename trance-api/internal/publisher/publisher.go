package publisher

import (
	"fmt"

	"github.com/jaydee029/V-trance/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func New(conn *amqp.Connection) *PbClient {
	pb := pubsub.New(conn)

	return &PbClient{
		pubsub: pb,
	}
}

type PbClient struct {
	pubsub *pubsub.PubSub
}

func (pb *PbClient) PublishTask(exchange, key string, val pubsub.Task, logger *zap.Logger) error {
	err := pb.pubsub.PublishGob(exchange, key, val)
	if err != nil {
		logger.Info("Error publish the event:", zap.Error(err))
		return err
	}
	return nil
}

func InitBroker(conn *amqp.Connection, exchange string) error {

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error creating a channel: %w", err)
	}

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error declaring %s: %w", exchange, err)
	}

	return nil
}
