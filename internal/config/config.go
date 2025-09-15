package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddr            string `env:"RUN_ADDRESS"`
	DatabaseConnection string `env:"DATABASE_URI"`
	SecretKey          string `env:"SECRET_KEY"`
	TokenExp           time.Duration
}

func NewConfig() (*Config, error) {
	var config Config

	config.SecretKey = "supersecretkey"
	config.TokenExp = time.Hour * 72

	flag.StringVar(&config.RunAddr, "a", "localhost:8181", "address and port to run server")
	flag.StringVar(&config.DatabaseConnection, "d", "", "Database connection string")

	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
