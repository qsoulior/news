package handler

import (
	"encoding/json"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/internal/service"
	"github.com/rabbitmq/amqp091-go"
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

func (n *news) Handle(msg *amqp091.Delivery) {
	var news entity.News
	err := json.Unmarshal(msg.Body, &news)
	if err != nil {
		n.Logger.Error().Err(err).Msg("")
		return
	}

	err = n.Service.Create(news)
	if err != nil {
		n.Logger.Error().Err(err).Msg("")
	}
}
