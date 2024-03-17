package service

import "context"

type News interface {
	Parse(ctx context.Context, query string, page string) (string, error)
	Release(ctx context.Context) (int, error)
}

type Page interface {
	Get(ctx context.Context) (string, error)
	Set(ctx context.Context, page string) error
}
