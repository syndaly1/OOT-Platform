package config

import "os"

type Config struct {
	GRPCPort    string
	DatabaseURL string
}

func LoadConfig() *Config {
	return &Config{
		GRPCPort:    getEnv("CONTENT_GRPC_PORT", "50052"),
		DatabaseURL: getEnv("CONTENT_DATABASE_URL", "postgres://postgres:postgres@localhost:5433/content_db?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
