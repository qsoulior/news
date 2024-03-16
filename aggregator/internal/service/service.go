package service

import (
	"context"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/internal/repo"
)

type News interface {
	Create(ctx context.Context, news entity.News) error
	CreateMany(ctx context.Context, news []entity.News) error
	GetByID(ctx context.Context, id string) (*entity.News, error)
	GetByQuery(ctx context.Context, query repo.Query, opts repo.Options) ([]entity.News, error)
	Parse(ctx context.Context, query string) error
}
