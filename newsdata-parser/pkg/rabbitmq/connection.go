package rabbitmq

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionConfig struct {
	URL          string
	AttemptCount int
	AttemptDelay time.Duration
}

type Connection struct {
	ConnectionConfig
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConnection(cfg ConnectionConfig) (*Connection, error) {
	_, err := amqp.ParseURI(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("amqp.ParseURI: %w", err)
	}

	return &Connection{ConnectionConfig: cfg}, nil
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

	c.conn, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("c.conn.Channel: %w", err)
	}

	return nil
}

func (c *Connection) Close() error {
	var err error
	if c.ch != nil {
		err = c.ch.Close()
		if err != nil {
			return fmt.Errorf("c.ch.Close: %w", err)
		}
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return fmt.Errorf("c.conn.Close: %w", err)
		}
	}

	return nil
}
