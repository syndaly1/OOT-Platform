package config

import (
	"os"
	"time"
)

type Config struct {
	DB    DBConfig
	Redis RedisConfig
	NATS  NATSConfig
	JWT   JWTConfig
	SMTP  SMTPConfig
	GRPC  GRPCConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host string
	Port string
}

type NATSConfig struct {
	URL string
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type GRPCConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("AUTH_DB_HOST", "localhost"),
			Port:     getEnv("AUTH_DB_PORT", "5432"),
			User:     getEnv("AUTH_DB_USER", "postgres"),
			Password: getEnv("AUTH_DB_PASSWORD", "postgres"),
			Name:     getEnv("AUTH_DB_NAME", "auth_db"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "secret"),
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 7 * 24 * time.Hour,
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
		GRPC: GRPCConfig{
			Port: getEnv("AUTH_GRPC_PORT", "50051"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
