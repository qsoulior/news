package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	*Config
	Ch    *amqp.Channel
	conn  *amqp.Connection
	errCh chan *amqp.Error
}

func New(ctx context.Context, cfg *Config) (*Connection, error) {
	c := &Connection{Config: cfg}

	err := c.open(ctx)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) observe(ctx context.Context) {
	go func(ctx context.Context) {
		c.Logger.Info().Msg("observer started")
		select {
		case <-ctx.Done():
			return
		case <-c.errCh:
			c.Close()
			c.reopen(ctx)
		}
	}(ctx)
}

func (c *Connection) open(ctx context.Context) error {
	if _, err := amqp.ParseURI(c.URL); err != nil {
		return fmt.Errorf("amqp.ParseURI: %w", err)
	}

	var err error
	timer := time.NewTimer(0)
	for i := c.AttemptCount; i > 0; i-- {
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			if err = c.attemptOpen(); err == nil {
				defer c.observe(ctx)
				return nil
			}

			c.Logger.Error().
				Err(err).
				Int("left", i).
				Dur("delay", c.AttemptDelay).
				Msg("attempt to establish a connection")

			timer.Reset(c.AttemptDelay)
		}
	}

	return err
}

func (c *Connection) reopen(ctx context.Context) error {
	for timer := time.NewTimer(0); ; timer.Reset(c.AttemptDelay) {
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			err := c.attemptOpen()
			if err == nil {
				defer c.observe(ctx)
				return nil
			}

			c.Logger.Error().
				Err(err).
				Dur("delay", c.AttemptDelay).
				Msg("attempt to establish a connection")
		}
	}
}

func (c *Connection) attemptOpen() error {
	var err error

	c.conn, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.Ch, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("c.conn.Channel: %w", err)
	}

	c.errCh = c.Ch.NotifyClose(make(chan *amqp.Error, 1))
	c.Logger.Info().Msg("connection established")
	return nil
}

func (c *Connection) Close() error {
	// if c.Ch != nil && !c.Ch.IsClosed() {
	// 	if err := c.Ch.Close(); err != nil {
	// 		return fmt.Errorf("c.Ch.Close: %w", err)
	// 	}
	// }

	if c.conn != nil && !c.conn.IsClosed() {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("c.conn.Close: %w", err)
		}
	}

	return nil
}
