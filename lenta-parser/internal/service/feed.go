package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/httpclient/httpresponse"
)

type newsFeed struct {
	*newsAbstract
}

func NewNewsFeed(baseAPI string) *newsFeed {
	client := httpclient.New()

	abstract := &newsAbstract{
		client:  client,
		baseAPI: baseAPI,
	}

	feed := &newsFeed{
		newsAbstract: abstract,
	}

	abstract.news = feed
	return feed
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

func (n *newsFeed) parseURLs(ctx context.Context, query string, page string) ([]*newsURL, error) {
	// rubrics
	u, _ := url.Parse(n.baseAPI + "/v3/rubrics")
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
	u, _ = url.Parse(n.baseAPI + "/v3/topics/by_rubrics")
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
