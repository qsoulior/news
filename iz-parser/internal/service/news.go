package service

import "github.com/qsoulior/news/aggregator/entity"

type news struct {
}

func NewNews() *news {
	return &news{}
}

func (n *news) Parse(query string, page string) ([]entity.News, string, error)
