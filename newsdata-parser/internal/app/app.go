package app

import (
	"github.com/qsoulior/news/newsdata-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	parser := service.NewNews(service.NewsConfig{
		BaseAPI:   cfg.API.URL,
		AccessKey: cfg.API.AccessKey,
	})

	app.Run(&app.Config{RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ), Redis: app.ConfigRedis(cfg.Redis)}, parser)
}
