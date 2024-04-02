package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

const PAGE_LAYOUT = "20060102"

type newsArchive struct {
	*news
	*newsView
}

func NewNewsArchive(appID string, client *httpclient.Client, url string, browser *rod.Browser) *newsArchive {
	news := &news{
		appID:  appID,
		client: client,
	}

	newsView := &newsView{
		URL:     url,
		browser: browser,
	}

	archive := &newsArchive{
		news:     news,
		newsView: newsView,
	}

	return archive
}

func (n *newsArchive) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
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

	urls, err := n.parseURLs(ctx, "/"+page)
	if err != nil {
		return nil, "", err
	}

	news, err := n.parseMany(ctx, urls)
	if err != nil {
		return nil, "", err
	}

	nextPage := pageObj.AddDate(0, 0, -1).Format(PAGE_LAYOUT)
	return news, nextPage, nil
}
