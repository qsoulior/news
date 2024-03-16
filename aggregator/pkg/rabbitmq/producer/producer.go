package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
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

func (p *producer) Produce(ctx context.Context, exchange string, routingKey string, msg rabbitmq.Message) error {
	var (
		timeoutCtx context.Context
		cancel     context.CancelFunc
	)

	if p.timeout > 0 {
		timeoutCtx, cancel = context.WithTimeout(ctx, p.timeout)
	} else {
		timeoutCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	err := p.conn.Ch.PublishWithContext(
		timeoutCtx,
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
