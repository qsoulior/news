package app

import (
	"time"

	"github.com/qsoulior/news/iz-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/rs/zerolog"
)

func Run(cfg *Config) {
	appID := "iz"

	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})
	logger := zerolog.New(out).With().Timestamp().Logger()

	client := httpclient.New(
		httpclient.URL(cfg.API.URL),
		httpclient.Headers(map[string]string{
			"Referer": cfg.API.URL,
		}),
	)

	searchParser := service.NewNewsSearch(appID, client, &logger)
	archiveParser := service.NewNewsArchive(appID, client, &logger)
	feedParser := service.NewNewsFeed(appID, client, &logger)

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
