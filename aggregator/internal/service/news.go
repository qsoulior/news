package service

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/internal/repo"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
)

type (
	NewsConfig struct {
		Repo       repo.News
		Producer   rabbitmq.Producer
		Exchange   string
		RoutingKey string
	}
)

type news struct {
	NewsConfig
}

func NewNews(cfg NewsConfig) News {
	return &news{cfg}
}

func (n *news) Create(news entity.News) error {
	if err := n.Repo.Create(context.Background(), news); err != nil {
		return fmt.Errorf("n.repo.Create: %w", err)
	}

	return nil
}

func (n *news) CreateMany(news []entity.News) error {
	if err := n.Repo.CreateMany(context.Background(), news); err != nil {
		return fmt.Errorf("n.repo.CreateMany: %w", err)
	}

	return nil
}

func (n *news) GetByID(id string) (*entity.News, error) {
	news, err := n.Repo.GetByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("n.repo.GetByID: %w", err)
	}

	return news, nil
}

type (
	Query   = repo.Query
	Options = repo.Options
)

func (n *news) GetByQuery(query Query, opts Options) ([]entity.News, error) {
	news, err := n.Repo.GetByQuery(context.Background(), query, opts)
	if err != nil {
		return nil, fmt.Errorf("n.repo.GetByQuery: %w", err)
	}

	return news, nil
}

func (n *news) Parse(query string) error {
	err := n.Producer.Produce(n.Exchange, n.RoutingKey, rabbitmq.Message{
		ContentType:  "text/plain",
		DeliveryMode: 2,
		Body:         []byte(query),
	})

	if err != nil {
		return fmt.Errorf("n.AMQP.Producer.Produce: %w", err)
	}

	return nil
}