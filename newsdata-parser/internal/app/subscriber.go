package app

import (
	"fmt"

	"github.com/qsoulior/news/newsdata-parser/internal/service"
	"github.com/qsoulior/news/newsdata-parser/internal/transport/amqp"
	"github.com/qsoulior/news/newsdata-parser/pkg/rabbitmq"
	"github.com/rs/zerolog"
)

type SubscriberConfig struct {
	Logger  *zerolog.Logger
	Service service.News
	AMQP    struct {
		Queue    string
		Consumer *rabbitmq.Consumer
	}
}

type subscriber struct {
	SubscriberConfig
}

func NewSubscriber(cfg SubscriberConfig) *subscriber {
	return &subscriber{cfg}
}

func (s *subscriber) Run() error {
	handler := amqp.NewNews(amqp.NewsConfig{
		Logger:  s.Logger,
		Service: s.Service,
	})

	err := s.AMQP.Consumer.Consume(s.AMQP.Queue, handler.Handle)
	return fmt.Errorf("c.amqp.Consume: %w", err)
}
