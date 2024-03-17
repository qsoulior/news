package repo

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/parser/pkg/redis"
	rdb "github.com/redis/go-redis/v9"
)

type newsRedis struct {
	*redis.Redis
}

func NewNewsRedis(redis *redis.Redis) News {
	return &newsRedis{redis}
}

func (n *newsRedis) Create(ctx context.Context, jsonStr string) error {
	err := n.Client.RPush(ctx, "news", jsonStr).Err()
	if err != nil {
		return fmt.Errorf("n.Client.RPush: %w", err)
	}
	return nil
}

func (n *newsRedis) Get(ctx context.Context, index int) (string, error) {
	jsonStr, err := n.Client.LIndex(ctx, "news", int64(index)).Result()
	if err != nil {
		if err == rdb.Nil {
			return "", ErrNotExist
		}

		return "", fmt.Errorf("n.Client.LIndex: %w", err)
	}

	return jsonStr, nil
}

func (n *newsRedis) GetFirst(ctx context.Context) (string, error) {
	return n.Get(ctx, 0)
}

func (n *newsRedis) GetLast(ctx context.Context) (string, error) {
	return n.Get(ctx, -1)
}

func (n *newsRedis) GetAll(ctx context.Context) ([]string, error) {
	jsonStrs, err := n.Client.LRange(ctx, "news", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("n.Client.LRange: %w", err)
	}

	return jsonStrs, nil
}

func (n *newsRedis) PopFirst(ctx context.Context) (string, error) {
	jsonStr, err := n.Client.LPop(ctx, "news").Result()
	if err != nil {
		if err == rdb.Nil {
			return "", ErrNotExist
		}

		return "", fmt.Errorf("n.Client.LPop: %w", err)
	}
	return jsonStr, nil
}

func (n *newsRedis) PopLast(ctx context.Context) (string, error) {
	jsonStr, err := n.Client.RPop(ctx, "news").Result()
	if err != nil {
		if err == rdb.Nil {
			return "", ErrNotExist
		}

		return "", fmt.Errorf("n.Client.RPop: %w", err)
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

func (n *newsRedis) DeleteFirst(ctx context.Context) error {
	err := n.Client.LPop(ctx, "news").Err()
	if err != nil {
		return fmt.Errorf("n.Client.LPop: %w", err)
	}
	return nil
}

func (n *newsRedis) DeleteLast(ctx context.Context) error {
	err := n.Client.RPop(ctx, "news").Err()
	if err != nil {
		return fmt.Errorf("n.Client.RPop: %w", err)
	}

	return nil
}

func (n *newsRedis) DeleteAll(ctx context.Context) error {
	err := n.Client.Del(ctx, "news").Err()
	if err != nil {
		return fmt.Errorf("n.Client.Del: %w", err)
	}

	return nil
}
