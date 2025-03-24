package main

import (
	"context"
	"service1/internal/config"
	"service1/internal/handler"
	"service1/internal/repository"
	"service1/internal/storage"
	"service1/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitUserUse(cfg *config.Config) (*gin.Engine, error) {
	// Подключение к PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		logrus.Fatalf("Ошибка при подключении к базе данных: %v", err)
		return nil, err
	}

	// Настройка логирования
	logger := logrus.New()

	// Инициализация репозитория
	repo := repository.NewUserRepository(pool, logger)

	// Подключение к MinIO
	s3, err := storage.NewMinioStorage(cfg)
	if err != nil {
		return nil, err
	}

	// Подключение к Redis
	redis, err := storage.NewRedisStorage(cfg)
	if err != nil {
		logrus.Warn("Redis не подключен: ", err)
		redis = nil
	}

	// Создание Usecase
	uc := usecase.NewUserUsecase(repo, s3, redis)

	// Инициализация хендлеров
	_, router := handler.NewUserHandler(*uc)

	// Подключение Swagger документации
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router, nil
}

func main() {
	cfg := config.NewConfig()
	router, err := InitUserUse(cfg)
	if err != nil {
		logrus.Fatal("Ошибка инициализации сервера: ", err)
	}

	// Запуск сервера
	logrus.Infof("Запуск сервера на %s", cfg.USER_SERVICE)
	if err := router.Run(cfg.USER_SERVICE); err != nil {
		logrus.Fatal("Ошибка при запуске сервера: ", err)
	}
}
