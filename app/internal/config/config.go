package config

import (
	"os"
	"strconv"
)

const (
	AppVersion = "1.0.0"
	AppName    = "platform-api"
)

type Config struct {
	HTTPPort     string
	MetricsPort  string
	AppName      string
	Environment  string
	LogLevel     string
	PodName      string
	PodNamespace string
	NodeName     string
	EnableCORS   bool
}

func Load() *Config {
	return &Config{
		HTTPPort:     envString("PORT", "8080"),
		MetricsPort:  envString("METRICS_PORT", "9090"),
		AppName:      envString("APP_NAME", AppName),
		Environment:  envString("ENVIRONMENT", "development"),
		LogLevel:     envString("LOG_LEVEL", "info"),
		PodName:      envString("POD_NAME", "local"),
		PodNamespace: envString("POD_NAMESPACE", "default"),
		NodeName:     envString("NODE_NAME", "local"),
		EnableCORS:   envBool("ENABLE_CORS", true),
	}
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}
