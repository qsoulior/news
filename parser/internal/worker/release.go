package worker

import (
	"context"
	"time"

	"github.com/qsoulior/news/parser/internal/service"
	"github.com/rs/zerolog"
)

type release struct {
	*worker
	news service.News
}

func NewRelease(delay time.Duration, logger *zerolog.Logger, news service.News) *release {
	worker := &worker{
		delay:  delay,
		logger: logger,
	}
	return &release{worker: worker, news: news}
}

func (r *release) Run(ctx context.Context) error {
	timer := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
			count, err := r.news.Release(ctx)
			if err != nil {
				r.logger.Error().Err(err).Int("count", count).Msg("")
			}
			r.logger.Info().Dur("delay", r.delay).Int("count", count).Msg("released")
			timer.Reset(r.delay)
		}
	}
}
