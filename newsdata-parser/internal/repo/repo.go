package repo

import "context"

type News interface {
	Create(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) error
}
