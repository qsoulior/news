package app

import (
	"github.com/qsoulior/news/iz-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run(cfg *Config) {
	consumerParser := service.NewNews()
	workerParser := service.NewNews()

	app.Run(&app.Config{RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ), Redis: app.ConfigRedis(cfg.Redis)}, consumerParser, workerParser)
}
