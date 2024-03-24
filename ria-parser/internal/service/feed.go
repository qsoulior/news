package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsFeed struct {
	*news
	browser *rod.Browser
}

func NewNewsFeed(baseAPI string, appID string, browser *rod.Browser) *newsFeed {
	client := httpclient.New()

	news := &news{
		client:  client,
		baseAPI: baseAPI,
		appID:   appID,
	}

	feed := &newsFeed{
		news:    news,
		browser: browser,
	}

	return feed
}

const PAGE_LAYOUT = "20060102"

func (n *newsFeed) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	var (
		pageObj time.Time
		err     error
	)

	if page == "" {
		pageObj = time.Now()
		page = pageObj.Format(PAGE_LAYOUT)
	} else {
		pageObj, err = time.Parse(PAGE_LAYOUT, page)
		if err != nil {
			return nil, "", fmt.Errorf("time.Parse: %w", err)
		}
	}

	urls, err := n.parseURLs(ctx, page)
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

	nextPage := pageObj.AddDate(0, 0, -1).Format(PAGE_LAYOUT)
	return news, nextPage, nil
}

func (n *newsFeed) parseURLs(ctx context.Context, path string) ([]string, error) {
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

	url, err := url.JoinPath(n.baseAPI, path)
	if err != nil {
		return nil, fmt.Errorf("utl.JoinPath: %w", err)
	}

	err = page.Navigate(url)
	if err != nil {
		return nil, fmt.Errorf("n.page.Navigate: %w", err)
	}

	urls, err := n.parseView(page, 0)
	if err != nil {
		return nil, fmt.Errorf("n.parseView: %w", err)
	}

	return urls, nil
}
