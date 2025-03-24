package config

import (
	"os"
)

type Config struct {
	USER_SERVICE      string
	DATABASE_URL      string
	MinioEndpoint     string
	MinioRootUser     string
	MinioRootPassword string
	S3Bucket          string
	RedisAddr         string
}

func NewConfig() *Config {
	return &Config{
		USER_SERVICE:      getEnv("USER_SERVICE", "localhost:8080"),
		DATABASE_URL:      getEnv("DATABASE_URL", "postgres://postgres:070823@postgresUser:5432/users?sslmode=disable"),
		MinioEndpoint:     getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioRootUser:     getEnv("MINIO_ROOT_USER", "minioadmin"),
		MinioRootPassword: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
		S3Bucket:          getEnv("S3_BUCKET", "my-bucket"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
