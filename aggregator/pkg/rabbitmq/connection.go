package rabbitmq

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Connection struct {
	Ch     *amqp.Channel
	conn   *amqp.Connection
	logger *zerolog.Logger
}

func New(cfg *Config) (*Connection, error) {
	c := &Connection{logger: cfg.Logger}

	err := c.open(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) open(cfg *Config) error {
	_, err := amqp.ParseURI(cfg.URL)
	if err != nil {
		return fmt.Errorf("amqp.ParseURI: %w", err)
	}

	for i := cfg.AttemptCount; i > 0; i-- {
		if err = c.attemptOpen(cfg); err == nil {
			return nil
		}

		c.logger.Error().
			Err(err).
			Int("left", i).
			Dur("delay", cfg.AttemptDelay).
			Msg("attempt to establish a connection")

		time.Sleep(cfg.AttemptDelay)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) attemptOpen(cfg *Config) error {
	var err error

	c.conn, err = amqp.Dial(cfg.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.Ch, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("c.conn.Channel: %w", err)
	}

	return nil
}

func (c *Connection) Err() <-chan *amqp.Error {
	return c.conn.NotifyClose(make(chan *amqp.Error, 1))
}

func (c *Connection) Close() error {
	var err error
	if c.Ch != nil {
		err = c.Ch.Close()
		if err != nil {
			return fmt.Errorf("c.Ch.Close: %w", err)
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
