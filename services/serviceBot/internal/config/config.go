package config

import "os"

type Config struct {
	TELEGRAM_BOT_TOKEN string
	USER_SERVICE       string
	MATCH_SERVICE      string
}

func NewConfig() *Config {
	return &Config{
		TELEGRAM_BOT_TOKEN: getEnv("TELEGRAM_BOT_TOKEN", ""),
		USER_SERVICE:       getEnv("USER_SERVICE", ""),
		MATCH_SERVICE:      getEnv("MATCH_SERVICE", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
