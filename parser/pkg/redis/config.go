package redis

import (
	"time"
)

type RedisConfig struct {
	URL          string
	AttemptCount int
	AttemptDelay time.Duration
}
