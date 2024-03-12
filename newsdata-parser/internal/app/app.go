package app

import (
	"github.com/qsoulior/news/newsdata-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	consumerParser := service.NewNews(service.NewsConfig{
		BaseAPI:   cfg.API.Consumer.URL,
		AccessKey: cfg.API.Consumer.AccessKey,
	})

	workerParser := service.NewNews(service.NewsConfig{
		BaseAPI:   cfg.API.Worker.URL,
		AccessKey: cfg.API.Worker.AccessKey,
	})

	app.Run(&app.Config{RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ), Redis: app.ConfigRedis(cfg.Redis)}, consumerParser, workerParser)
}
