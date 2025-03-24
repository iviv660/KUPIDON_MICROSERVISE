package adapter

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"serviceNotification/internal/entity"
	"serviceNotification/internal/utilites"

	"gopkg.in/telebot.v4"
)

type UserClient interface {
	GetUserByID(userID int64) (*entity.User, error)
}

type TelegramBot struct {
	b           *telebot.Bot
	uc          UserClient
	likeTracker map[int64]int64 //Хранит, кто кого лайкнул (userID -> likerID)
}

type Bothandle struct {
	t TelegramBot
}

func NewBothandle(t TelegramBot) *Bothandle {
	return &Bothandle{t: t}
}

var likes int

func NewTelegramBot(token string, userClient UserClient) (*TelegramBot, error) {
	botAPI, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10},
	})
	if err != nil {
		log.Printf("Error creating Telegram bot: %v", err)
		return nil, err
	}
	log.Printf("Authorized on account %s", botAPI.Me.FirstName)

	bot := &TelegramBot{
		b:           botAPI,
		uc:          userClient,
		likeTracker: make(map[int64]int64),
	}

	return bot, nil
}

func (bot *TelegramBot) SendMessage(msg entity.Message) error {
	if msg.FromUserID == 0 || msg.ToUserID == 0 || msg.Text == "" {
		return errors.New("Не верные данные")
	}
	ToRecipient := &telebot.User{ID: msg.ToUserID}

	bot.likeTracker[msg.ToUserID] = msg.FromUserID

	profileKeys := &telebot.ReplyMarkup{ResizeKeyboard: true}
	btnShowProfile := profileKeys.Text("Показать анкету")
	profileKeys.Reply(profileKeys.Row(btnShowProfile))

	switch msg.Text {
	case "like":
		_, err := bot.b.Send(ToRecipient, "Вас лайкнули! Нажмите смотреть анкету?", profileKeys)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			return err

		}
	}

	return nil
}

func (bot *Bothandle) BotStart() {
	bot.t.b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		if likes == 1 {
			if ctx.Text() == "❤" {
				chat, err := bot.t.b.ChatByID(bot.t.likeTracker[ctx.Sender().ID])
				if err != nil {
					return ctx.Send("Произошла какая-то ошибка")
				}
				msg := fmt.Sprintf("Взаимный лайк начинай общаться: [Начать общение!] (tg://user?id=%o)", ctx.Sender().ID)
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				bot.t.b.Send(chat, msg, &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})

			}
			if ctx.Text() == "👎" {
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}
		}
		if ctx.Text() == "Показать анкету" {
			userID := ctx.Sender().ID
			likerID, exists := bot.t.likeTracker[userID]
			if !exists {
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				ctx.Send("Нет анкеты для показа")
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}

			user, err := bot.t.uc.GetUserByID(likerID)
			if err != nil {
				ctx.Send("Ошибка загрузки анкеты.")
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				ctx.Send("Нет анкеты для показа")
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}

			imageBytes, err := utilites.DownloadImageAsBytes(user.Photo)
			if err != nil {
				log.Println(err)
				ctx.Send("Ошибка загрузки фотографии из бд")
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				ctx.Send("Нет анкеты для показа")
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}

			Answer := &telebot.Photo{
				File:    telebot.FromReader(bytes.NewReader(imageBytes)),
				Caption: fmt.Sprintf("%s, %d, %s - %s", user.Name, user.Age, user.City, user.Description),
			}

			likes = 1
			key := [][]telebot.ReplyButton{
				{{Text: "❤"}, {Text: "👎"}},
			}
			return ctx.Send(Answer, &telebot.ReplyMarkup{ReplyKeyboard: key, ResizeKeyboard: true})
		}
		return nil
	})

	log.Println("Бот запущен...")
	bot.t.b.Start()
}
