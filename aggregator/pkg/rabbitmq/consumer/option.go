package consumer

type Option func(*consumer)

func Ack(auto bool) Option {
	return func(c *consumer) {
		c.autoAck = auto
	}
}
