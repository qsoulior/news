package app

import (
	"context"
	"log"

	"github.com/go-rod/rod"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/ria-parser/internal/service"
)

func Run(cfg *Config) {
	appID := "ria"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	browser := rod.New().Context(ctx)
	if err := browser.Connect(); err != nil {
		log.Fatalf("failed to init browser: %s", err)
	}

	defer func() {
		if err := browser.Close(); err != nil {
			log.Fatalf("failed to close browser: %s", err)
		}
	}()

	client := httpclient.New()

	searchParser := service.NewNewsSearch(appID, client, cfg.API.URL, browser)
	// archiveParser := service.NewNewsArchive(appID, client, cfg.API.URL, browser)
	feedParser := service.NewNewsFeed(appID, client, cfg.API.URL)

	app.Run(
		&app.Config{
			ID:       appID,
			RabbitMQ: app.ConfigRabbitMQ(cfg.RabbitMQ),
			Redis:    app.ConfigRedis(cfg.Redis),
		},
		searchParser,
		nil,
		feedParser,
	)
}
