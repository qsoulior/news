package service

import (
	"context"

	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsSearch struct {
	*newsAbstract
}

func NewNewsSearch(baseAPI string, appID string) *newsSearch {
	client := httpclient.New()

	abstract := &newsAbstract{
		client:  client,
		baseAPI: baseAPI,
		appID:   appID,
	}

	search := &newsSearch{
		newsAbstract: abstract,
	}

	abstract.news = search
	return search
}

func (n *newsSearch) parseURLs(ctx context.Context, query string, page string) ([]*newsURL, error) {
	return nil, nil
}
