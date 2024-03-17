package app

import (
	"github.com/qsoulior/news/iz-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	searchParser := service.NewNewsSearch(cfg.API.URL)
	feedParser := service.NewNewsFeed(cfg.API.URL)

	app.Run(
		&app.Config{
			ID:       "iz",
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		searchParser,
		feedParser,
	)
}
