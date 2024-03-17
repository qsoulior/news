package app

import (
	"context"
	"time"

	"github.com/qsoulior/news/parser/internal/service"
	"github.com/rs/zerolog"
)

type releaserConfig struct {
	Delay  time.Duration
	Logger *zerolog.Logger
	News   service.News
}

type releaser struct {
	releaserConfig
}

func newReleaser(cfg releaserConfig) *releaser {
	return &releaser{cfg}
}

func (r *releaser) Run(ctx context.Context) {
	timer := time.NewTimer(r.Delay)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			err := r.News.Release(ctx)
			if err != nil {
				r.Logger.Error().Err(err).Msg("")
			}
			timer.Reset(r.Delay)
		}
	}
}
