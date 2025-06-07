package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port string `env:"PORT" envDefault:"8090"`
}

// NewFromEnvironment reads Environment Variables and returns a pointer to a Config struct
func NewFromEnvironment() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
