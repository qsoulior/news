package app

import (
	"github.com/qsoulior/news/newsdata-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
)

func Run(cfg *Config) {
	appID := "newsdata"

	client := httpclient.New(
		httpclient.URL(cfg.API.URL),
	)

	searchParser := service.NewNews(appID, cfg.API.Search.AccessKey, client)
	archiveParser := service.NewNews(appID, cfg.API.Archive.AccessKey, client)

	app.Run(
		&app.Config{
			ID:            appID,
			SearchParser:  searchParser,
			ArchiveParser: archiveParser,
		},
		&app.Options{
			RabbitURL: cfg.RabbitMQ.URL,
			RedisURL:  cfg.Redis.URL,
		},
	)
}
