package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
)

type newsArchive struct {
	*news
	URL string
}

func NewNewsArchive(appID string, url string, client *httpclient.Client) *newsArchive {
	news := &news{
		appID:  appID,
		client: client,
	}

	archive := &newsArchive{
		news: news,
		URL:  url,
	}

	return archive
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

type RubricDTO struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type RubricResponse struct {
	Rubrics []RubricDTO `json:"rubrics"`
}

type Topic struct {
	Headline struct {
		Info struct {
			Title    string `json:"title"`
			Modified int    `json:"modified"`
		} `json:"info"`
		Links struct {
			Public string `json:"public"`
			Self   string `json:"self"`
		} `json:"links"`
		Type string `json:"type"`
	} `json:"headline"`
}

type TopicResponse struct {
	Topics []Topic `json:"topics"`
}

func (n *newsArchive) parseURLs(ctx context.Context, page string) ([]*newsURL, error) {
	// rubrics
	u, _ := url.Parse(n.URL + "/v3/rubrics")
	rubricResp, err := n.client.Get(ctx, u.String(), map[string]string{
		"User-Agent": gofakeit.UserAgent(),
	})
	if err != nil {
		return nil, fmt.Errorf("n.client.Get: %w", err)
	}
	if rubricResp.StatusCode != http.StatusOK {
		return nil, newStatusError(rubricResp.StatusCode)
	}

	rubricData, err := httpresponse.JSON[RubricResponse](rubricResp)
	if err != nil {
		return nil, fmt.Errorf("httpresponse.JSON[RubricResponse]: %w", err)
	}

	// topics
	u, _ = url.Parse(n.URL + "/v3/topics/by_rubrics")
	values := u.Query()

	for _, rubric := range rubricData.Rubrics {
		values.Add("rubric[]", rubric.Slug)
	}

	values.Set("limit", "100")
	values.Set("offset", page+"00")
	u.RawQuery = values.Encode()

	topicResp, err := n.client.Get(ctx, u.String(), map[string]string{
		"User-Agent": gofakeit.UserAgent(),
	})
	if err != nil {
		return nil, fmt.Errorf("n.client.Get: %w", err)
	}
	if topicResp.StatusCode != http.StatusOK {
		return nil, newStatusError(topicResp.StatusCode)
	}

	topicData, err := httpresponse.JSON[TopicResponse](topicResp)
	if err != nil {
		return nil, fmt.Errorf("httpresponse.JSON[TopicResponse]: %w", err)
	}

	urls := make([]*newsURL, 0, len(topicData.Topics))
	for _, topic := range topicData.Topics {
		if topic.Headline.Type == "news" {
			urls = append(urls, &newsURL{
				URL:         topic.Headline.Links.Public,
				PublishedAt: time.Unix(int64(topic.Headline.Info.Modified), 0),
			})
		}
	}

	return urls, nil
}
