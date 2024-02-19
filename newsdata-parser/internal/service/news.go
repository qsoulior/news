package service

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/qsoulior/news/newsdata-parser/pkg/httpclient"
	"github.com/qsoulior/news/newsdata-parser/pkg/httpclient/httpresponse"
)

const (
	COUNTRY = "ru"
)

type NewsDTO struct {
	ArticleID  string    `json:"article_id"`
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	SourceID   string    `json:"source_id"`
	PubDate    time.Time `json:"pub_date"`
	Creator    []string  `json:"creator"`
	Keywords   []string  `json:"keywords"`
	Categories []string  `json:"categories"`
	Content    string    `json:"content"`
}

type news struct {
	client *httpclient.Client
}

type NewsConfig struct {
	BaseAPI   string
	AccessKey string
}

func NewNews(cfg NewsConfig) News {
	client := httpclient.NewClient(
		httpclient.Headers(map[string]string{
			"X-ACCESS-KEY": cfg.AccessKey,
		}),
		httpclient.URL(cfg.BaseAPI),
	)

	return &news{client: client}
}

func (n *news) Parse(query string) error {
	u, err := url.Parse("/news")
	if err != nil {
		log.Fatal(err)
	}

	values := u.Query()
	values.Set("country", COUNTRY)
	if query != "" {
		values.Set("q", query)
	}

	u.RawQuery = values.Encode()
	resp, err := n.client.Get(u.String(), nil)
	if err != nil {
		return fmt.Errorf("n.client.Get: %w", err)
	}

	news, err := httpresponse.JSON[NewsDTO](resp)
	if err != nil {
		return fmt.Errorf("httpresponse.JSON: %w", err)
	}

	fmt.Println(news)
	return nil
}
