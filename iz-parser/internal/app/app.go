package app

import (
	"github.com/qsoulior/news/iz-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	consumerParser := service.NewNewsFeed(service.NewsConfig{
		BaseAPI: cfg.API.URL,
	})

	workerParser := service.NewNewsFeed(service.NewsConfig{
		BaseAPI: cfg.API.URL,
	})

	app.Run(
		&app.Config{
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		consumerParser,
		workerParser,
	)
}
