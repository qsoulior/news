package rabbitmq

import (
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

func New(cfg *Config) (*Connection, error) {
	c := &Connection{Config: cfg}

	err := c.open()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) observe() {
	go func() {
		if err := <-c.errCh; err != nil {
			c.Close()
			c.reopen()
		}
	}()
}

func (c *Connection) open() error {
	_, err := amqp.ParseURI(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.ParseURI: %w", err)
	}

	for i := c.AttemptCount; i > 0; i-- {
		if err = c.attemptOpen(); err == nil {
			return nil
		}

		c.Logger.Error().
			Err(err).
			Int("left", i).
			Dur("delay", c.AttemptDelay).
			Msg("attempt to establish a connection")

		time.Sleep(c.AttemptDelay)
	}

	if err != nil {
		return err
	}

	defer c.observe()
	return nil
}

func (c *Connection) reopen() {
	for {
		if err := c.attemptOpen(); err == nil {
			defer c.observe()
			return
		}
		time.Sleep(1 * time.Minute)
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
	return nil
}

func (c *Connection) Close() error {
	if c.Ch != nil && !c.Ch.IsClosed() {
		if err := c.Ch.Close(); err != nil {
			return fmt.Errorf("c.Ch.Close: %w", err)
		}
	}

	if c.conn != nil && !c.conn.IsClosed() {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("c.conn.Close: %w", err)
		}
	}

	return nil
}
