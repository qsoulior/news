package app

import (
	"github.com/qsoulior/news/lenta-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

func Run(cfg *Config) {
	appID := "lenta"

	client := httpclient.New()
	searchParser := service.NewNewsSearch(appID, cfg.API.Search.URL, client)
	archiveParser := service.NewNewsArchive(appID, cfg.API.Archive.URL, client)
	feedParser := service.NewNewsFeed(appID, cfg.API.Feed.URL, client)

	app.Run(
		&app.Config{
			ID:       appID,
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		searchParser,
		archiveParser,
		feedParser,
	)
}
