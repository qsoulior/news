package amqp

import (
	"context"

	"github.com/qsoulior/news/aggregator/internal/service"
	"github.com/qsoulior/news/aggregator/internal/transport/amqp/handler"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
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

	return func(ctx context.Context, msg *rabbitmq.Delivery) {
		logger.Info().Str("app", msg.AppId).Str("id", msg.MessageId).Msg("message accepted")
		news.Handle(ctx, msg)
	}
}
