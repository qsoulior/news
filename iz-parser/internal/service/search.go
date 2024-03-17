package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsSearch struct {
	*newsAbstract
}

func NewNewsSearch(baseAPI string, appID string) *newsSearch {
	client := httpclient.New(
		httpclient.URL(baseAPI),
		httpclient.Headers(map[string]string{
			"Referer": "https://iz.ru/feed",
		}),
	)

	abstract := &newsAbstract{
		client: client,
		appID:  appID,
	}

	search := &newsSearch{
		newsAbstract: abstract,
	}

	abstract.news = search
	return search
}

func (n *newsSearch) parseURLs(ctx context.Context, query string, page string) ([]string, error) {
	u, _ := url.Parse("/search")
	values := u.Query()
	values.Set("text", query)
	values.Set("sort", "0")
	values.Set("type", "1")
	values.Set("from", page+"0")
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

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}

	urls := doc.
		Find(".view-search__title a[href]").
		Map(func(i int, s *goquery.Selection) string {
			href, _ := s.Attr("href")
			return href
		})

	return urls, nil
}
