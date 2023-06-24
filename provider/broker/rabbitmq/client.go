package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	config  Config
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewClient(config Config) (*Broker, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	conn, err := amqp.Dial(config.RabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %w", err)
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		return nil, fmt.Errorf("qos: %w", err)
	}

	queue, err := channel.QueueDeclare(
		config.RabbitmqNewPostsNtfQueueName, true, false, false, false, nil,
	)

	return &Broker{config: config, conn: conn, channel: channel, queue: queue}, nil
}

func (b Broker) Publish(payload []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.config.RabbitmqPublishDeadlineDur)
	defer cancel()

	return b.channel.PublishWithContext(
		ctx,
		"",
		b.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         payload,
		})
}

func (b Broker) Consume() (<-chan amqp.Delivery, error) {
	return b.channel.Consume(
		b.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (b Broker) Close() error {
	return b.conn.Close()
}
