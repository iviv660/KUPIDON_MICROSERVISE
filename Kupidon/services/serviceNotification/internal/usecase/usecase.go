package usecase

import (
	"serviceNotification/internal/entity"
)

// TelegramBotSender интерфейс для отправки сообщений в Telegram
type TelegramBotSender interface {
	SendMessage(entity.Message) error
}

// BotUsecase - бизнес-логика для обработки сообщений
type BotUsecase struct {
	sender TelegramBotSender
}

// NewBotUsecase создает новый экземпляр BotUsecase
func NewBotUsecase(sender TelegramBotSender) *BotUsecase {
	return &BotUsecase{
		sender: sender,
	}
}

// HandleEvent обрабатывает событие и отправляет сообщение в Telegram
func (u *BotUsecase) SendMessage(msg entity.Message) error {
	u.sender.SendMessage(msg)
	return nil
}
