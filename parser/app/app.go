package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/qsoulior/news/aggregator/pkg/rabbitmq"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/consumer"
	"github.com/qsoulior/news/aggregator/pkg/rabbitmq/producer"
	"github.com/qsoulior/news/parser/internal/repo"
	"github.com/qsoulior/news/parser/internal/service"
	"github.com/qsoulior/news/parser/internal/transport/amqp"
	"github.com/qsoulior/news/parser/pkg/redis"
	"github.com/rs/zerolog"
)

var wg sync.WaitGroup

func Run(cfg *Config, searchParser service.Parser, feedParser service.Parser) {
	zerolog.DurationFieldUnit = time.Second
	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})
	logger := zerolog.New(out).With().Timestamp().Logger()

	// notify context
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// redis client
	redisLog := logger.With().Str("module", "redis").Logger()
	redis, err := redis.New(sigCtx, &redis.RedisConfig{
		URL:          cfg.Redis.URL,
		AttemptCount: 5,
		AttemptDelay: 10 * time.Second,
		Logger:       &redisLog,
	})
	if err != nil {
		redisLog.Error().Err(err).Msg("")
		return
	}
	redisLog.Info().Str("url", cfg.Redis.URL).Msg("started")

	defer func() {
		err := redis.Close()
		if err != nil {
			redisLog.Error().Err(err).Msg("")
			return
		}
		redisLog.Info().Msg("graceful shutdown")
	}()

	newsRepo := repo.NewNewsRedis(redis)
	pageRepo := repo.NewPageRedis(redis)

	// rabbit connection
	rmqLog := logger.With().Str("module", "rmq").Logger()
	rmqConn, queue, err := runRMQ(sigCtx, &rmqLog, cfg.RabbitMQ.URL, "query."+cfg.ID)
	if err != nil {
		rmqLog.Error().Err(err).Msg("")
		return
	}
	rmqLog.Info().Msg("started")

	// rabbit producer
	rmqProducer := producer.New(rmqConn)
	feedService := service.NewNews(service.NewsConfig{
		Repo:   newsRepo,
		Parser: feedParser,

		Producer:   rmqProducer,
		Exchange:   "",
		RoutingKey: "news",
		AppID:      cfg.ID,
	})

	pageService := service.NewPage(service.PageConfig{
		Repo: pageRepo,
	})

	// rabbit consumer
	searchService := service.NewNews(service.NewsConfig{
		Repo:   newsRepo,
		Parser: searchParser,

		Producer:   rmqProducer,
		Exchange:   "",
		RoutingKey: "news",
		AppID:      cfg.ID,
	})

	consumerLog := logger.With().Str("module", "consumer").Logger()
	runConsumer(sigCtx, &consumerLog, searchService, rmqConn, queue)
	consumerLog.Info().Msg("started")

	// worker
	workerLog := logger.With().Str("module", "worker").Logger()
	runWorker(sigCtx, &workerLog, feedService, pageService)
	workerLog.Info().Msg("started")

	// releaser
	releaserLog := logger.With().Str("module", "releaser").Logger()
	runReleaser(sigCtx, &releaserLog, feedService)
	releaserLog.Info().Msg("started")

	wg.Wait()
}

func runRMQ(ctx context.Context, logger *zerolog.Logger, url string, queueName string) (*rabbitmq.Connection, string, error) {
	rmqConn, err := rabbitmq.New(ctx, &rabbitmq.Config{
		URL:          url,
		AttemptCount: 5,
		AttemptDelay: 10 * time.Second,
		Logger:       logger,
	})
	if err != nil {
		return nil, "", fmt.Errorf("rabbitmq.New: %w", err)
	}

	queue, err := rmqConn.Ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, "", fmt.Errorf("rmqConn.Ch.QueueDeclare: %w", err)
	}

	err = rmqConn.Ch.QueueBind(queue.Name, "", "query", false, nil)
	if err != nil {
		return nil, "", fmt.Errorf("rmqConn.Ch.QueueBind: %w", err)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		<-ctx.Done()
		logger.Info().Msg("term signal accepted")

		err := rmqConn.Close()
		if err != nil {
			logger.Error().Err(err).Msg("graceful shutdown")
		}
	}()

	return rmqConn, queue.Name, nil
}

func runConsumer(ctx context.Context, logger *zerolog.Logger, news service.News, conn *rabbitmq.Connection, queue string) {
	amqpRouter := amqp.NewRouter(logger, news)
	rmqConsumer := consumer.New(conn, amqpRouter)

	go func() {
		wg.Add(1)
		defer wg.Done()
		for timer := time.NewTimer(0); ; timer.Reset(5 * time.Second) {
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				err := rmqConsumer.Consume(ctx, queue)
				if err == nil {
					logger.Info().Msg("graceful shutdown")
					return
				}
				logger.Error().Err(err).Msg("")
			}
		}
	}()
}

func runWorker(ctx context.Context, logger *zerolog.Logger, news service.News, page service.Page) {
	worker := newWorker(workerConfig{
		Delay:  5 * time.Second,
		Logger: logger,
		News:   news,
		Page:   page,
	})

	go func() {
		wg.Add(1)
		defer wg.Done()
		err := worker.Run(ctx)
		if err == nil {
			logger.Info().Msg("graceful shutdown")
			return
		}
		logger.Error().Err(err).Msg("")
	}()
}

func runReleaser(ctx context.Context, logger *zerolog.Logger, news service.News) {
	releaser := newReleaser(releaserConfig{
		Delay:  30 * time.Second,
		Logger: logger,
		News:   news,
	})

	go func() {
		wg.Add(1)
		defer wg.Done()
		releaser.Run(ctx)
		logger.Info().Msg("graceful shutdown")
	}()
}
