package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/qsoulior/news/newsdata-parser/internal/entity"
	"github.com/qsoulior/news/newsdata-parser/internal/repo"
	"github.com/qsoulior/news/newsdata-parser/pkg/httpclient"
	"github.com/qsoulior/news/newsdata-parser/pkg/httpclient/httpresponse"
	"github.com/qsoulior/news/newsdata-parser/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
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

func (dto *NewsDTO) Entity() *entity.News {
	entity := &entity.News{
		Title:       dto.Title,
		Link:        dto.Link,
		Source:      dto.SourceID,
		PublishedAt: dto.PubDate,
		Authors:     make([]string, len(dto.Creator)),
		Tags:        make([]string, len(dto.Keywords)),
		Categories:  make([]string, len(dto.Categories)),
		Content:     dto.Content,
	}

	copy(entity.Authors, dto.Creator)
	copy(entity.Tags, dto.Keywords)
	copy(entity.Categories, dto.Categories)

	return entity
}

type NewsResponseSuccess struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Results      []NewsDTO `json:"results"`
	NextPage     string    `json:"nextPage"`
}

type NewsResponseError struct {
	Status  string `json:"status"`
	Results struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"results"`
}

type news struct {
	NewsConfig
	client *httpclient.Client
}

type NewsConfig struct {
	BaseAPI   string
	AccessKey string
	AMQP      struct {
		Producer   *rabbitmq.Producer
		Exchange   string
		RoutingKey string
	}
	repo repo.News
}

func NewNews(cfg NewsConfig) News {
	client := httpclient.NewClient(
		httpclient.Headers(map[string]string{
			"X-ACCESS-KEY": cfg.AccessKey,
		}),
		httpclient.URL(cfg.BaseAPI),
	)

	return &news{
		NewsConfig: cfg,
		client:     client,
	}
}

func (n *news) Parse(query string, page string) (string, error) {
	u, _ := url.Parse("/news")
	values := u.Query()
	values.Set("country", COUNTRY)
	if query != "" {
		values.Set("q", query)
	}
	if page != "" {
		values.Set("page", page)
	}

	u.RawQuery = values.Encode()
	resp, err := n.client.Get(u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("n.client.Get: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return n.parseResult(resp)
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusConflict, http.StatusUnsupportedMediaType, http.StatusUnprocessableEntity, http.StatusTooManyRequests:
		return n.parseError(resp)
	case http.StatusInternalServerError:
		return "", ErrInternalServer
	}

	return "", ErrUnexpectedCode
}

func (n *news) parseResult(resp *http.Response) (string, error) {
	data, err := httpresponse.JSON[NewsResponseSuccess](resp)
	if err != nil {
		return "", fmt.Errorf("httpresponse.JSON: %w", err)
	}

	for _, result := range data.Results {
		body, err := json.Marshal(result.Entity())
		if err != nil {
			return "", fmt.Errorf("json.Marshal: %w", err)
		}

		err = n.AMQP.Producer.Produce(n.AMQP.Exchange, n.AMQP.RoutingKey, amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
		})
		if err != nil {
			// TODO: amqp.Produce error handling
			err := n.repo.Create(context.Background(), string(body))
			if err != nil {
				return "", fmt.Errorf("n.Repo.News.Create: %w", err)
			}

			return "", fmt.Errorf("n.amqp.Produce: %w", err)
		}
	}
	return data.NextPage, nil
}

func (n *news) parseError(resp *http.Response) (string, error) {
	data, err := httpresponse.JSON[NewsResponseError](resp)
	if err != nil {
		return "", fmt.Errorf("httpresponse.JSON: %w", err)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return "", &ResponseError{ErrRateLimit, data.Results.Code}
	}
	return "", &ResponseError{ErrRequestInvalid, data.Results.Code}
}
