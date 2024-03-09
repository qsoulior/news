package app

type (
	Config struct {
		RabbitMQ ConfigRabbitMQ `yaml:"rabbitmq"`
		Redis    ConfigRedis    `yaml:"redis"`
	}

	ConfigRabbitMQ struct {
		URL string `yaml:"url"`
	}

	ConfigRedis struct {
		URL string `yaml:"url"`
	}
)
