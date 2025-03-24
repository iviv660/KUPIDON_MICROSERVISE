package delivery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"serviceNotification/internal/entity"
	"serviceNotification/internal/usecase"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	usecase *usecase.BotUsecase
	Reader  *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupId string, usecase *usecase.BotUsecase) (*KafkaConsumer, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" || groupId == "" {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupId,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		usecase: usecase,
		Reader:  reader,
	}, nil
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down...")
			return
		default:
			msg, err := c.Reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				continue
			}

			log.Printf("Received message: %s", string(msg.Value))

			// Парсим сообщение
			err = c.processMessage(string(msg.Value))
			if err != nil {
				log.Printf("Error processing message: %v", err)
				continue // Если ошибка обработки, не подтверждаем сообщение
			}

			// Подтверждаем, что сообщение обработано
			err = c.Reader.CommitMessages(ctx, msg)
			if err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

func (c *KafkaConsumer) processMessage(message string) error {
	if len(message) == 0 {
		return fmt.Errorf("error: received empty message")
	}

	// Обрабатываем лайк-сообщение
	if len(message) >= 5 && message[:5] == "like_" {
		var fromUserID, toUserID int64
		n, err := fmt.Sscanf(message, "like_%d_%d", &fromUserID, &toUserID)
		if err != nil || n != 2 {
			return fmt.Errorf("error parsing like message: %v", err)
		}

		text := "like"
		log.Printf("Processing like: from %d to %d", fromUserID, toUserID)

		err = c.usecase.SendMessage(entity.Message{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Text:       text,
		})
		if err != nil {
			return fmt.Errorf("failed to send like message: %v", err)
		}

		return nil
	}

	return fmt.Errorf("unrecognized message type: %s", message)
}
