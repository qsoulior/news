package main

import (
	"flag"

	"github.com/qsoulior/news/aggregator/internal/app"
	"github.com/rs/zerolog"
)

func main() {
	out := zerolog.NewConsoleWriter()
	logger := zerolog.New(out).With().Timestamp().Logger()

	var path string
	flag.StringVar(&path, "c", "", "config file path")
	flag.Parse()

	if path == "" {
		flag.PrintDefaults()
		return
	}

	cfg, err := app.NewConfig(path)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	app.Run(cfg, &logger)
}
