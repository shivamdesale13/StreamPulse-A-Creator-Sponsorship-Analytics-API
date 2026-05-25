package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	JWTSecret      string
	WorkerPoolSize int
	QueueBuffer    int
}

func Load() *Config {
	workers, _ := strconv.Atoi(getEnv("WORKER_POOL_SIZE", "4"))
	buffer, _ := strconv.Atoi(getEnv("QUEUE_BUFFER", "100"))

	return &Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "streampulse-dev-secret-change-in-prod"),
		WorkerPoolSize: workers,
		QueueBuffer:    buffer,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
