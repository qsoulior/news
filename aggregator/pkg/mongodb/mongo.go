package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URL          string
	AttemptCount int
	AttemptDelay time.Duration
}

type Mongo struct {
	MongoConfig
	Client *mongo.Client
}

func New(cfg MongoConfig) (*Mongo, error) {
	err := options.Client().ApplyURI(cfg.URL).Validate()
	if err != nil {
		return nil, fmt.Errorf("options.Client.ApplyURI: %w", err)
	}

	return &Mongo{MongoConfig: cfg}, nil
}

func (m *Mongo) Connect(ctx context.Context) error {
	var err error

	opts := options.Client().ApplyURI(m.URL)
	m.Client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("mongo.Connect: %w", err)
	}

	for i := m.AttemptCount; i > 0; i-- {
		if err = m.Client.Ping(ctx, nil); err == nil {
			return nil
		}

		log.Error().
			Err(fmt.Errorf("m.Client.Ping: %w", err)).
			Int("left", i).
			Dur("delay", m.AttemptDelay).
			Msg("")

		time.Sleep(m.AttemptDelay)
	}

	if err != nil {
		return fmt.Errorf("m.Client.Ping: %w", err)
	}

	return nil
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
