package service

import (
	"context"
	"time"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsURL struct {
	URL         string
	PublishedAt time.Time
}

type news struct {
	baseAPI string
	appID   string
	client  *httpclient.Client
}

func (n *news) parseOne(ctx context.Context, url string) (*entity.News, error) {
	return nil, nil
}
