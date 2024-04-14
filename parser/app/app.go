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
	"github.com/qsoulior/news/parser/internal/worker"
	"github.com/qsoulior/news/parser/pkg/redis"
	"github.com/rs/zerolog"
)

var wg sync.WaitGroup

type Options struct {
	Logger        *zerolog.Logger
	SearchParser  service.Parser
	ArchiveParser service.Parser
	FeedParser    service.Parser
}

func Run(cfg *Config, opts *Options) {
	zerolog.DurationFieldUnit = time.Second
	var log zerolog.Logger

	if opts.Logger != nil {
		log = *opts.Logger
	} else {
		out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.RFC3339
		})
		log = zerolog.New(out)
	}

	logger := log.With().Timestamp().Logger()

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
		redisLog.Error().Err(err).Send()
		return
	}
	redisLog.Info().Str("url", cfg.Redis.URL).Msg("started")

	defer func() {
		// redis graceful shutdown
		err := redis.Close()
		if err != nil {
			redisLog.Error().Err(err).Msg("graceful shutdown")
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
		rmqLog.Error().Err(err).Send()
		return
	}
	rmqLog.Info().Msg("started")

	// rabbit producer
	rmqProducer := producer.New(rmqConn)

	// search consumer
	if opts.SearchParser == nil {
		logger.Error().Msg("search parser is nil")
		return
	}

	searchService := service.NewNews(service.NewsConfig{
		Repo:   newsRepo,
		Parser: opts.SearchParser,

		Producer:   rmqProducer,
		Exchange:   "",
		RoutingKey: "news",
		AppID:      cfg.ID,
	})
	runSearcher(sigCtx, &logger, searchService, rmqConn, queue)

	// archive worker
	if opts.ArchiveParser != nil {
		archiveService := service.NewNews(service.NewsConfig{
			Repo:   newsRepo,
			Parser: opts.ArchiveParser,

			Producer:   rmqProducer,
			Exchange:   "",
			RoutingKey: "news",
			AppID:      cfg.ID,
		})
		pageService := service.NewPage(service.PageConfig{
			Repo: pageRepo,
		})
		runArchiver(sigCtx, &logger, archiveService, pageService)
	}

	// feed worker
	if opts.FeedParser != nil {
		feedService := service.NewNews(service.NewsConfig{
			Repo:   newsRepo,
			Parser: opts.FeedParser,

			Producer:   rmqProducer,
			Exchange:   "",
			RoutingKey: "news",
			AppID:      cfg.ID,
		})
		runFeeder(sigCtx, &logger, feedService)
	}

	// release worker
	releaseService := service.NewNews(service.NewsConfig{
		Repo:   newsRepo,
		Parser: nil,

		Producer:   rmqProducer,
		Exchange:   "",
		RoutingKey: "news",
		AppID:      cfg.ID,
	})
	runReleaser(sigCtx, &logger, releaseService)

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

		// rabbit conn graceful shutdown
		err := rmqConn.Close()
		if err != nil {
			logger.Error().Err(err).Msg("graceful shutdown")
			return
		}
		logger.Info().Msg("graceful shutdown")
	}()

	return rmqConn, queue.Name, nil
}

func runSearcher(ctx context.Context, logger *zerolog.Logger, news service.News, conn *rabbitmq.Connection, queue string) {
	log := logger.With().Str("module", "searcher").Logger()

	amqpRouter := amqp.NewRouter(&log, news)
	rmqConsumer := consumer.New(conn, amqpRouter, consumer.Ack(false))

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
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				log.Info().Msg("graceful shutdown")
			}
		}
	}()

	log.Info().Msg("started")
}

func runWorker(ctx context.Context, logger *zerolog.Logger, worker worker.Worker) {
	go func() {
		wg.Add(1)
		defer wg.Done()
		err := worker.Run(ctx)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
		logger.Info().Msg("graceful shutdown")
	}()

	logger.Info().Msg("started")
}

func runArchiver(ctx context.Context, logger *zerolog.Logger, news service.News, page service.Page) {
	log := logger.With().Str("module", "archiver").Logger()
	worker := worker.NewArchive(5*time.Second, &log, news, page)

	runWorker(ctx, &log, worker)
}

func runReleaser(ctx context.Context, logger *zerolog.Logger, news service.News) {
	log := logger.With().Str("module", "releaser").Logger()
	worker := worker.NewRelease(10*time.Minute, &log, news)

	runWorker(ctx, &log, worker)
}

func runFeeder(ctx context.Context, logger *zerolog.Logger, news service.News) {
	log := logger.With().Str("module", "feeder").Logger()
	worker := worker.NewFeed(15*time.Second, &log, news)

	runWorker(ctx, &log, worker)
}
