package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"service1/internal/config"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

type FileStorage interface {
	UploadFile(ctx context.Context, file multipart.File, fileName string, fileSize int64) (string, error)
}

func NewMinioStorage(cfg *config.Config) (*MinioStorage, error) {
	// Проверка обязательных настроек
	if cfg.MinioRootUser == "" || cfg.MinioRootPassword == "" || cfg.MinioEndpoint == "" || cfg.S3Bucket == "" {
		return nil, fmt.Errorf("configuration error: missing required MinIO settings")
	}

	// Подключение к MinIO
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioRootUser, cfg.MinioRootPassword, ""),
		Secure: false, // Измените на true, если используете HTTPS
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	// Проверка существования корзины
	exists, err := client.BucketExists(context.Background(), cfg.S3Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	// Если корзина не существует, создаем ее
	if !exists {
		log.Printf("Bucket %s does not exist. Creating...", cfg.S3Bucket)
		err = client.MakeBucket(context.Background(), cfg.S3Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully", cfg.S3Bucket)
	}

	// Устанавливаем публичную политику для корзины
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::` + cfg.S3Bucket + `/*"
			}
		]
	}`

	err = client.SetBucketPolicy(context.Background(), cfg.S3Bucket, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to set public policy on bucket: %w", err)
	}

	log.Printf("Public read access granted for bucket %s", cfg.S3Bucket)

	return &MinioStorage{
		Client: client,
		Bucket: cfg.S3Bucket,
	}, nil
}

func (s *MinioStorage) UploadFile(ctx context.Context, file multipart.File, fileName string, fileSize int64) (string, error) {
	// Задаем тайм-аут для контекста на время загрузки
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Загружаем файл в MinIO
	_, err := s.Client.PutObject(ctx, s.Bucket, fileName, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("error uploading to MinIO: %w", err)
	}

	// Формируем URL для доступа к файлу
	url := fmt.Sprintf("%s/%s/%s", "http://minio:9000", s.Bucket, fileName)

	return url, nil
}
