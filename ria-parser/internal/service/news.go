package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"golang.org/x/sync/semaphore"
)

type newsURL struct {
	URL         string
	PublishedAt time.Time
}

type news struct {
	baseAPI string
	appID   string
	client  *httpclient.Client
}

func (n *news) parseOne(ctx context.Context, url string) (*entity.News, error) {
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
	newsCh := make(chan *entity.News, len(urls))
	news := make([]entity.News, 0, len(urls))

	var wg sync.WaitGroup
	maxProcs := int64(runtime.GOMAXPROCS(0)) * 2
	sem := semaphore.NewWeighted(maxProcs)

	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, url := range urls {
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, fmt.Errorf("sem.Acquire: %w", err)
		}

		go func(ctx context.Context) {
			defer sem.Release(1)
			wg.Add(1)
			defer wg.Done()
			newsItem, err := n.parseOne(ctx, url)
			if err != nil {
				// TODO: Logger for parseOne
				log.Println(err)
			} else {
				newsCh <- newsItem
			}
		}(cancelCtx)

	}

	go func() {
		wg.Wait()
		close(newsCh)
	}()

	for newsItem := range newsCh {
		news = append(news, *newsItem)
	}

	return news, nil
}

func (n *news) parseView(page *rod.Page) ([]string, error) {
	err := n.loadView(page)
	if err != nil {
		return nil, fmt.Errorf("n.loadView: %w", err)
	}

	list, err := page.Element(".list")
	if err != nil {
		return nil, fmt.Errorf("page.Element: %w", err)
	}

	err = list.WaitStable(1 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("list.WaitStable: %w", err)
	}

	listHTML, err := list.HTML()
	if err != nil {
		return nil, fmt.Errorf("list.HTML: %w", err)
	}

	listDocument, err := goquery.NewDocumentFromReader(strings.NewReader(listHTML))
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}

	urls := listDocument.Find(".list-item__title[href]").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Attr("href")
		return href
	})

	return urls, nil
}

func (n *news) loadView(page *rod.Page) error {
	listMore, err := page.Element(".list-more")
	if err != nil {
		return fmt.Errorf("n.page.Element: %w", err)
	}

	err = listMore.Timeout(5 * time.Second).WaitStableRAF()
	if err != nil {
		return fmt.Errorf("listMore.WaitStable: %w", err)
	}

	err = listMore.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("listMore.Click: %w", err)
	}

	for {
		err = listMore.Timeout(5 * time.Second).WaitStableRAF()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			return fmt.Errorf("listMore.WaitStableRAF: %w", err)
		}

		visible, err := listMore.Visible()
		if err != nil {
			return fmt.Errorf("listMore.Visible: %w", err)
		}

		if !visible {
			break
		}

		loading, err := listMore.Matches(".loading")
		if err != nil {
			return fmt.Errorf("listMore.Matches: %w", err)
		}

		if loading {
			continue
		}

		err = listMore.Timeout(5 * time.Second).ScrollIntoView()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			return fmt.Errorf("listMore.ScrollIntoView: %w", err)
		}
	}

	return nil
}
