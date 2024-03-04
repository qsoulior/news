package app

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		Host    string   `yaml:"host"`
		Port    string   `yaml:"port"`
		Origins []string `yaml:"origins"`
	} `yaml:"http"`

	RabbitMQ struct {
		URL string `yaml:"url"`
	} `yaml:"rabbitmq"`

	MongoDB struct {
		URL string `yaml:"url"`
	} `yaml:"mongodb"`
}

func NewConfig(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return cfg, nil
}
