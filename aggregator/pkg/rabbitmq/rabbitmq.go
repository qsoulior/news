package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type Consumer interface {
	Consume(queue string) error
}

type Producer interface {
	Produce(exchange string, routingKey string, msg amqp091.Publishing) error
}
