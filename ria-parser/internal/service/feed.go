package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type newsFeed struct {
	*news
	browser *rod.Browser
}

func NewNewsFeed(baseAPI string, appID string, browser *rod.Browser) *newsFeed {
	client := httpclient.New()

	news := &news{
		client:  client,
		baseAPI: baseAPI,
		appID:   appID,
	}

	feed := &newsFeed{
		news:    news,
		browser: browser,
	}

	return feed
}

const PAGE_LAYOUT = "20060102"

func (n *newsFeed) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	var (
		pageObj time.Time
		err     error
	)

	if page == "" {
		pageObj = time.Now()
		page = pageObj.Format(PAGE_LAYOUT)
	} else {
		pageObj, err = time.Parse(PAGE_LAYOUT, page)
		if err != nil {
			return nil, "", fmt.Errorf("time.Parse: %w", err)
		}
	}

	urls, err := n.parseURLs(ctx, page)
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

	nextPage := pageObj.AddDate(0, 0, -1).Format(PAGE_LAYOUT)
	return news, nextPage, nil
}

func (n *newsFeed) parseURLs(ctx context.Context, path string) ([]string, error) {
	page, err := stealth.Page(n.browser)
	if err != nil {
		return nil, fmt.Errorf("stealth.Page: %w", err)
	}
	defer page.Close()

	page = page.Context(ctx)

	err = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: gofakeit.UserAgent()})
	if err != nil {
		return nil, fmt.Errorf("n.page.SetUserAgent: %w", err)
	}

	url, err := url.JoinPath(n.baseAPI, path)
	if err != nil {
		return nil, fmt.Errorf("utl.JoinPath: %w", err)
	}

	err = page.Navigate(url)
	if err != nil {
		return nil, fmt.Errorf("n.page.Navigate: %w", err)
	}

	err = n.loadPage(page)
	if err != nil {
		return nil, fmt.Errorf("n.loadPage: %w", err)
	}

	list, err := page.Element(".list")
	if err != nil {
		return nil, fmt.Errorf("page.Element: %w", err)
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

func (n *newsFeed) loadPage(page *rod.Page) error {
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
