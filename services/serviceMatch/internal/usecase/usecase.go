package usecase

import (
	"context"
	"fmt"
	"log"
)

type MatchRepository interface {
	SaveLike(fromUserID, toUserID int64) error
	CheckMatch(fromUserID, toUserID int64) (bool, error)
}

type MatchKafka interface {
	SendMessage(ctx context.Context, message string) error
}

type Usecase struct {
	repo          MatchRepository
	kafkaProducer MatchKafka
}

func NewUseCase(repo MatchRepository, kafkaProducer MatchKafka) *Usecase {
	return &Usecase{repo: repo, kafkaProducer: kafkaProducer}
}

// Like - процесс лайкания и проверки совпадений
func (uc *Usecase) Like(ctx context.Context, fromUserID, toUserID int64) error {
	log.Printf("User %d liked user %d", fromUserID, toUserID)

	// Сохранение лайка в репозитории
	if err := uc.repo.SaveLike(fromUserID, toUserID); err != nil {
		return fmt.Errorf("failed to save like: %w", err)
	}

	// Создание события лайка и асинхронная отправка в Kafka
	likeEvent := fmt.Sprintf("like_%d_%d", fromUserID, toUserID)
	go func() {
		if err := uc.kafkaProducer.SendMessage(ctx, likeEvent); err != nil {
			log.Printf("Error sending like event to Kafka: %v", err)
		}
	}()
	return nil
}
