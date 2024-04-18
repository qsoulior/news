package rabbitmq

import (
	"time"
)

type Config struct {
	URL          string
	AttemptCount int
	AttemptDelay time.Duration
}
