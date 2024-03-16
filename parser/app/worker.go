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
	if err != nil && !errors.Is(err, service.ErrNotExist) {
		return fmt.Errorf("w.Services.Page.Get: %w", err)
	}

	w.Logger.Info().Str("page", page).Msg("init page")
	w.work(ctx, page)
	return nil
}

func (w *worker) work(ctx context.Context, page string) {
	var delay time.Duration = 0
	timer := time.NewTimer(delay)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			nextPage, err := w.News.Parse(ctx, "", page)
			if err == nil {
				err = w.Page.Set(nextPage)
				if err != nil {
					w.Logger.Error().Str("next_page", nextPage).Err(err).Msg("")
				}

				w.Logger.Info().Str("page", page).Str("next_page", nextPage).Msg("parsed")
				page = nextPage
				delay = w.Delay
			} else {
				w.Logger.Error().Str("page", page).Err(err).Msg("")
				if delay >= 0 {
					delay *= 2
				} else {
					delay = w.Delay
				}
			}

			timer.Reset(delay)
		}
	}
}
