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

func Run(cfg *Config, consumerParser service.Parser, workerParser service.Parser) {
	out := zerolog.NewConsoleWriter()
	logger := zerolog.New(out).With().Timestamp().Logger()

	// notify context
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// redis client
	redisLog := logger.With().Str("name", "redis").Logger()
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
	rmqLog := logger.With().Str("name", "rmq").Logger()
	rmqConn, queue, err := runRMQ(sigCtx, &rmqLog, cfg.RabbitMQ)
	if err != nil {
		rmqLog.Error().Err(err).Msg("")
		return
	}
	rmqLog.Info().Msg("started")

	// rabbit producer
	rmqProducer := producer.New(rmqConn)
	workerService := service.NewNews(service.NewsConfig{
		Repo:   newsRepo,
		Parser: workerParser,

		Producer:   rmqProducer,
		Exchange:   "",
		RoutingKey: "news",
	})

	pageService := service.NewPage(service.PageConfig{
		Repo: pageRepo,
	})

	// rabbit consumer
	consumerService := service.NewNews(service.NewsConfig{
		Repo:   newsRepo,
		Parser: consumerParser,

		Producer:   rmqProducer,
		Exchange:   "",
		RoutingKey: "news",
	})

	consumerLog := logger.With().Str("name", "consumer").Logger()
	runConsumer(sigCtx, &consumerLog, consumerService, rmqConn, queue)
	consumerLog.Info().Msg("started")

	// worker
	workerLog := logger.With().Str("name", "worker").Logger()
	runWorker(sigCtx, &workerLog, workerService, pageService)
	workerLog.Info().Msg("started")

	wg.Wait()
}

func runRMQ(ctx context.Context, logger *zerolog.Logger, cfg ConfigRabbitMQ) (*rabbitmq.Connection, string, error) {
	rmqConn, err := rabbitmq.New(ctx, &rabbitmq.Config{
		URL:          cfg.URL,
		AttemptCount: 5,
		AttemptDelay: 10 * time.Second,
		Logger:       logger,
	})
	if err != nil {
		return nil, "", fmt.Errorf("rabbitmq.New: %w", err)
	}

	queue, err := rmqConn.Ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		return nil, "", fmt.Errorf("rmqConn.Ch.QueueDeclare: %w", err)
	}

	err = rmqConn.Ch.QueueBind(queue.Name, "", "queries", false, nil)
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
	worker := NewWorker(WorkerConfig{
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
