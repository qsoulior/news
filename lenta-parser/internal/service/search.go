package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DataHenHQ/useragent"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
)

type newsSearch struct {
	*news
	url string
}

func NewNewsSearch(appID string, url string, client *httpclient.Client) *newsSearch {
	news := &news{
		client: client,
		appID:  appID,
	}

	search := &newsSearch{
		news: news,
		url:  url,
	}

	return search
}

type MatchDTO struct {
	URL     string `json:"url"`
	PubDate int    `json:"pubdate"`
}

type MatchResponse struct {
	Matches []MatchDTO `json:"matches"`
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

func (n *newsSearch) parseURLs(ctx context.Context, query string, from string) ([]*newsURL, error) {
	u, _ := url.Parse(n.url + "/search/v2/process")
	values := u.Query()
	values.Set("query", query)
	values.Set("from", from)
	values.Set("size", "100")
	values.Set("sort", "2")
	values.Set("domain", "1")
	values.Set("type", "1")
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

	data, err := httpresponse.JSON[MatchResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("httpresponse.JSON[MatchResponse]: %w", err)
	}

	urls := make([]*newsURL, 0, len(data.Matches))
	for _, match := range data.Matches {
		urls = append(urls, &newsURL{
			URL:         match.URL,
			PublishedAt: time.Unix(int64(match.PubDate), 0),
		})
	}

	return urls, nil
}
