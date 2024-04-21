package app

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		RabbitMQ ConfigRabbitMQ `yaml:"rabbitmq"`
		Redis    ConfigRedis    `yaml:"redis"`
		Service  ConfigService  `yaml:"service"`
	}

	ConfigService struct {
		URL          string        `yaml:"url"`
		FeedDelay    time.Duration `yaml:"feed_delay"`
		ArchiveDelay time.Duration `yaml:"archive_delay"`
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
