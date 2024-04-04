package app

import (
	"github.com/qsoulior/news/lenta-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

func Run(cfg *Config) {
	appID := "lenta"

	client := httpclient.New()
	searchParser := service.NewNewsSearch(appID, cfg.API.SearchURL, client)
	archiveParser := service.NewNewsArchive(appID, cfg.API.FeedURL, client)

	app.Run(
		&app.Config{
			ID:       appID,
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		searchParser,
		archiveParser,
	)
}
