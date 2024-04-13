package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DataHenHQ/useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsURL struct {
	URL         string
	PublishedAt time.Time
}

type news struct {
	appID  string
	client *httpclient.Client
}

func (n *news) parseOne(ctx context.Context, url *newsURL) (*entity.News, error) {
	ua, err := useragent.Desktop()
	if err != nil {
		return nil, fmt.Errorf("useragent.Desktop: %w", err)
	}

	resp, err := n.client.Get(ctx, url.URL, map[string]string{
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

	news := &entity.News{
		NewsHead: entity.NewsHead{
			Source:      n.appID,
			PublishedAt: url.PublishedAt,
		},
		Link: resp.Request.URL.String(),
		Tags: make([]string, 0),
	}

	topic := doc.Find(".topic-page__container")

	news.Title = topic.Find(".topic-body__title").Text()
	news.Description = topic.Find(".topic-body__title-yandex").Text()
	news.Categories = topic.
		Find(".topic-header .topic-header__rubric").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	var text strings.Builder
	textItems := topic.Find(".topic-body__content .topic-body__content-text")
	textItems.Each(func(i int, s *goquery.Selection) {
		text.WriteString(s.Text())
		if i < textItems.Size()-1 {
			text.WriteRune('\n')
		}
	})

	news.Content = text.String()

	news.Authors = topic.
		Find(".topic-authors .topic-authors__author").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	return news, nil
}

func (n *news) parseMany(ctx context.Context, urls []*newsURL) ([]entity.News, error) {
	news := make([]entity.News, 0, len(urls))
	for _, url := range urls {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		newsItem, err := n.parseOne(ctx, url)
		if err != nil {
			log.Println(err)
			continue
		}
		news = append(news, *newsItem)
	}

	return news, nil
}
