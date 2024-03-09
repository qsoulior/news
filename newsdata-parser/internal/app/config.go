package app

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/qsoulior/news/parser/app"
)

type (
	Config struct {
		*app.Config

		API ConfigAPI `yaml:"api"`
	}

	ConfigAPI struct {
		URL       string `yaml:"url"`
		AccessKey string `yaml:"access_key"`
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
