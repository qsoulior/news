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
		URL    string `yaml:"url"`
		Search struct {
			AccessKey string `yaml:"access_key"`
		} `yaml:"search"`
		Archive struct {
			AccessKey string `yaml:"access_key"`
		} `yaml:"archive"`
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
