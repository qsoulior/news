package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
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
	urls, err := n.parseURLs(ctx, query)
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

	return news, "", nil
}

func (n *newsSearch) parseURLs(ctx context.Context, query string) ([]string, error) {
	page, err := stealth.Page(n.browser)
	if err != nil {
		return nil, fmt.Errorf("stealth.Page: %w", err)
	}
	defer page.Close()

	page = page.Context(ctx)

	err = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: gofakeit.UserAgent()})
	if err != nil {
		return nil, fmt.Errorf("n.page.SetUserAgent: %w", err)
	}

	u, _ := url.Parse(n.baseAPI + "/search")
	values := u.Query()
	values.Set("query", query)
	u.RawQuery = values.Encode()

	err = page.Navigate(u.String())
	if err != nil {
		return nil, fmt.Errorf("n.page.Navigate: %w", err)
	}

	urls, err := n.parseView(page, 0)
	if err != nil {
		return nil, fmt.Errorf("n.parseView: %w", err)
	}

	return urls, nil
}
