package app

import (
	"github.com/qsoulior/news/lenta-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	appID := "lenta"
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
