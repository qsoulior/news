package app

import (
	"github.com/qsoulior/news/iz-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

func Run(cfg *Config) {
	appID := "iz"

	client := httpclient.New(
		httpclient.URL(cfg.API.URL),
		httpclient.Headers(map[string]string{
			"Referer": cfg.API.URL,
		}),
	)

	searchParser := service.NewNewsSearch(appID, client)
	archiveParser := service.NewNewsArchive(appID, client)

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
