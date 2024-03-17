package app

import (
	"github.com/qsoulior/news/lenta-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	consumerParser := service.NewNewsSearch(cfg.API.URL)
	workerParser := service.NewNewsFeed(cfg.API.URL)

	app.Run(
		&app.Config{
			ID:       "lenta",
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		consumerParser,
		workerParser,
	)
}
