package amqp

import (
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/consumer"
	"github.com/qsoulior/news/aggregator/service"
	"github.com/qsoulior/news/aggregator/transport/amqp/handler"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Config struct {
	Logger      *zerolog.Logger
	NewsService service.News
}

func NewRouter(logger *zerolog.Logger, service service.News) consumer.Handler {
	news := handler.NewNews(handler.NewsConfig{
		Logger:  logger,
		Service: service,
	})

	return func(msg *amqp091.Delivery) {
		news.Handle(msg)
	}
}
