package service

import (
	"context"
	"encoding/json"
	"errors"
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

		err = n.store(ctx, body)
		if err != nil {
			return "", err
		}
	}

	return page, nil
}

func (n *news) Release(ctx context.Context) error {
	for {
		jsonStr, err := n.Repo.Pop(ctx)
		if err != nil && !errors.Is(err, ErrNotExist) {
			return fmt.Errorf("n.Repo.Pop: %w", err)
		}

		if jsonStr == "" {
			return nil
		}

		err = n.store(ctx, []byte(jsonStr))
		if err != nil {
			return err
		}
	}
}

func (n *news) store(ctx context.Context, body []byte) error {
	err := n.Producer.Produce(ctx, n.Exchange, n.RoutingKey, rabbitmq.Message{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Body:         body,
	})

	if err != nil {
		err := n.Repo.Create(ctx, string(body))
		if err != nil {
			return fmt.Errorf("n.Repo.Create: %w", err)
		}
	}

	return nil
}
