package service

import (
	"context"

	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsFeed struct {
	*newsAbstract
}

func NewNewsFeed(baseAPI string) *newsFeed {
	client := httpclient.New(
		httpclient.URL(baseAPI),
	)

	abstract := &newsAbstract{
		client: client,
	}

	feed := &newsFeed{
		newsAbstract: abstract,
	}

	abstract.news = feed
	return feed
}

func (n *newsFeed) parseURLs(ctx context.Context, query string, page string) ([]string, error) {
	return nil, nil
}
