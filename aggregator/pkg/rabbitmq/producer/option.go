package producer

import "time"

type Option func(*producer)

func Timeout(d time.Duration) Option {
	return func(c *producer) {
		c.timeout = d
	}
}
