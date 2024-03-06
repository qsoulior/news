package consumer

import (
	"fmt"

	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
)

type consumer struct {
	conn    *rabbitmq.Connection
	handler rabbitmq.Handler

	autoAck bool
}

func New(conn *rabbitmq.Connection, handler rabbitmq.Handler, opts ...Option) *consumer {
	consumer := &consumer{
		conn:    conn,
		handler: handler,
		autoAck: true,
	}

	for _, opt := range opts {
		opt(consumer)
	}

	return consumer
}

func (c *consumer) Consume(queue string) error {
	msgs, err := c.conn.Ch.Consume(
		queue,
		"",
		c.autoAck,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("c.conn.ch.Consume: %w", err)
	}

	for msg := range msgs {
		c.handler(&msg)
		if !c.autoAck {
			msg.Ack(false)
		}
	}

	return nil
}
