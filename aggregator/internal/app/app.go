package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qsoulior/news/aggregator/internal/repo"
	"github.com/qsoulior/news/aggregator/internal/service"
	"github.com/qsoulior/news/aggregator/internal/transport/amqp"
	"github.com/qsoulior/news/aggregator/internal/transport/http"
	"github.com/qsoulior/news/aggregator/pkg/httpserver"
	"github.com/qsoulior/news/aggregator/pkg/mongodb"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/consumer"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/producer"
	"github.com/rs/zerolog"
)

func Run(cfg *Config, logger *zerolog.Logger) {
	mongo, err := mongodb.New(mongodb.MongoConfig{
		URL:          cfg.MongoDB.URL,
		AttemptCount: 5,
		AttemptDelay: 5 * time.Second,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	db := mongo.Client.Database("app")
	newsRepo := repo.NewNewsMongo(db)

	rmqConn, err := rabbitmq.New(rabbitmq.Config{
		URL:          cfg.RabbitMQ.URL,
		AttemptCount: 5,
		AttemptDelay: 5 * time.Second,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	err = rmqConn.Ch.ExchangeDeclare("queries", "fanout", false, false, false, false, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	queue, err := rmqConn.Ch.QueueDeclare("news", true, false, false, false, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	rmqProducer := producer.New(rmqConn)
	newsService := service.NewNews(service.NewsConfig{
		Producer:   rmqProducer,
		Exchange:   "queries",
		RoutingKey: "",
		Repo:       newsRepo,
	})

	httpRouter := http.NewRouter(logger, newsService)
	httpServer := httpserver.New(httpRouter, httpserver.Addr(cfg.HTTP.Host, cfg.HTTP.Port))

	amqpRouter := amqp.NewRouter(logger, newsService)
	rmqConsumer := consumer.New(rmqConn, amqpRouter)

	httpServer.Start()
	go func() {
		err := rmqConsumer.Consume(queue.Name)
		logger.Error().Err(err).Msg("")
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		logger.Info().Msg("graceful shutdown")
	case err := <-httpServer.Err():
		logger.Error().Err(err).Msg("")
	case err := <-rmqConn.Err():
		logger.Error().Err(err).Msg("")
	}

	// shutdown
	err = httpServer.Stop(10 * time.Second)
	if err != nil {
		logger.Error().Err(err).Msg("")
	}

	err = rmqConn.Close(10 * time.Second)
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
}
