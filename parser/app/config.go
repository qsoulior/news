package app

type (
	Config struct {
		ID       string
		RabbitMQ ConfigRabbitMQ
		Redis    ConfigRedis
	}

	ConfigRabbitMQ struct {
		URL string
	}

	ConfigRedis struct {
		URL string
	}
)
