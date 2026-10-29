package config

import (
	"os"
	"strings"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type Config struct {
	AppEnv             string
	CORSAllowedOrigins []string
	Port               string
	RedisURL           string
	TrustedProxies     []string
}

func Load() Config {
	return Config{
		AppEnv:             envOrDefault("APP_ENV", EnvironmentDevelopment),
		CORSAllowedOrigins: splitList(envOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		Port:               envOrDefault("PORT", "3000"),
		RedisURL:           required("REDIS_URL"),
		TrustedProxies:     splitList(envOrDefault("TRUSTED_PROXIES", "127.0.0.1,::1")),
	}
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(strings.TrimSpace(c.AppEnv), EnvironmentProduction)
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func required(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic(key + "_ENV_NOT_SET")
	}
	return value
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			result = append(result, item)
		}
	}

	return result
}
