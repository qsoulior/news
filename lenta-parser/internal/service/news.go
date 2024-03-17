package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsURL struct {
	URL         string
	PublishedAt time.Time
}

type news interface {
	parseURLs(ctx context.Context, query string, page string) ([]*newsURL, error)
}

type newsAbstract struct {
	news
	baseAPI string
	client  *httpclient.Client
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

func (n *newsAbstract) parseOne(ctx context.Context, url *newsURL) (*entity.News, error) {
	resp, err := n.client.Get(ctx, url.URL, map[string]string{
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

	news := &entity.News{
		Source:      "lenta.ru",
		Link:        resp.Request.URL.String(),
		PublishedAt: url.PublishedAt,
	}

	topic := doc.Find(".topic-page__container")

	news.Title = topic.Find(".topic-body__title").Text()
	news.Description = topic.Find(".topic-body__title-yandex").Text()
	news.Categories = topic.
		Find(".topic-header .topic-header__rubric").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	var text strings.Builder
	textItems := topic.Find("topic-body__content topic-body__content-text")
	textItems.Each(func(i int, s *goquery.Selection) {
		text.WriteString(s.Text())
		if i < textItems.Size()-1 {
			text.WriteRune('\n')
		}
	})

	news.Content = text.String()

	news.Authors = topic.
		Find(".topic-authors topic-authors__author").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	return news, nil
}
