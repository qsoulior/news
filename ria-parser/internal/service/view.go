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
	"github.com/go-rod/stealth"
)

type newsView struct {
	URL     string
	browser *rod.Browser
}

func (n *newsView) parseURLs(ctx context.Context, path string) ([]string, error) {
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

	err = page.Navigate(n.URL + path)
	if err != nil {
		return nil, fmt.Errorf("n.page.Navigate: %w", err)
	}

	urls, err := n.parseView(page)
	if err != nil {
		return nil, fmt.Errorf("n.parseView: %w", err)
	}

	return urls, nil
}

func (n *newsView) parseView(page *rod.Page) ([]string, error) {
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

func (n *newsView) loadView(page *rod.Page) error {
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
