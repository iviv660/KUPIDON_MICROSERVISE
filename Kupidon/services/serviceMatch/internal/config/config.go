// config.go
package config

import (
	"os"
	"strings"
)

type Config struct {
	DATABASE_URL     string
	KAFKA_URL        []string
	SERVICE_MATCH    string
	KAFKA_LIKE_TOPIC string
	KAFKA_GROUP_ID   string
}

func NewConfig() *Config {
	kafkaURL := getEnv("KAFKA_URL", "")
	brokers := strings.Split(kafkaURL, ",")

	return &Config{
		DATABASE_URL:     getEnv("DATABASE_URL", "postgres://postgres:070823@postgresMatch:5432/match?sslmode=disable"),
		SERVICE_MATCH:    getEnv("SERVICE_MATCH", ":8081"),
		KAFKA_URL:        brokers,
		KAFKA_LIKE_TOPIC: getEnv("KAFKA_LIKE_TOPIC", ""),
		KAFKA_GROUP_ID:   getEnv("KAFKA_GROUP_ID", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
