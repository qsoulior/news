package app

type (
	Config struct {
		API struct {
			URL       string
			AccessKey string
		}

		RabbitMQ ConfigRabbitMQ

		Redis struct {
			URL string
		}
	}

	ConfigRabbitMQ struct {
		URL string
	}
)
