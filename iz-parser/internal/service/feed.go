package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
	"github.com/rs/zerolog"
)

type NewsConfig struct {
	BaseAPI string
	Logger  *zerolog.Logger
}

type newsFeed struct {
	NewsConfig
	client *httpclient.Client
}

func NewNewsFeed(cfg NewsConfig) *newsFeed {
	client := httpclient.New(
		// httpclient.CookieJar(jar),
		httpclient.URL(cfg.BaseAPI),
		httpclient.Headers(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
			"Referer":    "https://iz.ru/feed",
		}),
	)

	return &newsFeed{
		NewsConfig: cfg,
		client:     client,
	}
}

type ViewDTO struct {
	Command string `json:"command"`
	Method  string `json:"method"`
	Data    string `json:"data"`
}

func (n *newsFeed) Parse(query string, page string) ([]entity.News, string, error) {
	u, _ := url.Parse("/views/ajax?_wrapper_format=drupal_ajax")

	var reqData url.Values
	reqData.Set("view_name", "content_field")
	reqData.Set("view_display_id", "page_feed")
	reqData.Set("page", page)

	resp, err := n.client.Post(u.String(), strings.NewReader(reqData.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	})

	if err != nil {
		return nil, "", fmt.Errorf("n.client.Post: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", NewStatusError(resp.StatusCode)
	}

	return n.parseResult(resp, page)
}

func (n *newsFeed) parseResult(resp *http.Response, page string) ([]entity.News, string, error) {
	respData, err := httpresponse.JSON[[]ViewDTO](resp)
	if err != nil {
		log.Println(err)
	}

	index := slices.IndexFunc(*respData, func(item ViewDTO) bool {
		return item.Command == "insert" && (item.Method == "infiniteScrollInsertView" || item.Method == "replaceWith")
	})

	if index == -1 {
		return nil, "", errors.New("response does not contain valid data")
	}

	data := (*respData)[index].Data
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}

	urls := doc.
		Find(".lenta_news__day .lenta_news__day__list__item[href]").
		Map(func(i int, s *goquery.Selection) string {
			href, _ := s.Attr("href")
			return href
		})

	news := make([]entity.News, 0, len(urls))
	for _, url := range urls {
		newsItem, err := n.parseOne(url)
		if err != nil {
			n.Logger.Error().Err(err).Str("url", url).Msg("")
			continue
		}
		news = append(news, *newsItem)
	}

	nextPage, err := strconv.Atoi(page)
	if err != nil {
		return news, "0", nil
	}

	return news, strconv.Itoa(nextPage + 1), nil
}

func (n *newsFeed) parseOne(url string) (*entity.News, error) {
	resp, err := n.client.Get(url, nil)
	if err != nil {
		return nil, fmt.Errorf("n.client.Get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, NewStatusError(resp.StatusCode)
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
		Find(".article_page__left__top__author div[itemprop=\"author\"] span[itemprop=\"name\"]").
		Map(func(i int, s *goquery.Selection) string { return s.Text() })

	datetimeStr, ok := article.Find(".article_page__left__top__time time").Attr("datetime")
	if ok {
		datetimeLayout := "2006-01-02T15:04:05Z"
		news.PublishedAt, err = time.Parse(datetimeLayout, datetimeStr)
		if err != nil {
			n.Logger.Warn().
				Err(err).
				Str("url", url).
				Str("layout", datetimeLayout).
				Str("value", datetimeStr).
				Msg("incorrect datetime layout")
		}
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
