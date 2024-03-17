package repo

import (
	"context"
)

type News interface {
	Create(ctx context.Context, jsonStr string) error

	Get(ctx context.Context, index int) (string, error)
	GetFirst(ctx context.Context) (string, error)
	GetLast(ctx context.Context) (string, error)
	GetAll(ctx context.Context) ([]string, error)

	PopFirst(ctx context.Context) (string, error)
	PopLast(ctx context.Context) (string, error)
	PopAll(ctx context.Context) ([]string, error)

	DeleteFirst(ctx context.Context) error
	DeleteLast(ctx context.Context) error
	DeleteAll(ctx context.Context) error
}

type Page interface {
	Get(ctx context.Context) (string, error)
	Update(ctx context.Context, value string) error
}
