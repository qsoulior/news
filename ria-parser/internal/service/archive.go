package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/queue"
	"github.com/rs/zerolog"
)

const PAGE_LAYOUT = "20060102"

type newsArchive struct {
	*news
	*newsView
	urls queue.Queue[string]
}

func NewNewsArchive(
	appID string,
	client *httpclient.Client,
	url string,
	browser *rod.Browser,
	logger *zerolog.Logger,
) *newsArchive {
	log := logger.With().Str("service", "archive").Logger()

	news := &news{
		appID:  appID,
		client: client,
		logger: &log,
	}

	newsView := &newsView{
		URL:     url,
		browser: browser,
	}

	archive := &newsArchive{
		news:     news,
		newsView: newsView,
		urls:     queue.New[string](100),
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

	const limit = 20
	urls := make([]string, 0, limit)
	nextPage := page

	if n.urls.Len() > 0 {
		for i := 0; i < limit; i++ {
			url, ok := n.urls.Pop()
			if !ok {
				break
			}

			urls = append(urls, url)
		}

		if n.urls.Len() == 0 {
			nextPage = pageObj.AddDate(0, 0, -1).Format(PAGE_LAYOUT)
		}
	} else {
		parsedUrls, err := n.parseURLs(ctx, "/"+page)
		if err != nil {
			return nil, "", err
		}

		parsedLen := len(parsedUrls)
		urls = append(urls, parsedUrls[:min(parsedLen, limit)]...)
		if parsedLen > limit {
			n.urls.Push(parsedUrls[limit:]...)
		}
	}

	news, err := n.parseMany(ctx, urls)
	if err != nil {
		return nil, "", fmt.Errorf("n.parseMany: %w", err)
	}

	return news, nextPage, nil
}
