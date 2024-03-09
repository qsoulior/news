package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qsoulior/news/parser/internal/service"
	"github.com/rs/zerolog"
)

type WorkerConfig struct {
	Delay  time.Duration
	Logger *zerolog.Logger

	News service.News
	Page service.Page
}

type worker struct {
	WorkerConfig
}

func NewWorker(cfg WorkerConfig) *worker {
	return &worker{cfg}
}

func (w *worker) Run(ctx context.Context) error {
	page, err := w.Page.Get()
	if !errors.Is(err, service.ErrNotExist) {
		return fmt.Errorf("w.Services.Page.Get: %w", err)
	}

	w.Logger.Info().Str("page", page).Msg("init page")
	w.work(ctx, page)
	return nil
}

func (w *worker) work(ctx context.Context, page string) {
	var delay time.Duration = 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(delay)
			nextPage, err := w.News.Parse("", page)
			if err != nil {
				w.Logger.Error().Err(err).Msg("")
				delay *= 2
				continue
			}

			err = w.Page.Set(nextPage)
			if err != nil {
				w.Logger.Error().Err(err).Msg("")
			}

			page = nextPage
			delay = w.Delay
		}
	}
}
