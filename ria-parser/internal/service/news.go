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

type news struct {
	appID  string
	client *httpclient.Client
}

func (n *news) parseOne(ctx context.Context, url string) (*entity.News, error) {
	ua, err := useragent.Desktop()
	if err != nil {
		return nil, fmt.Errorf("useragent.Desktop: %w", err)
	}
	resp, err := n.client.Get(ctx, url, map[string]string{
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
		Source: n.appID,
		Link:   resp.Request.URL.String(),
	}

	article := doc.Find(".article")

	news.Categories = article.
		Find(".article__supertag-header .article__supertag-header-title").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	datetimeLayout := "15:04 02.01.2006"
	datetimeStr := strings.TrimSpace(
		article.Find(".article__info-date .article__info-date-modified").Children().Remove().End().Text(),
	)

	if datetimeStr == "" {
		datetimeStr = article.Find(".article__info-date a").Text()
	} else {
		parts := strings.SplitN(datetimeStr, " ", 2)
		if len(parts) > 1 {
			datetimeStr = strings.Trim(parts[1], "()")
		}
	}

	news.PublishedAt, err = time.Parse(datetimeLayout, datetimeStr)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}

	news.Title = article.Find(".article__title").Text()
	news.Description = article.Find(".article__second-title").Text()

	news.Authors = article.
		Find(".article__author .article__author-name").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	var text strings.Builder
	textItems := article.Find(".article__body .article__text, .article__body .article__quote-text")
	textItems.Each(func(i int, s *goquery.Selection) {
		text.WriteString(s.Text())
		if i < textItems.Size()-1 {
			text.WriteRune('\n')
		}
	})
	news.Content = text.String()

	news.Tags = article.
		Find(".article__tags a").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	return news, nil
}

func (n *news) parseMany(ctx context.Context, urls []string) ([]entity.News, error) {
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
