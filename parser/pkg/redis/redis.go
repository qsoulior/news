package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Redis struct {
	Client *redis.Client
	logger *zerolog.Logger
}

func New(ctx context.Context, cfg *RedisConfig) (*Redis, error) {
	r := &Redis{logger: zerolog.Ctx(ctx)}

	err := r.open(cfg)
	if err != nil {
		return nil, err
	}

	err = r.ping(ctx, cfg)
	if err != nil {
		if err := r.Close(); err != nil {
			return nil, err
		}

		return nil, err
	}

	return r, nil
}

func (r *Redis) open(cfg *RedisConfig) error {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return fmt.Errorf("redis.ParseURL: %w", err)
	}

	r.Client = redis.NewClient(opt)
	return nil
}

func (r *Redis) ping(ctx context.Context, cfg *RedisConfig) error {
	var err error
	timer := time.NewTimer(0)
	for i := cfg.AttemptCount; i > 0; i-- {
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			if err = r.Client.Ping(ctx).Err(); err == nil {
				return nil
			}

			r.logger.Error().
				Err(err).
				Int("left", i).
				Dur("delay", cfg.AttemptDelay).
				Msg("attempt to establish a connection")

			timer.Reset(cfg.AttemptDelay)
		}
	}

	return err
}

func (r *Redis) Close() error {
	if r.Client != nil {
		err := r.Client.Close()
		if err != nil {
			return fmt.Errorf("r.Client.Close: %w", err)
		}
	}

	return nil
}
