package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
	AppID      string
}

func NewNews(cfg NewsConfig) *news {
	return &news{
		NewsConfig: cfg,
	}
}

func (n *news) Parse(ctx context.Context, query string, page string) (int, string, error) {
	count := 0

	results, page, err := n.Parser.Parse(ctx, query, page)
	if err != nil {
		return count, "", fmt.Errorf("n.Parser.Parse: %w", err)
	}

	for _, result := range results {
		body, err := json.Marshal(result)
		if err != nil {
			return count, "", fmt.Errorf("json.Marshal: %w", err)
		}

		if err := n.produce(ctx, body); err == nil {
			count++
			continue
		}

		if err := n.Repo.Create(ctx, string(body)); err != nil {
			return count, "", fmt.Errorf("n.Repo.Create: %w", err)
		}
		count++
	}

	return count, page, nil
}

func (n *news) Release(ctx context.Context) (int, error) {
	count := 0
	for {
		jsonStr, err := n.Repo.GetFirst(ctx)
		if err != nil && !errors.Is(err, ErrNotExist) {
			return count, fmt.Errorf("n.Repo.GetFirst: %w", err)
		}

		if jsonStr == "" {
			return count, nil
		}

		err = n.produce(ctx, []byte(jsonStr))
		if err != nil {
			return count, err
		}

		count++

		err = n.Repo.DeleteFirst(ctx)
		if err != nil {
			return count, fmt.Errorf("n.Repo.DeleteFirst: %w", err)
		}
	}
}

func (n *news) produce(ctx context.Context, body []byte) error {
	err := n.Producer.Produce(ctx, n.Exchange, n.RoutingKey, rabbitmq.Message{
		AppId:        n.AppID,
		MessageId:    uuid.NewString(),
		ContentType:  "application/json",
		DeliveryMode: 2,
		Body:         body,
	})

	if err != nil {
		return fmt.Errorf("n.Producer.Produce: %w", err)
	}

	return nil
}
