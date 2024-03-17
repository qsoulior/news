package app

import (
	"github.com/qsoulior/news/iz-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	appID := "iz"
	searchParser := service.NewNewsSearch(cfg.API.URL, appID)
	feedParser := service.NewNewsFeed(cfg.API.URL, appID)

	app.Run(
		&app.Config{
			ID:       appID,
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		searchParser,
		feedParser,
	)
}
