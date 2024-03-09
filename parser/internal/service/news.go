package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/parser/internal/repo"
)

type Parser interface {
	Parse(query string, page string) ([]entity.News, string, error)
}

type news struct {
	NewsConfig
}

type NewsConfig struct {
	Repo   repo.News
	Parser Parser

	Producer   rabbitmq.Producer
	Exchange   string
	RoutingKey string
}

func NewNews(cfg NewsConfig) *news {
	return &news{
		NewsConfig: cfg,
	}
}

func (n *news) Parse(query string, page string) (string, error) {
	results, page, err := n.Parser.Parse(query, page)
	if err != nil {
		return "", fmt.Errorf("n.Parser.Parse: %w", err)
	}

	for _, result := range results {
		body, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("json.Marshal: %w", err)
		}

		err = n.Producer.Produce(n.Exchange, n.RoutingKey, rabbitmq.Message{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
		})
		if err != nil {
			// TODO: amqp.Produce error handling
			err := n.Repo.Create(context.Background(), string(body))
			if err != nil {
				return "", fmt.Errorf("n.Repo.News.Create: %w", err)
			}

			return "", fmt.Errorf("n.AMQP.Producer.Produce: %w", err)
		}
	}

	return page, nil
}
