package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPServer HTTPServer `envconfig:"HTTP"`
	GRPCServer GRPCServer `envconfig:"GRPC"`
	PgConn     PgConn     `envconfig:"PG"`
	App        App        `envconfig:"APP"`
}

func New() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig.Process(): %w", err)
	}

	return cfg, nil
}
