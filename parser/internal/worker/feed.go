package worker

import (
	"context"
	"time"

	"github.com/qsoulior/news/parser/internal/service"
	"github.com/rs/zerolog"
)

type feed struct {
	*worker
	news service.News
}

func NewFeed(delay time.Duration, logger *zerolog.Logger, news service.News) *feed {
	worker := &worker{
		delay:  delay,
		logger: logger,
	}
	return &feed{worker: worker, news: news}
}

func (f *feed) Run(ctx context.Context) error {
	var delay time.Duration = 0
	timer := time.NewTimer(delay)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
			_, err := f.news.Parse(ctx, "", "")
			if err == nil {
				delay = f.delay
				f.logger.Info().Dur("delay", delay).Msg("parsed")
			} else {
				if delay > 0 {
					delay *= 2
				} else {
					delay = f.delay
				}
				f.logger.Error().Err(err).Dur("delay", delay).Send()
			}

			timer.Reset(delay)
		}
	}
}
