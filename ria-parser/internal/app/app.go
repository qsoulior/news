package app

import (
	"context"
	"fmt"
	"log"

	"github.com/go-rod/rod"
	"github.com/qsoulior/news/parser/app"
	"github.com/qsoulior/news/ria-parser/internal/service"
)

func Run(cfg *Config) {
	appID := "ria"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	browser := rod.New().Context(ctx)
	err := browser.Connect()
	if err != nil {
		log.Fatal(fmt.Errorf("browser.Connect: %w", err))
	}

	defer func() {
		err := browser.Close()
		if err != nil {
			log.Println(fmt.Errorf("browser.Close: %w", err))
		}
	}()

	searchParser := service.NewNewsSearch(cfg.API.SearchURL, appID, browser)
	feedParser := service.NewNewsFeed(cfg.API.FeedURL, appID, browser)

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
