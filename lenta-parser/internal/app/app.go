package app

import (
	"github.com/qsoulior/news/lenta-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	searchParser := service.NewNewsSearch(cfg.API.SearchURL)
	feedParser := service.NewNewsFeed(cfg.API.FeedURL)

	app.Run(
		&app.Config{
			ID:       "lenta",
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		searchParser,
		feedParser,
	)
}
