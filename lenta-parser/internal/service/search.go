package service

import (
	"context"

	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsSearch struct {
	*newsAbstract
}

func NewNewsSearch(baseAPI string) *newsSearch {
	client := httpclient.New(
		httpclient.URL(baseAPI),
	)

	abstract := &newsAbstract{
		client: client,
	}

	search := &newsSearch{
		newsAbstract: abstract,
	}

	abstract.news = search
	return search
}

func (n *newsSearch) parseURLs(ctx context.Context, query string, page string) ([]string, error) {
	return nil, nil
}
