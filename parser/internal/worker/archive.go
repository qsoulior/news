package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qsoulior/news/parser/internal/service"
	"github.com/rs/zerolog"
)

type archive struct {
	*worker
	news service.News
	page service.Page
}

func NewArchive(delay time.Duration, logger *zerolog.Logger, news service.News, page service.Page) *archive {
	worker := &worker{
		delay:  delay,
		logger: logger,
	}

	return &archive{worker: worker, news: news, page: page}
}

func (a *archive) Run(ctx context.Context) error {
	page, err := a.page.Get(ctx)
	if err != nil && !errors.Is(err, service.ErrNotExist) {
		return fmt.Errorf("w.Services.Page.Get: %w", err)
	}

	a.logger.Info().Str("page", page).Msg("init page")
	a.work(ctx, page)
	return nil
}

func (a *archive) work(ctx context.Context, page string) {
	var delay time.Duration = 0
	timer := time.NewTimer(delay)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			nextPage, err := a.news.Parse(ctx, "", page)
			if err == nil {
				err = a.page.Set(ctx, nextPage)
				if err != nil {
					a.logger.Error().Str("next_page", nextPage).Err(err).Send()
				}

				delay = a.delay
				a.logger.Info().Str("page", page).Str("next_page", nextPage).Dur("delay", delay).Msg("parsed")
				page = nextPage
			} else {
				if delay > 0 {
					delay *= 2
				} else {
					delay = a.delay
				}
				a.logger.Error().Str("page", page).Err(err).Dur("delay", delay).Send()
			}

			timer.Reset(delay)
		}
	}
}
