package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
)

type newsArchive struct {
	*news
}

func NewNewsArchive(appID string, client *httpclient.Client) *newsArchive {
	news := &news{
		appID:  appID,
		client: client,
	}

	archive := &newsArchive{
		news: news,
	}

	return archive
}

type ViewDTO struct {
	Command string `json:"command"`
	Method  string `json:"method"`
	Data    string `json:"data"`
}

func (n *newsArchive) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	urls, err := n.parseURLs(ctx, page)
	if err != nil {
		return nil, "", err
	}

	news, err := n.parseMany(ctx, urls)
	if err != nil {
		return nil, "", fmt.Errorf("n.parseMany: %w", err)
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

func (n *newsArchive) parseURLs(ctx context.Context, page string) ([]string, error) {
	u, _ := url.Parse("/views/ajax?_wrapper_format=drupal_ajax")

	reqData := make(url.Values, 3)
	reqData.Set("view_name", "content_field")
	reqData.Set("view_display_id", "page_feed")
	reqData.Set("page", page)

	resp, err := n.client.Post(ctx, u.String(), strings.NewReader(reqData.Encode()), map[string]string{
		"User-Agent":   gofakeit.UserAgent(),
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	})

	if err != nil {
		return nil, fmt.Errorf("n.client.Post: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp.StatusCode)
	}

	respData, err := httpresponse.JSON[[]ViewDTO](resp)
	if err != nil {
		return nil, fmt.Errorf("httpresponse.JSON[[]ViewDTO]: %w", err)
	}

	index := slices.IndexFunc(*respData, func(item ViewDTO) bool {
		return item.Command == "insert" && (item.Method == "infiniteScrollInsertView" || item.Method == "replaceWith")
	})

	if index == -1 {
		return nil, errors.New("response does not contain valid data")
	}

	data := (*respData)[index].Data
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}

	urls := doc.
		Find(".lenta_news__day__list__item[href]").
		Map(func(i int, s *goquery.Selection) string {
			href, _ := s.Attr("href")
			return href
		})

	return urls, nil
}
