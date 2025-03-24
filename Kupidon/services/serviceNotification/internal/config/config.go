package config

import (
	"os"
)

type Config struct {
	KafkaBrokers   []string
	KafkaLikeTopic string
	TelegramToken  string
	GroupId        string
	UserURL        string
}

func NewConfig() *Config {
	return &Config{
		KafkaBrokers:   []string{getEnv("KAFKA_URL", "localhost:9092")},
		KafkaLikeTopic: getEnv("KAFKA_LIKE_TOPIC", "likes-topic"),
		TelegramToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
		GroupId:        getEnv("GROUP_ID", ""),
		UserURL:        getEnv("USER_SERVICE", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
