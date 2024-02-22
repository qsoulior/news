package repo

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/newsdata-parser/pkg/redis"
)

type newsRedis struct {
	*redis.Redis
}

func NewNewsRedis(redis *redis.Redis) News {
	return &newsRedis{redis}
}

func (n *newsRedis) Create(ctx context.Context, jsonStr string) error {
	err := n.Client.LPush(ctx, "news", jsonStr).Err()
	if err != nil {
		return fmt.Errorf("n.Client.LPush: %w", err)
	}
	return nil
}

func (n *newsRedis) Pop(ctx context.Context) (string, error) {
	jsonStr, err := n.Client.LPop(ctx, "news").Result()
	if err != nil {
		return "", fmt.Errorf("n.Client.Pop: %w", err)
	}
	return jsonStr, nil
}

func (n *newsRedis) PopAll(ctx context.Context) ([]string, error) {
	pipe := n.Client.TxPipeline()
	cmd := pipe.LRange(ctx, "news", 0, -1)
	pipe.Del(ctx, "news")

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("pipe.Exec: %w", err)
	}

	return cmd.Val(), nil
}
