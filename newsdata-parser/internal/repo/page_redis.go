package repo

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/newsdata-parser/pkg/redis"
)

type pageRedis struct {
	*redis.Redis
}

func NewPageRedis(redis *redis.Redis) Page {
	return &pageRedis{redis}
}

func (p *pageRedis) Update(ctx context.Context, value string) error {
	err := p.Client.Set(ctx, "page", value, 0).Err()
	if err != nil {
		return fmt.Errorf("p.Client.Set: %w", err)
	}

	return nil
}
