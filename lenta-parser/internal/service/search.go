package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
)

type newsSearch struct {
	*newsAbstract
}

func NewNewsSearch(baseAPI string) *newsSearch {
	client := httpclient.New()

	abstract := &newsAbstract{
		client:  client,
		baseAPI: baseAPI,
	}

	search := &newsSearch{
		newsAbstract: abstract,
	}

	abstract.news = search
	return search
}

type MatchDTO struct {
	URL     string `json:"url"`
	PubDate int    `json:"pubdate"`
}

type MatchResponse struct {
	Matches []MatchDTO `json:"matches"`
}

func (n *newsSearch) parseURLs(ctx context.Context, query string, page string) ([]*newsURL, error) {
	u, _ := url.Parse(n.baseAPI + "/search/v2/process")
	values := u.Query()
	values.Set("query", query)
	values.Set("from", page+"00")
	values.Set("size", "100")
	values.Set("sort", "2")
	values.Set("domain", "1")
	values.Set("type", "1")
	u.RawQuery = values.Encode()

	resp, err := n.client.Get(ctx, u.String(), map[string]string{
		"User-Agent": gofakeit.UserAgent(),
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
