package broker

import (
	"io"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IBroker interface {
	io.Closer

	Publish(payload []byte) error
	Consume() (<-chan amqp.Delivery, error)
}
