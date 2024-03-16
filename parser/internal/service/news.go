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
	Parse(ctx context.Context, query string, page string) ([]entity.News, string, error)
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

func (n *news) Parse(ctx context.Context, query string, page string) (string, error) {
	results, page, err := n.Parser.Parse(ctx, query, page)
	if err != nil {
		return "", fmt.Errorf("n.Parser.Parse: %w", err)
	}

	for _, result := range results {
		body, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("json.Marshal: %w", err)
		}

		err = n.Producer.Produce(ctx, n.Exchange, n.RoutingKey, rabbitmq.Message{
			ContentType:  "application/json",
			DeliveryMode: 2,
			Body:         body,
		})
		if err != nil {
			// TODO: amqp.Produce error handling
			if err := n.Repo.Create(ctx, string(body)); err != nil {
				return "", fmt.Errorf("n.Repo.News.Create: %w", err)
			}

			return "", fmt.Errorf("n.AMQP.Producer.Produce: %w", err)
		}
	}

	return page, nil
}
