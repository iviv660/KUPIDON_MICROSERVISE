package storage

import (
	"context"
	"fmt"
	"service1/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type redisStorage struct {
	Client *redis.Client
}

func NewRedisStorage(cfg *config.Config) (RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Проверяем подключение
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %w", err)
	}

	return &redisStorage{Client: client}, nil
}

func (r *redisStorage) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.Client.Get(ctx, key)
}

func (r *redisStorage) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Client.Set(ctx, key, value, expiration)
}

func (r *redisStorage) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.Client.Del(ctx, keys...)
}
