package main

import (
	"context"
	"log"
	"serviceNotification/internal/adapter"
	clientsUser "serviceNotification/internal/client"
	"serviceNotification/internal/config"
	"serviceNotification/internal/delivery"
	"serviceNotification/internal/usecase"
)

func main() {
	cfg := config.NewConfig()
	userClient := clientsUser.NewHTTPUserServiseClient(cfg.UserURL)
	sender, err := adapter.NewTelegramBot(cfg.TelegramToken, userClient)
	if err != nil {
		log.Fatal(err)
	}
	bot := adapter.NewBothandle(*sender)
	go bot.BotStart()
	uc := usecase.NewBotUsecase(sender)
	kfk, err := delivery.NewKafkaConsumer(cfg.KafkaBrokers, cfg.KafkaLikeTopic, cfg.GroupId, uc)
	if err != nil {
		log.Fatal(err)
	}
	kfk.Start(context.Background())
}
