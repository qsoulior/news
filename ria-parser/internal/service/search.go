package service

import (
	"context"
	"strconv"

	"github.com/go-rod/rod"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsSearch struct {
	*news
	browser *rod.Browser
}

func NewNewsSearch(baseAPI string, appID string, browser *rod.Browser) *newsSearch {
	client := httpclient.New()

	news := &news{
		client:  client,
		baseAPI: baseAPI,
		appID:   appID,
	}

	search := &newsSearch{
		news:    news,
		browser: browser,
	}

	return search
}

func (n *newsSearch) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
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

func (n *newsSearch) parseURLs(ctx context.Context, query string, page string) ([]string, error) {
	return nil, nil
}
