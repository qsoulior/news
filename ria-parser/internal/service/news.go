package service

import (
	"context"
	"strconv"
	"time"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsURL struct {
	URL         string
	PublishedAt time.Time
}

type news interface {
	parseURLs(ctx context.Context, query string, page string) ([]string, error)
}

type newsAbstract struct {
	news
	baseAPI string
	appID   string
	client  *httpclient.Client
}

func (n *newsAbstract) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	urls, err := n.parseURLs(ctx, query, page)
	if err != nil {
		return nil, "", err
	}

	news := make([]entity.News, 0, len(urls))
	for _, url := range urls {
		newsItem, err := n.parseOne(ctx, url)
		if err != nil {
			continue
		}
		news = append(news, *newsItem)
	}

	if page == "" {
		return news, "1", nil
	}

	nextPage, err := strconv.Atoi(page)
	if err != nil {
		return news, "0", nil
	}

	return news, strconv.Itoa(nextPage + 1), nil
}

func (n *newsAbstract) parseOne(ctx context.Context, url string) (*entity.News, error) {
	return nil, nil
}
