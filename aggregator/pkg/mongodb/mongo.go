package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	logger *zerolog.Logger
}

func New(ctx context.Context, cfg *Config) (*Mongo, error) {
	m := &Mongo{logger: cfg.Logger}

	err := m.connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = m.ping(ctx, cfg)
	if err != nil {
		if err := m.Disconnect(ctx); err != nil {
			return nil, err
		}

		return nil, err
	}

	return m, nil
}

func (m *Mongo) connect(ctx context.Context, cfg *Config) error {
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetBSONOptions(&options.BSONOptions{
			NilSliceAsEmpty: true,
		})

	err := opts.Validate()
	if err != nil {
		return fmt.Errorf("opts.Validate: %w", err)
	}

	m.Client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("mongo.Connect: %w", err)
	}

	return nil
}

func (m *Mongo) ping(ctx context.Context, cfg *Config) error {
	var err error
	timer := time.NewTimer(0)
	for i := cfg.AttemptCount; i > 0; i-- {
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			if err = m.Client.Ping(ctx, nil); err == nil {
				return nil
			}

			m.logger.Error().
				Err(err).
				Int("left", i).
				Dur("delay", cfg.AttemptDelay).
				Msg("attempt to establish a connection")

			timer.Reset(cfg.AttemptDelay)
		}
	}

	return err
}

func (m *Mongo) Disconnect(ctx context.Context) error {
	if m.Client != nil {
		err := m.Client.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("m.Client.Disconnect: %w", err)
		}
	}

	return nil
}
