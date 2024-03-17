package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsURL struct {
	url         string
	publishedAt time.Time
}

type news interface {
	parseURLs(ctx context.Context, query string, page string) ([]newsURL, error)
}

type newsAbstract struct {
	news
	client *httpclient.Client
}

func (n *newsAbstract) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
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

func (n *newsAbstract) parseOne(ctx context.Context, url newsURL) (*entity.News, error) {
	resp, err := n.client.Get(ctx, url.url, map[string]string{
		"User-Agent": gofakeit.UserAgent(),
	})
	if err != nil {
		return nil, fmt.Errorf("n.client.Get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()
	_, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}

	news := &entity.News{
		Source:      "lenta.ru",
		Link:        resp.Request.URL.String(),
		PublishedAt: url.publishedAt,
	}

	// TODO
	return news, nil
}
