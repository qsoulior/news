package repo

import (
	"context"

	"github.com/qsoulior/news/aggregator/entity"
)

type News interface {
	Create(ctx context.Context, news entity.News) error
	ReplaceOrCreate(ctx context.Context, news entity.News) error
	CreateMany(ctx context.Context, news []entity.News) error
	GetByID(ctx context.Context, id string) (*entity.News, error)
	GetByQuery(ctx context.Context, query Query, opts Options) ([]entity.News, error)
}

type Options struct {
	Limit int
	Skip  int
}

type Query struct {
	Title  string
	Source string
}
