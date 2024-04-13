package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
)

const (
	COUNTRY = "ru"
)

type NewsDTO struct {
	ArticleID   string      `json:"article_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Link        string      `json:"link"`
	SourceID    string      `json:"source_id"`
	PubDate     NewsPubDate `json:"pubDate"`
	Creator     []string    `json:"creator"`
	Keywords    []string    `json:"keywords"`
	Categories  []string    `json:"category"`
	Content     string      `json:"content"`
}

type NewsPubDate time.Time

func (p *NewsPubDate) UnmarshalJSON(b []byte) error {
	var t time.Time

	s := strings.Trim(string(b), "\"")
	if s != "null" {
		var err error
		t, err = time.Parse(time.DateTime, s)
		if err != nil {
			return err
		}
	}

	*p = NewsPubDate(t)
	return nil
}

func (dto *NewsDTO) Entity() *entity.News {
	entity := &entity.News{
		NewsHead: entity.NewsHead{
			Title:       dto.Title,
			Description: dto.Description,
			Source:      dto.SourceID,
			PublishedAt: time.Time(dto.PubDate),
		},
		Link:       dto.Link,
		Authors:    make([]string, len(dto.Creator)),
		Tags:       make([]string, len(dto.Keywords)),
		Categories: make([]string, len(dto.Categories)),
		Content:    dto.Content,
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
	appID     string
	accessKey string
	client    *httpclient.Client
}

func NewNews(appID string, accessKey string, client *httpclient.Client) *news {
	return &news{
		appID:     appID,
		accessKey: accessKey,
		client:    client,
	}
}

func (n *news) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	u, _ := url.Parse("/news")
	values := u.Query()
	values.Set("apikey", n.accessKey)
	values.Set("country", COUNTRY)
	if query != "" {
		values.Set("q", query)
	}
	if page != "" {
		values.Set("page", page)
	}

	u.RawQuery = values.Encode()
	resp, err := n.client.Get(ctx, u.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("n.client.Get: %w", err)
	}

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
		entity := result.Entity()
		entity.Source = n.appID
		news[i] = *entity
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
