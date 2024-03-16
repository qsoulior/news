package handler

import (
	"context"

	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/parser/internal/service"
	"github.com/rs/zerolog"
)

type NewsConfig struct {
	Logger  *zerolog.Logger
	Service service.News
}

type news struct {
	NewsConfig
}

func NewNews(cfg NewsConfig) *news {
	return &news{cfg}
}

func (n *news) Handle(ctx context.Context, msg *rabbitmq.Delivery) {
	_, err := n.Service.Parse(ctx, string(msg.Body), "")
	if err != nil {
		n.Logger.Error().Err(err).Msg("")
	}
}
