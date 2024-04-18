package mongodb

import (
	"time"
)

type Config struct {
	URI          string
	AttemptCount int
	AttemptDelay time.Duration
}
