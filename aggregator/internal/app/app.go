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

func Run(cfg *Config) {
	// notify context
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	zerolog.DurationFieldUnit = time.Second
	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})

	logger := zerolog.New(out).With().Timestamp().Logger()
	ctx := logger.WithContext(sigCtx)

	// mongo client
	mongoLog := logger.With().Str("module", "mongo").Logger()
	mongo, err := mongodb.New(mongoLog.WithContext(ctx), &mongodb.Config{
		URI:          cfg.MongoDB.URI,
		AttemptCount: 5,
		AttemptDelay: 5 * time.Second,
	})
	if err != nil {
		mongoLog.Error().Err(err).Send()
		return
	}
	mongoLog.Info().Str("uri", cfg.MongoDB.URI).Msg("started")

	defer func() {
		// mongo graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := mongo.Disconnect(ctx)
		if err != nil {
			mongoLog.Error().Err(err).Msg("graceful shutdown")
			return
		}
		mongoLog.Info().Msg("graceful shutdown")
	}()

	db := mongo.Client.Database("app")
	newsRepo := repo.NewNewsMongo(db)

	// rabbit connection
	rmqLog := logger.With().Str("module", "rmq").Logger()
	rmqConn, err := runRMQ(rmqLog.WithContext(ctx), cfg.RabbitMQ)
	if err != nil {
		rmqLog.Error().Err(err).Send()
		return
	}
	rmqLog.Info().Msg("started")

	// rabbit producer
	rmqProducer := producer.New(rmqConn)
	newsService := service.NewNews(service.NewsConfig{
		Producer:   rmqProducer,
		Exchange:   "query",
		RoutingKey: "",
		Repo:       newsRepo,
	})

	// rabbit consumer
	runConsumer(ctx, newsService, rmqConn)

	// http server
	runServer(ctx, newsService, cfg.HTTP)

	wg.Wait()
}

func runRMQ(ctx context.Context, cfg ConfigRabbitMQ) (*rabbitmq.Connection, error) {
	logger := zerolog.Ctx(ctx)
	rmqConn, err := rabbitmq.New(ctx, &rabbitmq.Config{
		URL:          cfg.URL,
		AttemptCount: 5,
		AttemptDelay: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("rabbitmq.New: %w", err)
	}

	err = rmqConn.Ch.ExchangeDeclare("query", "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("rmqConn.Ch.ExchangeDeclare: %w", err)
	}

	_, err = rmqConn.Ch.QueueDeclare("news", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("rmqConn.Ch.QueueDeclare: %w", err)
	}

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		<-ctx.Done()
		logger.Info().Msg("term signal accepted")

		// rabbit conn graceful shutdown
		err := rmqConn.Close()
		if err != nil {
			logger.Error().Err(err).Msg("graceful shutdown")
			return
		}
		logger.Info().Msg("graceful shutdown")
	}(ctx)

	return rmqConn, nil
}

func runConsumer(ctx context.Context, news service.News, conn *rabbitmq.Connection) {
	log := zerolog.Ctx(ctx).With().Str("module", "consumer").Logger()

	amqpRouter := amqp.NewRouter(&log, news)
	rmqConsumer := consumer.New(conn, amqpRouter, consumer.Ack(false))

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for timer := time.NewTimer(0); ; timer.Reset(5 * time.Second) {
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				err := rmqConsumer.Consume(ctx, "news")
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				log.Info().Msg("graceful shutdown")
			}
		}
	}(ctx)

	log.Info().Msg("started")
}

func runServer(ctx context.Context, news service.News, cfg ConfigHTTP) {
	log := zerolog.Ctx(ctx).With().Str("module", "server").Logger()

	httpRouter := http.NewRouter(news)
	httpServer := httpserver.New(httpRouter, httpserver.Addr(cfg.Host, cfg.Port))

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		httpServer.Start(ctx)

		select {
		case <-ctx.Done():
			log.Info().Msg("term signal accepted")
		case err := <-httpServer.Err():
			log.Error().Err(err).Send()
		}

		// http server graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := httpServer.Stop(ctx)
		if err != nil {
			log.Error().Err(err).Msg("graceful shutdown")
			return
		}
		log.Info().Msg("graceful shutdown")
	}(log.WithContext(ctx))

	log.Info().Msg("started")
}
