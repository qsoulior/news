package service

import (
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/internal/repo"
)

type News interface {
	Create(news entity.News) error
	CreateMany(news []entity.News) error
	GetByID(id string) (*entity.News, error)
	GetByQuery(query repo.Query, opts repo.Options) ([]entity.News, error)
	Parse(query string) error
}
