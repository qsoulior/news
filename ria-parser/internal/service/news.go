package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

const BASE_COUNT = 20

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
	return nil, nil
}

func (n *news) parseView(page *rod.Page, limit int) ([]string, error) {
	err := n.loadView(page, limit)
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

func (n *news) loadView(page *rod.Page, limit int) error {
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

	count := BASE_COUNT * 2
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

		if !visible || (limit > 0 && count >= limit) {
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
