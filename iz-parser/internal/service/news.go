package service

import (
	"context"
	"errors"
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

type news interface {
	parseURLs(ctx context.Context, query string, page string) ([]string, error)
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

func (n *newsAbstract) parseOne(ctx context.Context, url string) (*entity.News, error) {
	resp, err := n.client.Get(ctx, url, map[string]string{
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
		Source: "iz",
		Link:   resp.Request.URL.String(),
	}

	article := doc.Find("[role=\"article\"]")

	news.Title = article.Find("[itemprop=\"headline\"] span").Text()

	// TODO: add description to entity
	_ = article.Find("[itemprop=\"alternativeHeadline\"]").Text()

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
