package repo

import "context"

type newsRedis struct {
}

func NewNewsRedis() News {
	return &newsRedis{}
}

func (n *newsRedis) Create(ctx context.Context, key string, value string) error {
	return nil
}

func (n *newsRedis) Get(ctx context.Context, key string) error {
	return nil
}
