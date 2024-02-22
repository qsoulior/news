package repo

import (
	"context"
)

type News interface {
	Create(ctx context.Context, jsonStr string) error
	Pop(ctx context.Context) (string, error)
	PopAll(ctx context.Context) ([]string, error)
}

type Page interface {
	Update(ctx context.Context, value string) error
}
