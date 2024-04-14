package app

import (
	"context"
	"time"

	"github.com/go-rod/rod"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/ria-parser/internal/service"
	"github.com/rs/zerolog"
)

func Run(cfg *Config) {
	appID := "ria"

	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})
	logger := zerolog.New(out).With().Timestamp().Logger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	browser := rod.New().Context(ctx)
	if err := browser.Connect(); err != nil {
		logger.Fatal().Err(err).Msg("failed to init browser")
	}

	defer func() {
		if err := browser.Close(); err != nil {
			logger.Fatal().Err(err).Msg("failed to close browser")
		}
	}()

	client := httpclient.New()

	searchParser := service.NewNewsSearch(appID, client, cfg.API.URL, browser, &logger)
	archiveParser := service.NewNewsArchive(appID, client, cfg.API.URL, browser, &logger)
	feedParser := service.NewNewsFeed(appID, client, cfg.API.URL, &logger)

	app.Run(
		&app.Config{
			ID:       appID,
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		&app.Options{
			SearchParser:  searchParser,
			ArchiveParser: archiveParser,
			FeedParser:    feedParser,
		},
	)
}
