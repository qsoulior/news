package repo

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/newsdata-parser/pkg/redis"
	redislib "github.com/redis/go-redis/v9"
)

type pageRedis struct {
	*redis.Redis
}

func NewPageRedis(redis *redis.Redis) Page {
	return &pageRedis{redis}
}

func (p *pageRedis) Get(ctx context.Context) (string, error) {
	page, err := p.Client.Get(ctx, "page").Result()
	if err != nil {
		if err == redislib.Nil {
			return "", ErrNotExist
		}

		return "", fmt.Errorf("p.Client.Get: %w", err)
	}

	return page, nil
}

func (p *pageRedis) Update(ctx context.Context, value string) error {
	err := p.Client.Set(ctx, "page", value, 0).Err()
	if err != nil {
		return fmt.Errorf("p.Client.Set: %w", err)
	}

	return nil
}
