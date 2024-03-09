package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
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
}

func NewNews(cfg NewsConfig) *news {
	client := httpclient.New(
		httpclient.URL(cfg.BaseAPI),
	)

	return &news{
		NewsConfig: cfg,
		client:     client,
	}
}

func (n *news) Parse(query string, page string) ([]entity.News, string, error) {
	u, _ := url.Parse("/news")
	values := u.Query()
	values.Set("apikey", n.AccessKey)
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
		return nil, "", fmt.Errorf("n.client.Get: %w", err)
	}
	fmt.Println(resp.Request.URL)

	switch resp.StatusCode {
	case http.StatusOK:
		return n.parseResult(resp)
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusConflict, http.StatusUnsupportedMediaType, http.StatusUnprocessableEntity, http.StatusTooManyRequests:
		return nil, "", n.parseError(resp)
	case http.StatusInternalServerError:
		return nil, "", ErrInternalServer
	}

	return nil, "", ErrUnexpectedCode
}

func (n *news) parseResult(resp *http.Response) ([]entity.News, string, error) {
	data, err := httpresponse.JSON[NewsResponseSuccess](resp)
	if err != nil {
		return nil, "", fmt.Errorf("httpresponse.JSON: %w", err)
	}

	news := make([]entity.News, len(data.Results))
	for i, result := range data.Results {
		news[i] = *result.Entity()
	}

	return news, data.NextPage, nil
}

func (n *news) parseError(resp *http.Response) error {
	data, err := httpresponse.JSON[NewsResponseError](resp)
	if err != nil {
		return fmt.Errorf("httpresponse.JSON: %w", err)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return &ResponseError{ErrRateLimit, data.Results.Code}
	}
	return &ResponseError{ErrRequestInvalid, data.Results.Code}
}
