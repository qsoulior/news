package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Consume(queue string) error
}

type Producer interface {
	Produce(ctx context.Context, exchange string, routingKey string, msg Message) error
}

type Message = amqp.Publishing
type Delivery = amqp.Delivery
type Handler func(ctx context.Context, msg *Delivery)
