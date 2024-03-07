package mongodb

import (
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	URI          string
	AttemptCount int
	AttemptDelay time.Duration
	Logger       *zerolog.Logger
}
