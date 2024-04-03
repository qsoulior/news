package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-rod/rod"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsSearch struct {
	*news
	*newsView
}

func NewNewsSearch(appID string, client *httpclient.Client, url string, browser *rod.Browser) *newsSearch {
	news := &news{
		appID:  appID,
		client: client,
	}

	newsView := &newsView{
		URL:     url,
		browser: browser,
	}

	search := &newsSearch{
		news:     news,
		newsView: newsView,
	}

	return search
}

func (n *newsSearch) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	u, _ := url.Parse("/search")
	values := u.Query()
	values.Set("query", query)
	u.RawQuery = values.Encode()

	urls, err := n.parseURLs(ctx, u.String())
	if err != nil {
		return nil, "", err
	}

	news, err := n.parseMany(ctx, urls)
	if err != nil {
		return nil, "", fmt.Errorf("n.parseMany: %w", err)
	}

	return news, "", nil
}
