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

type Config struct {
	ID            string
	Logger        *zerolog.Logger
	SearchParser  service.Parser
	ArchiveParser service.Parser
	FeedParser    service.Parser
}

type Options struct {
	RabbitURL    string
	RedisURL     string
	ReleaseDelay *time.Duration
	ArchiveDelay *time.Duration
	FeedDelay    *time.Duration
}

var (
	DefaultReleaseDelay = 15 * time.Minute
	DefaultArchiveDelay = 5 * time.Second
	DefaultFeedDelay    = 1 * time.Minute
)

func (o *Options) setDefault() {
	if o.ReleaseDelay == nil {
		o.ReleaseDelay = &DefaultReleaseDelay
	}
	if o.ArchiveDelay == nil {
		o.ArchiveDelay = &DefaultArchiveDelay
	}
	if o.FeedDelay == nil {
		o.FeedDelay = &DefaultFeedDelay
	}
}

func Run(cfg *Config, opts *Options) {
	opts.setDefault()

	// notify context
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	zerolog.DurationFieldUnit = time.Second
	var log zerolog.Logger

	if cfg.Logger != nil {
		log = *cfg.Logger
	} else {
		out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.RFC3339
		})
		log = zerolog.New(out)
	}

	logger := log.With().Timestamp().Logger()
	ctx := logger.WithContext(sigCtx)

	// redis client
	redisLog := logger.With().Str("module", "redis").Logger()
	redis, err := redis.New(redisLog.WithContext(ctx), &redis.RedisConfig{
		URL:          opts.RedisURL,
		AttemptCount: 5,
		AttemptDelay: 10 * time.Second,
	})
	if err != nil {
		redisLog.Error().Err(err).Send()
		return
	}
	redisLog.Info().Str("url", opts.RedisURL).Msg("started")

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
	rmqConn, queue, err := runRMQ(rmqLog.WithContext(ctx), opts.RabbitURL, "query."+cfg.ID)
	if err != nil {
		rmqLog.Error().Err(err).Send()
		return
	}
	rmqLog.Info().Msg("started")

	// rabbit producer
	rmqProducer := producer.New(rmqConn)

	// search consumer
	if cfg.SearchParser == nil {
		logger.Error().Msg("search parser is nil")
		return
	}

	searchService := service.NewNews(service.NewsConfig{
		Repo:   newsRepo,
		Parser: cfg.SearchParser,

		Producer:   rmqProducer,
		Exchange:   "",
		RoutingKey: "news",
		AppID:      cfg.ID,
	})
	runSearcher(ctx, searchService, rmqConn, queue)

	// archive worker
	if cfg.ArchiveParser != nil {
		archiveService := service.NewNews(service.NewsConfig{
			Repo:   newsRepo,
			Parser: cfg.ArchiveParser,

			Producer:   rmqProducer,
			Exchange:   "",
			RoutingKey: "news",
			AppID:      cfg.ID,
		})
		pageService := service.NewPage(service.PageConfig{
			Repo: pageRepo,
		})
		runArchiver(ctx, archiveService, pageService, *opts.ArchiveDelay)
	}

	// feed worker
	if cfg.FeedParser != nil {
		feedService := service.NewNews(service.NewsConfig{
			Repo:   newsRepo,
			Parser: cfg.FeedParser,

			Producer:   rmqProducer,
			Exchange:   "",
			RoutingKey: "news",
			AppID:      cfg.ID,
		})
		runFeeder(ctx, feedService, *opts.FeedDelay)
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
	runReleaser(ctx, releaseService, *opts.ReleaseDelay)

	wg.Wait()
}

func runRMQ(ctx context.Context, url string, queueName string) (*rabbitmq.Connection, string, error) {
	logger := zerolog.Ctx(ctx)
	rmqConn, err := rabbitmq.New(ctx, &rabbitmq.Config{
		URL:          url,
		AttemptCount: 5,
		AttemptDelay: 10 * time.Second,
	})
	if err != nil {
		return nil, "", fmt.Errorf("rabbitmq.New: %w", err)
	}

	err = rmqConn.Ch.ExchangeDeclare("query", "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, "", fmt.Errorf("rmqConn.Ch.ExchangeDeclare: %w", err)
	}

	queue, err := rmqConn.Ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, "", fmt.Errorf("rmqConn.Ch.QueueDeclare: %w", err)
	}

	err = rmqConn.Ch.QueueBind(queue.Name, "", "query", false, nil)
	if err != nil {
		return nil, "", fmt.Errorf("rmqConn.Ch.QueueBind: %w", err)
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

	return rmqConn, queue.Name, nil
}

func runSearcher(ctx context.Context, news service.News, conn *rabbitmq.Connection, queue string) {
	log := zerolog.Ctx(ctx).With().Str("module", "searcher").Logger()

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
				err := rmqConsumer.Consume(ctx, queue)
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

func runWorker(ctx context.Context, worker worker.Worker) {
	logger := zerolog.Ctx(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		err := worker.Run(ctx)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
		logger.Info().Msg("graceful shutdown")
	}(ctx)

	logger.Info().Msg("started")
}

func runArchiver(ctx context.Context, news service.News, page service.Page, delay time.Duration) {
	log := zerolog.Ctx(ctx).With().Str("module", "archiver").Logger()
	worker := worker.NewArchive(delay, &log, news, page)

	runWorker(log.WithContext(ctx), worker)
}

func runReleaser(ctx context.Context, news service.News, delay time.Duration) {
	log := zerolog.Ctx(ctx).With().Str("module", "releaser").Logger()
	worker := worker.NewRelease(delay, &log, news)

	runWorker(log.WithContext(ctx), worker)
}

func runFeeder(ctx context.Context, news service.News, delay time.Duration) {
	log := zerolog.Ctx(ctx).With().Str("module", "feeder").Logger()
	worker := worker.NewFeed(delay, &log, news)

	runWorker(log.WithContext(ctx), worker)
}
