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
	EnableHTTPS        bool   `env:"ENABLE_HTTPS"`
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

	if flag.Lookup("a") == nil {
		flag.StringVar(&config.RunAddr, "a", "localhost:8181", "address and port to run server")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&config.DatabaseConnection, "d", "", "database connection string")
	}
	if flag.Lookup("e") == nil {
		flag.StringVar(&config.MiniIOConnection.Endpoint, "e", "localhost:9000", "mini-io address and port to run")
	}
	if flag.Lookup("k") == nil {
		flag.StringVar(&config.MiniIOConnection.AccessKeyID, "k", "minioadmin", "mini-io access key")
	}
	if flag.Lookup("m") == nil {
		flag.StringVar(&config.MiniIOConnection.SecretAccessKey, "m", "minioadmin", "mini-io secret key")
	}
	if flag.Lookup("u") == nil {
		flag.BoolVar(&config.MiniIOConnection.UseSSL, "u", false, "set to true if using HTTPS")
	}
	if flag.Lookup("b") == nil {
		flag.StringVar(&config.MiniIOConnection.BucketName, "b", "keeper", "miniio bucket name")
	}
	if flag.Lookup("s") == nil {
		flag.BoolVar(&config.EnableHTTPS, "s", false, "enable HTTPS connection")
	}

	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
