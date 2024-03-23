package service

import (
	"context"
)

type newsFeed struct {
	*newsAbstract
}

func NewNewsFeed(baseAPI string, appID string) *newsFeed {
	abstract := &newsAbstract{
		baseAPI: baseAPI,
		appID:   appID,
	}

	feed := &newsFeed{
		newsAbstract: abstract,
	}

	abstract.news = feed
	return feed
}

func (n *newsFeed) parseURLs(ctx context.Context, query string, page string) ([]*newsURL, error) {
	return nil, nil
}
