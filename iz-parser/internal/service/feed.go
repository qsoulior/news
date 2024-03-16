package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
)

type newsFeed struct {
	*newsAbstract
}

func NewNewsFeed(baseAPI string) *newsFeed {
	client := httpclient.New(
		httpclient.URL(baseAPI),
		httpclient.Headers(map[string]string{
			"Referer": "https://iz.ru/feed",
		}),
	)

	abstract := &newsAbstract{
		client: client,
	}

	feed := &newsFeed{
		newsAbstract: abstract,
	}

	abstract.news = feed
	return feed
}

type ViewDTO struct {
	Command string `json:"command"`
	Method  string `json:"method"`
	Data    string `json:"data"`
}

func (n *newsFeed) parseURLs(query string, page string) ([]string, error) {
	u, _ := url.Parse("/views/ajax?_wrapper_format=drupal_ajax")

	var reqData url.Values
	reqData.Set("view_name", "content_field")
	reqData.Set("view_display_id", "page_feed")
	reqData.Set("page", page)

	resp, err := n.client.Post(u.String(), strings.NewReader(reqData.Encode()), map[string]string{
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
		log.Println(err)
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
