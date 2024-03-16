package service

import "context"

type News interface {
	Parse(ctx context.Context, query string, page string) (string, error)
}

type Page interface {
	Get() (string, error)
	Set(page string) error
}
