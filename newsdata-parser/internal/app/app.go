package app

import (
	"github.com/qsoulior/news/newsdata-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	appID := "newsdata"

	searchParser := service.NewNews(service.NewsConfig{
		AppID:     appID,
		BaseAPI:   cfg.API.Consumer.URL,
		AccessKey: cfg.API.Consumer.AccessKey,
	})

	feedParser := service.NewNews(service.NewsConfig{
		AppID:     appID,
		BaseAPI:   cfg.API.Worker.URL,
		AccessKey: cfg.API.Worker.AccessKey,
	})

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
