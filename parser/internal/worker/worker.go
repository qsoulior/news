package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Worker interface {
	Run(ctx context.Context) error
}

type worker struct {
	delay  time.Duration
	logger *zerolog.Logger
}
