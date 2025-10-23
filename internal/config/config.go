package config

import (
	"os"
	"time"
)

type Config struct {
	AWS struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		Bucket          string
		Endpoint        string
	}
	Server struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	App struct {
		BaseDirectory string
	}
}

func Load() *Config {
	cfg := &Config{}

	// AWS Configuration
	cfg.AWS.Region = getEnv("AWS_REGION", "us-east-1")
	cfg.AWS.AccessKeyID = getEnv("AWS_ACCESS_KEY_ID", "")
	cfg.AWS.SecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "")
	cfg.AWS.Bucket = getEnv("AWS_BUCKET", "")
	cfg.AWS.Endpoint = getEnv("AWS_ENDPOINT", "")

	// Server Configuration
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.ReadTimeout = 15 * time.Second
	cfg.Server.WriteTimeout = 15 * time.Second

	// App Configuration
	cfg.App.BaseDirectory = getEnv("BASE_DIRECTORY", "")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
