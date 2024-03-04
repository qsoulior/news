package app

type Config struct {
	HTTP struct {
		Host    string
		Port    string
		Origins []string
	}

	RabbitMQ struct {
		URL string
	}
}
