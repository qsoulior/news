package service

import "github.com/qsoulior/news/aggregator/entity"

type News interface {
	Get() []entity.News
	Create(news *entity.News) error
	CreateMany(news []entity.News) error
}
