package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qsoulior/news/parser/internal/service"
	"github.com/rs/zerolog"
)

type workerConfig struct {
	Delay  time.Duration
	Logger *zerolog.Logger

	News service.News
	Page service.Page
}

type worker struct {
	workerConfig
}

func newWorker(cfg workerConfig) *worker {
	return &worker{cfg}
}

func (w *worker) Run(ctx context.Context) error {
	page, err := w.Page.Get(ctx)
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
				err = w.Page.Set(ctx, nextPage)
				if err != nil {
					w.Logger.Error().Str("next_page", nextPage).Err(err).Msg("")
				}

				delay = w.Delay
				w.Logger.Info().Str("page", page).Str("next_page", nextPage).Dur("delay", delay).Msg("parsed")
				page = nextPage
			} else {
				if delay > 0 {
					delay *= 2
				} else {
					delay = w.Delay
				}
				w.Logger.Error().Str("page", page).Err(err).Dur("delay", delay).Msg("")
			}

			timer.Reset(delay)
		}
	}
}
