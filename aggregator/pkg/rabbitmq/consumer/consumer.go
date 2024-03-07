package consumer

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
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

func (c *consumer) Consume(ctx context.Context, queue string) error {
	msgs, err := c.conn.Ch.ConsumeWithContext(
		ctx,
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

	if ctx.Err() != nil {
		return nil
	}

	return fmt.Errorf("amqp.Delivery is closed: %w", amqp.ErrClosed)
}
