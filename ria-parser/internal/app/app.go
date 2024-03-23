package app

import (
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/ria-parser/internal/service"
)

func Run(cfg *Config) {
	appID := "ria"
	searchParser := service.NewNewsSearch(cfg.API.SearchURL, appID)
	feedParser := service.NewNewsFeed(cfg.API.FeedURL, appID)

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
