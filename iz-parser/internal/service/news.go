package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DataHenHQ/useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/rs/zerolog"
)

type news struct {
	appID  string
	client *httpclient.Client
	logger *zerolog.Logger
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
		NewsHead: entity.NewsHead{
			Source: n.appID,
		},
		Link: resp.Request.URL.String(),
	}

	article := doc.Find("[role=\"article\"]")

	news.Title = article.Find("[itemprop=\"headline\"] span").Text()
	news.Description = article.Find("[itemprop=\"alternativeHeadline\"]").Text()
	news.Authors = article.
		Find(".article_page__left__top__author [itemprop=\"name\"]").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	datetimeStr, ok := article.Find(".article_page__left__top__time time").Attr("datetime")
	if !ok {
		return nil, errors.New("empty datetime")
	}
	datetimeLayout := "2006-01-02T15:04:05Z"
	news.PublishedAt, err = time.Parse(datetimeLayout, datetimeStr)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}

	news.Tags = article.Find(".article_page__left__top__left__hash_tags a").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	news.Categories = article.Find(".rubrics_btn a").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})

	var text strings.Builder
	textItems := article.Find("[itemprop=\"articleBody\"] p")
	textItems.Each(func(i int, s *goquery.Selection) {
		text.WriteString(s.Text())
		if i < textItems.Size()-1 {
			text.WriteRune('\n')
		}
	})

	news.Content = text.String()
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
			n.logger.Warn().Err(err).Str("url", url).Send()
			continue
		}
		news = append(news, *newsItem)
	}

	return news, nil
}
