package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	URL          string
	AttemptCount int
	AttemptDelay time.Duration
}

type Redis struct {
	RedisConfig
	Client *redis.Client
}

func New(cfg RedisConfig) (*Redis, error) {
	_, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("redis.ParseURL: %w", err)
	}

	return &Redis{RedisConfig: cfg}, nil
}

func (r *Redis) Open(ctx context.Context) error {
	opt, _ := redis.ParseURL(r.URL)
	r.Client = redis.NewClient(opt)

	var err error
	for i := r.AttemptCount; i > 0; i-- {
		if err := r.Client.Ping(ctx).Err(); err == nil {
			return nil
		}

		time.Sleep(r.AttemptDelay)
	}

	if err != nil {
		return fmt.Errorf("r.client.Ping: %w", err)
	}

	return nil
}

func (r *Redis) Close() error {
	if r.Client != nil {
		err := r.Client.Close()
		if err != nil {
			return fmt.Errorf("r.client.Close: %w", err)
		}
	}

	return nil
}
