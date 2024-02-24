package amqp

import (
	"github.com/qsoulior/news/newsdata-parser/internal/service"
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
	_, err := n.Service.Parse(string(msg.Body), "")
	if err != nil {
		n.Logger.Error().Err(err).Msg("")
	}
}
