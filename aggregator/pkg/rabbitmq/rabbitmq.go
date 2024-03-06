package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Consumer interface {
	Consume(queue string) error
}

type Producer interface {
	Produce(exchange string, routingKey string, msg Message) error
}

type Message = amqp.Publishing
type Delivery = amqp.Delivery
type Handler func(msg *Delivery)
