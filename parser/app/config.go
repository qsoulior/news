package app

type Config struct {
	API struct {
		URL       string
		AccessKey string
	}

	RabbitMQ struct {
		URL string
	}
}
