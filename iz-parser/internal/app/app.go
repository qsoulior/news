package app

import (
	"net/http/cookiejar"
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

	cookiejar, err := cookiejar.New(nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init cookiejar")
	}

	client := httpclient.New(
		httpclient.URL(cfg.API.URL),
		httpclient.Headers(map[string]string{
			"Referer": cfg.API.URL,
		}),
		httpclient.CookieJar(cookiejar),
	)

	searchParser := service.NewNewsSearch(appID, client, &logger)
	archiveParser := service.NewNewsArchive(appID, client, &logger)
	feedParser := service.NewNewsFeed(appID, client, &logger)

	app.Run(
		&app.Config{
			ID:            appID,
			SearchParser:  searchParser,
			ArchiveParser: archiveParser,
			FeedParser:    feedParser,
			Logger:        &logger,
		},
		&app.Options{
			RabbitURL: cfg.RabbitMQ.URL,
			RedisURL:  cfg.Redis.URL,
		},
	)
}
