package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *Connection
}

func NewConsumer(conn *Connection) *Consumer {
	return &Consumer{conn}
}

type Handler func(msg *amqp.Delivery)

func (c *Consumer) Consume(queue string, handler Handler) error {
	msgs, err := c.conn.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("c.conn.ch.Consume: %w", err)
	}

	for msg := range msgs {
		handler(&msg)
		msg.Ack(false)
	}

	return nil
}
