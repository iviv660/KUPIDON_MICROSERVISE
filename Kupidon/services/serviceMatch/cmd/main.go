package main

import (
	"context"
	"fmt"
	"log"
	"service3/internal/config"
	"service3/internal/handler"
	"service3/internal/repository"
	"service3/internal/storage"
	"service3/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitMatchUse(cfg *config.Config) (*gin.Engine, error) {

	pool, err := initPostgresPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize database connection: %w", err)
	}

	// Инициализация репозитория
	repo := repository.NewRepository(pool)

	// Подключение к Kafka
	kfk, err := storage.New(cfg.KAFKA_URL, cfg.KAFKA_LIKE_TOPIC)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kafka producer: %w", err)
	}

	// Логика UseCase
	uc := usecase.NewUseCase(repo, kfk)

	// Gin router
	router := gin.Default()

	// Обработчики
	handler.NewMatchHandler(uc, router)

	return router, nil
}

func initPostgresPool(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	return pool, nil
}

func main() {
	cfg := config.NewConfig()

	// Инициализация и запуск приложения
	router, err := InitMatchUse(cfg)
	if err != nil {
		log.Fatalf("Error initializing match service: %v", err)
	}

	// Запуск сервера
	log.Println("Starting server on", cfg.SERVICE_MATCH)
	if err := router.Run(cfg.SERVICE_MATCH); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
