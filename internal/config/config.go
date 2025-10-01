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
	MiniIOConnection   MiniIOConnection
}

type MiniIOConnection struct {
	Endpoint        string `env:"ENDPOINT"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `env:"USE_SSL"`
	BucketName      string `env:"BUCKET_NAME"`
}

func NewConfig() (*Config, error) {
	var config Config

	config.SecretKey = "supersecretkey"
	config.TokenExp = time.Hour * 72

	flag.StringVar(&config.RunAddr, "a", "localhost:8181", "address and port to run server")
	flag.StringVar(&config.DatabaseConnection, "d", "", "Database connection string")
	flag.StringVar(&config.MiniIOConnection.Endpoint, "e", "localhost:9000", "mini-io address and port to run")
	flag.StringVar(&config.MiniIOConnection.AccessKeyID, "k", "minioadmin", "mini-io access key")
	flag.StringVar(&config.MiniIOConnection.SecretAccessKey, "s", "minioadmin", "mini-io secret key")
	flag.BoolVar(&config.MiniIOConnection.UseSSL, "u", false, "set to true if using HTTPS")
	flag.StringVar(&config.MiniIOConnection.BucketName, "b", "keeper", "miniio bucket name")

	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
