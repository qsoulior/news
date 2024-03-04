package rabbitmq

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL          string
	AttemptCount int
	AttemptDelay time.Duration
}

type Connection struct {
	Config
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func New(cfg Config) (*Connection, error) {
	_, err := amqp.ParseURI(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("amqp.ParseURI: %w", err)
	}

	return &Connection{Config: cfg}, nil
}

func (c *Connection) Open() error {
	var err error
	for i := c.AttemptCount; i > 0; i-- {
		if err = c.attemptOpen(); err == nil {
			return nil
		}

		time.Sleep(c.AttemptDelay)
	}

	if err != nil {
		return fmt.Errorf("c.attemptOpen: %w", err)
	}

	return nil
}

func (c *Connection) attemptOpen() error {
	var err error

	c.Conn, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.Ch, err = c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("c.conn.Channel: %w", err)
	}

	return nil
}

func (c *Connection) Close() error {
	var err error
	if c.Ch != nil {
		err = c.Ch.Close()
		if err != nil {
			return fmt.Errorf("c.ch.Close: %w", err)
		}
	}

	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			return fmt.Errorf("c.conn.Close: %w", err)
		}
	}

	return nil
}
