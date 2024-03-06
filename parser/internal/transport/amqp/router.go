package amqp

import (
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/parser/internal/service"
	"github.com/qsoulior/news/parser/internal/transport/amqp/handler"
	"github.com/rs/zerolog"
)

type Config struct {
	Logger      *zerolog.Logger
	NewsService service.News
}

func NewRouter(logger *zerolog.Logger, service service.News) rabbitmq.Handler {
	news := handler.NewNews(handler.NewsConfig{
		Logger:  logger,
		Service: service,
	})

	return func(msg *rabbitmq.Delivery) {
		news.Handle(msg)
	}
}
