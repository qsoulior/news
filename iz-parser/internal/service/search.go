package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/DataHenHQ/useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/rs/zerolog"
)

type newsSearch struct {
	*news
}

func NewNewsSearch(appID string, client *httpclient.Client, logger *zerolog.Logger) *newsSearch {
	log := logger.With().Str("service", "search").Logger()
	news := &news{
		appID:  appID,
		client: client,
		logger: &log,
	}

	search := &newsSearch{
		news: news,
	}

	return search
}

func (n *newsSearch) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	urls, err := n.parseURLs(ctx, query, page)
	if err != nil {
		return nil, "", err
	}

	news, err := n.parseMany(ctx, urls)
	if err != nil {
		return nil, "", fmt.Errorf("n.parseMany: %w", err)
	}

	return news, "", nil
}

func (n *newsSearch) parseURLs(ctx context.Context, query string, from string) ([]string, error) {
	u, _ := url.Parse("/search")
	values := u.Query()
	values.Set("text", query)
	values.Set("sort", "0")
	values.Set("type", "1")
	if from != "" {
		values.Set("from", from)
	}
	u.RawQuery = values.Encode()

	ua, err := useragent.Desktop()
	if err != nil {
		return nil, fmt.Errorf("useragent.Desktop: %w", err)
	}

	resp, err := n.client.Get(ctx, u.String(), map[string]string{
		"User-Agent": ua,
	})
	if err != nil {
		return nil, fmt.Errorf("n.client.Get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}

	urls := doc.
		Find(".view-search__title a[href]").
		Map(func(i int, s *goquery.Selection) string {
			href, _ := s.Attr("href")
			u, _ := url.Parse(href)
			return u.EscapedPath()
		})

	return urls, nil
}
