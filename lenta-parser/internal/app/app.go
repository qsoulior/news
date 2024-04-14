package app

import (
	"time"

	"github.com/qsoulior/news/lenta-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/rs/zerolog"
)

func Run(cfg *Config) {
	appID := "lenta"

	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})

	logger := zerolog.New(out).With().Timestamp().Logger()

	client := httpclient.New()
	searchParser := service.NewNewsSearch(appID, cfg.API.Search.URL, client, &logger)
	archiveParser := service.NewNewsArchive(appID, cfg.API.Archive.URL, client, &logger)
	feedParser := service.NewNewsFeed(appID, cfg.API.Feed.URL, client, &logger)

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
			Logger:        &logger,
		},
	)
}
