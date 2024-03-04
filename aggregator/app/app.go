package app

import (
	"github.com/qsoulior/news/aggregator/pkg/httpserver"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/consumer"
	"github.com/qsoulior/news/aggregator/transport/amqp"
	"github.com/qsoulior/news/aggregator/transport/http"
)

func Run(cfg *Config) {
	httpRouter := http.NewRouter(nil, nil)
	httpServer := httpserver.New(httpRouter)

	amqpRouter := amqp.NewRouter(nil, nil)
	rmqConn, err := rabbitmq.New(rabbitmq.Config{})
	if err != nil {
	}

	rmqConsumer := consumer.New(rmqConn, amqpRouter)

	err = httpServer.Start()
	err = rmqConsumer.Consume("")

}
