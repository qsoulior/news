package service

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/aggregator/repo"
	"github.com/rabbitmq/amqp091-go"
)

type NewsConfig struct {
	AMQP struct {
		Producer   rabbitmq.Producer
		Exchange   string
		RoutingKey string
	}
	repo repo.News
}

type news struct {
	NewsConfig
}

func NewNews(cfg NewsConfig) News {
	return &news{cfg}
}

func (n *news) Create(news entity.News) error {
	if err := n.repo.Create(context.Background(), news); err != nil {
		return fmt.Errorf("n.repo.Create: %w", err)
	}

	return nil
}

func (n *news) CreateMany(news []entity.News) error {
	if err := n.repo.CreateMany(context.Background(), news); err != nil {
		return fmt.Errorf("n.repo.CreateMany: %w", err)
	}

	return nil
}

func (n *news) GetByID(id string) (*entity.News, error) {
	news, err := n.repo.GetByID(context.Background(), id)
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
	news, err := n.repo.GetByQuery(context.Background(), query, opts)
	if err != nil {
		return nil, fmt.Errorf("n.repo.GetByQuery: %w", err)
	}

	return news, nil
}

func (n *news) Parse(query string) error {
	err := n.AMQP.Producer.Produce(n.AMQP.Exchange, n.AMQP.RoutingKey, amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte(query),
	})

	if err != nil {
		return fmt.Errorf("n.AMQP.Producer.Produce: %w", err)
	}

	return nil
}
