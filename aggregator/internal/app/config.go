package app

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP     ConfigHTTP     `yaml:"http"`
		RabbitMQ ConfigRabbitMQ `yaml:"rabbitmq"`
		MongoDB  ConfigMongoDB  `yaml:"mongodb"`
	}

	ConfigHTTP struct {
		Host    string   `yaml:"host"`
		Port    string   `yaml:"port"`
		Origins []string `yaml:"origins"`
	}

	ConfigRabbitMQ struct {
		URL string `yaml:"url"`
	}

	ConfigMongoDB struct {
		URI string `yaml:"url"`
	}
)

func NewConfig(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return cfg, nil
}
