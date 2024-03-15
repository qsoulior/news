package service

import "github.com/qsoulior/news/aggregator/entity"

type newsSearch struct {
}

func NewNewsSearch() *newsSearch {
	return &newsSearch{}
}

func (n *newsSearch) Parse(query string, page string) ([]entity.News, string, error)
