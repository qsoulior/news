package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/qsoulior/news/newsdata-parser/internal/service"
	"github.com/rs/zerolog"
)

type WorkerConfig struct {
	Delay    time.Duration
	Logger   *zerolog.Logger
	Services struct {
		News service.News
		Page service.Page
	}
}

type Worker struct {
	WorkerConfig
}

func NewWorker(cfg WorkerConfig) *Worker {
	return &Worker{cfg}
}

func (w *Worker) Run() error {
	var delay time.Duration = 0
	page, err := w.Services.Page.Get()
	if !errors.Is(err, service.ErrNotExist) {
		return fmt.Errorf("w.Services.Page.Get: %w", err)
	}
	w.Logger.Info().Str("page", page).Msg("init page")

	for {
		time.Sleep(delay)

		nextPage, err := w.Services.News.ParsePage(page)
		if err != nil {
			w.Logger.Error().Err(err).Msg("")
			delay *= 2
		}

		err = w.Services.Page.Set(nextPage)
		if err != nil {
			w.Logger.Error().Err(err).Msg("")
		}

		delay = w.Delay
	}
}
