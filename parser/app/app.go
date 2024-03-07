package app

import (
	"context"

	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/consumer"
	"github.com/qsoulior/news/parser/internal/service"
	"github.com/qsoulior/news/parser/internal/transport/amqp"
)

func Run(news service.News, page service.Page) {
	amqpRouter := amqp.NewRouter(nil, news)
	rmqConn, err := rabbitmq.New(&rabbitmq.Config{})
	if err != nil {
	}

	rmqConsumer := consumer.New(rmqConn, amqpRouter)

	err = rmqConsumer.Consume(context.Background(), "")
}
