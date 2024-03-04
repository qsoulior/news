package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type producer struct {
	conn *rabbitmq.Connection

	timeout time.Duration
}

func New(conn *rabbitmq.Connection, opts ...Option) *producer {
	producer := &producer{
		conn:    conn,
		timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(producer)
	}
	return producer
}

func (p *producer) Produce(exchange string, routingKey string, msg amqp.Publishing) error {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	if p.timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	err := p.conn.Ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		msg,
	)

	if err != nil {
		return fmt.Errorf("p.conn.ch.PublishWithContext: %w", err)
	}

	return nil
}
