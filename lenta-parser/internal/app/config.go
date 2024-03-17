package app

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		RabbitMQ ConfigRabbitMQ `yaml:"rabbitmq"`
		Redis    ConfigRedis    `yaml:"redis"`
		API      ConfigAPI      `yaml:"api"`
	}

	ConfigAPI struct {
		FeedURL   string `yaml:"feed_url"`
		SearchURL string `yaml:"search_url"`
	}

	ConfigRabbitMQ struct {
		URL string `yaml:"url"`
	}

	ConfigRedis struct {
		URL string `yaml:"url"`
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
