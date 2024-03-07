package redis

import (
	"time"

	"github.com/rs/zerolog"
)

type RedisConfig struct {
	URL          string
	AttemptCount int
	AttemptDelay time.Duration
	Logger       *zerolog.Logger
}
