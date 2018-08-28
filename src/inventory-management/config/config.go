package config

import (
	"net"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DbUrl string `split_words:"true" required:"true"`
	DbUser string `split_words:"true" required:"true"`
	DbPass string `split_words:"true" required:"true"`
}

func FromEnvironment() (*Config, error) {
	cfg := new(Config)
	err := envconfig.Process("", cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
