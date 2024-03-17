package main

import (
	"flag"
	"log"

	"github.com/qsoulior/news/lenta-parser/internal/app"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "config file path")
	flag.Parse()

	if path == "" {
		flag.PrintDefaults()
		return
	}

	cfg, err := app.NewConfig(path)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	app.Run(cfg)
}
