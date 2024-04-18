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

func (n *news) Create(ctx context.Context, news entity.News) error {
	if err := n.Repo.ReplaceOrCreate(ctx, news); err != nil {
		return fmt.Errorf("n.repo.Create: %w", err)
	}

	return nil
}

func (n *news) CreateMany(ctx context.Context, news []entity.News) error {
	if err := n.Repo.CreateMany(ctx, news); err != nil {
		return fmt.Errorf("n.repo.CreateMany: %w", err)
	}

	return nil
}

func (n *news) Get(ctx context.Context, id string) (*entity.News, error) {
	news, err := n.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("n.repo.GetByID: %w", err)
	}

	return news, nil
}

func (n *news) GetHead(ctx context.Context, query Query, opts Options) ([]entity.NewsHead, int, error) {
	if opts.raw.Sort.IsRelevance() && (query.Text == "" || (query.Text != "" && query.Title)) {
		opts.SetSort(0)
	}

	if query.DateTo != nil {
		dateTo := query.DateTo.AddDate(0, 0, 1)
		query.DateTo = &dateTo
	}

	news, count, err := n.Repo.GetByQuery(ctx, query, opts.raw)
	if err != nil {
		return nil, 0, fmt.Errorf("n.repo.GetByQuery: %w", err)
	}

	return news, count, nil
}

func (n *news) SendToParse(ctx context.Context, query string) error {
	err := n.Producer.Produce(ctx, n.Exchange, n.RoutingKey, rabbitmq.Message{
		ContentType:  "text/plain",
		DeliveryMode: 2,
		Body:         []byte(query),
	})

	if err != nil {
		return fmt.Errorf("n.AMQP.Producer.Produce: %w", err)
	}

	return nil
}
