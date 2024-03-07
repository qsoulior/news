package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
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

var wg sync.WaitGroup

func Run(cfg *Config, logger *zerolog.Logger) {
	// notify context
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// mongo client initialization
	mongo, err := mongodb.New(sigCtx, &mongodb.Config{
		URL:          cfg.MongoDB.URL,
		AttemptCount: 5,
		AttemptDelay: 5 * time.Second,
		Logger:       logger,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := mongo.Disconnect(ctx)
		logger.Error().Err(err).Msg("")
	}()

	db := mongo.Client.Database("app")
	newsRepo := repo.NewNewsMongo(db)

	wg.Add(3)

	// rabbit connection initialization
	rmqConn, err := runRMQ(sigCtx, logger, cfg.RabbitMQ)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	// rabbit producer initialization
	rmqProducer := producer.New(rmqConn)
	newsService := service.NewNews(service.NewsConfig{
		Producer:   rmqProducer,
		Exchange:   "queries",
		RoutingKey: "",
		Repo:       newsRepo,
	})

	runConsumer(sigCtx, logger, newsService, rmqConn)
	runServer(sigCtx, logger, newsService, cfg.HTTP)

	wg.Wait()
}

func runRMQ(ctx context.Context, logger *zerolog.Logger, cfg ConfigRabbitMQ) (*rabbitmq.Connection, error) {
	rmqLog := logger.With().Str("name", "rmq").Logger()
	rmqConn, err := rabbitmq.New(&rabbitmq.Config{
		URL:          cfg.URL,
		AttemptCount: 5,
		AttemptDelay: 5 * time.Second,
		Logger:       &rmqLog,
	})
	if err != nil {
		return nil, fmt.Errorf("rabbitmq.New: %w", err)
	}

	err = rmqConn.Ch.ExchangeDeclare("queries", "fanout", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("rmqConn.Ch.ExchangeDeclare: %w", err)
	}

	_, err = rmqConn.Ch.QueueDeclare("news", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("rmqConn.Ch.QueueDeclare: %w", err)
	}

	go func() {
		defer wg.Done()
		select {
		case err := <-rmqConn.Err():
			rmqLog.Error().Err(err).Msg("")
		case <-ctx.Done():
			rmqLog.Info().Msg("term signal accepted")
		}

		err := rmqConn.Close()
		if err != nil {
			rmqLog.Error().Err(err).Msg("graceful shutdown")
		}
	}()

	return rmqConn, nil
}

func runConsumer(ctx context.Context, logger *zerolog.Logger, news service.News, conn *rabbitmq.Connection) {
	consumerLog := logger.With().Str("name", "consumer").Logger()
	amqpRouter := amqp.NewRouter(&consumerLog, news)
	rmqConsumer := consumer.New(conn, amqpRouter)

	go func() {
		defer wg.Done()
		err := rmqConsumer.Consume(ctx, "news")
		if err != nil {
			consumerLog.Error().Err(err).Msg("")
		}
	}()
}

func runServer(ctx context.Context, logger *zerolog.Logger, news service.News, cfg ConfigHTTP) {
	serverLog := logger.With().Str("name", "server").Logger()
	httpRouter := http.NewRouter(&serverLog, news)
	httpServer := httpserver.New(httpRouter, httpserver.Addr(cfg.Host, cfg.Port))

	go func() {
		defer wg.Done()
		httpServer.Start()

		select {
		case err := <-httpServer.Err():
			serverLog.Error().Err(err).Msg("")
		case <-ctx.Done():
			serverLog.Info().Msg("term signal accepted")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := httpServer.Stop(ctx)
		if err != nil {
			serverLog.Error().Err(err).Msg("graceful shutdown")
		}
	}()
}
