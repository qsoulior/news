package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type newsFeed struct {
	*newsAbstract
	page *rod.Page
}

func NewNewsFeed(baseAPI string, appID string) *newsFeed {
	abstract := &newsAbstract{
		baseAPI: baseAPI,
		appID:   appID,
	}

	feed := &newsFeed{
		newsAbstract: abstract,
	}

	abstract.news = feed
	return feed
}

func (n *newsFeed) parseURLs(ctx context.Context, query string, page string) ([]string, error) {
	err := n.page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: gofakeit.UserAgent()})
	if err != nil {
		return nil, fmt.Errorf("n.page.SetUserAgent: %w", err)
	}

	err = n.page.Navigate("https://ria.ru/" + page)
	if err != nil {
		return nil, fmt.Errorf("n.page.Navigate: %w", err)
	}

	err = n.loadPage()
	if err != nil {
		return nil, fmt.Errorf("n.loadPage: %w", err)
	}

	list, err := n.page.Element(".list")
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

func (n *newsFeed) loadPage() error {
	listMore, err := n.page.Element(".list-more")
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
