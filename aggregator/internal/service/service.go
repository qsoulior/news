package service

import (
	"context"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/internal/repo"
)

type News interface {
	Create(ctx context.Context, news entity.News) error
	CreateMany(ctx context.Context, news []entity.News) error
	Get(ctx context.Context, id string) (*entity.News, error)
	GetHead(ctx context.Context, query repo.Query, opts Options) ([]entity.NewsHead, int, error)
	SendToParse(ctx context.Context, query string) error
}

type (
	Query = repo.Query
)
