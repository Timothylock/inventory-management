package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DbUrl  string `split_words:"true" required:"true"`
	DbUser string `split_words:"true" required:"true"`
	DbPass string `split_words:"true" required:"true"`
	DbName string `split_words:"true" required:"true"`

	UpcUrl   string `split_words:"true" required:"true"`
	UpcToken string `split_words:"true" required:"true"`

	EmailSMTPServ string `split_words:"true" required:"false"`
	EmailSMTPPort int    `split_words:"true" required:"false"`
	EmailUsername string `split_words:"true" required:"false"`
	EmailPassword string `split_words:"true" required:"false"`

	FrontendPath string `split_words:"true" required:"true"`
}

func FromEnvironment() (*Config, error) {
	cfg := new(Config)
	err := envconfig.Process("", cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
