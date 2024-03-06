package app

import (
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/consumer"
	"github.com/qsoulior/news/parser/transport/amqp"
)

func Run() {
	amqpRouter := amqp.NewRouter(nil, nil)
	rmqConn, err := rabbitmq.New(rabbitmq.Config{})
	if err != nil {
	}

	rmqConsumer := consumer.New(rmqConn, amqpRouter)

	err = rmqConsumer.Consume("")
}
