package usecase

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"serviceBot/internal/entity"
	"serviceBot/utilites"
	"strconv"
	"time"

	"gopkg.in/telebot.v4"
)

type UserService interface {
	CreateUser(name, city, gender, description string, age int, telegramID int64, file []byte, filename string) error
	SearchUser(MinAge, MaxAge int, City, Gender string) ([]entity.User, error)
	Delete(id int64) error
	GetUserByID(userID int64) (*entity.User, error)
}

type MatchService interface {
	LikeUser(fromUserID, toUserID int64) error
}

type UseCase struct {
	userService  UserService
	matchService MatchService
}

func NewUseCase(userService UserService, matchService MatchService) *UseCase {
	return &UseCase{userService: userService, matchService: matchService}
}

var users = make(map[int64]entity.User)
var usersLike = make([]entity.User, 0)
var outID int64

func (uc *UseCase) StartBot(token string) {
	if token == "" {
		log.Fatal("Token empty")
	}
	var state int
	var autho int
	var likes int

	b, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10},
	})

	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	b.Handle("/start", func(ctx telebot.Context) error {
		user, err := uc.userService.GetUserByID(ctx.Sender().ID)
		if err != nil || user == nil {
			state = 1
			ctx.Send("Привет, это бот для знакомств!")
			time.Sleep(1 * time.Second)
			ctx.Send("Давай создадим твою анкету!")
			time.Sleep(1 * time.Second)
			state = 1
			return ctx.Send("Как тебя зовут?", &telebot.ReplyMarkup{RemoveKeyboard: true})
		}

		ctx.Send("Вот так выглядит твоя анкета:", &telebot.ReplyMarkup{RemoveKeyboard: true})
		fmt.Println(user)
		imageBytes, err := utilites.DownloadImageAsBytes(user.Photo)
		if err != nil {
			log.Println(err)
			return ctx.Send("Ошибка загрузки фотографии из бд")
		}
		Answer := &telebot.Photo{
			File:    telebot.FromReader(bytes.NewReader(imageBytes)),
			Caption: fmt.Sprintf("%s, %d, %s - %s", user.Name, user.Age, user.City, user.Description),
		}
		autho = 1
		ctx.Send(Answer)

		profileKeys := [][]telebot.ReplyButton{
			{{Text: "1"}, {Text: "2"}, {Text: "3"}},
		}
		return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
	})

	///////////////////////////////////////////////////////////////////////////////////////////////////////
	b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		if ctx.Text() == "Cмотреть" {
			autho = 1
			likes = 0
			return nil
		}
		if ctx.Text() == "Показать анкету" {
			autho = 0
			likes = 0
			return nil
		}
		if ctx.Text() == "❤" || ctx.Text() == "👎" {
			autho = 1
			return nil
		}
		if likes == 1 {
			if len(usersLike) <= 0 {
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				likes = 0
				autho = 1
				ctx.Send("Анкеты закончились :(")
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}
			outUser := usersLike[len(usersLike)-1]
			outID = outUser.TelegramID
			usersLike = usersLike[:len(usersLike)-1]

			image, err := utilites.DownloadImageAsBytes(outUser.Photo)
			if err != nil {
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				likes = 0
				autho = 1
				ctx.Send("произошла ошибка в боте:(")
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}

			Answer := &telebot.Photo{
				File:    telebot.FromReader(bytes.NewReader(image)),
				Caption: fmt.Sprintf("%s, %d, %s - %s", outUser.Name, outUser.Age, outUser.City, outUser.Description),
			}
			likes = 2
			key := [][]telebot.ReplyButton{
				{{Text: "❤"}, {Text: "👎"}, {Text: "💤"}},
			}
			return ctx.Send(Answer, &telebot.ReplyMarkup{ReplyKeyboard: key, ResizeKeyboard: true})
		}
		if likes == 2 {

			if ctx.Text() == "❤" {
				err := uc.matchService.LikeUser(ctx.Sender().ID, outID)
				if err != nil {
					log.Println("fdsfsdfsd", err)
					profileKeys := [][]telebot.ReplyButton{
						{{Text: "1"}, {Text: "2"}, {Text: "3"}},
					}
					likes = 0
					autho = 1
					ctx.Send("произошла ошибка в боте:(")
					return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
				}
				if len(usersLike) <= 0 {
					profileKeys := [][]telebot.ReplyButton{
						{{Text: "1"}, {Text: "2"}, {Text: "3"}},
					}
					likes = 0
					autho = 1
					ctx.Send("Анкеты закончились :(")
					return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
				}
				outUser := usersLike[len(usersLike)-1]
				outID = outUser.TelegramID
				usersLike = usersLike[:len(usersLike)-1]

				image, err := utilites.DownloadImageAsBytes(outUser.Photo)
				if err != nil {
					profileKeys := [][]telebot.ReplyButton{
						{{Text: "1"}, {Text: "2"}, {Text: "3"}},
					}
					likes = 0
					autho = 1
					ctx.Send("произошла ошибка в боте:(")
					return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
				}

				Answer := &telebot.Photo{
					File:    telebot.FromReader(bytes.NewReader(image)),
					Caption: fmt.Sprintf("%s, %d, %s - %s", outUser.Name, outUser.Age, outUser.City, outUser.Description),
				}
				likes = 2
				key := [][]telebot.ReplyButton{
					{{Text: "❤"}, {Text: "👎"}, {Text: "💤"}},
				}
				return ctx.Send(Answer, &telebot.ReplyMarkup{ReplyKeyboard: key, ResizeKeyboard: true})
			}

			if ctx.Text() == "👎" {
				if len(usersLike) <= 0 {
					profileKeys := [][]telebot.ReplyButton{
						{{Text: "1"}, {Text: "2"}, {Text: "3"}},
					}
					likes = 0
					autho = 1
					ctx.Send("Анкеты закончились :(")
					return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
				}
				outUser := usersLike[len(usersLike)-1]
				outID = outUser.TelegramID
				usersLike = usersLike[:len(usersLike)-1]

				image, err := utilites.DownloadImageAsBytes(outUser.Photo)
				if err != nil {
					profileKeys := [][]telebot.ReplyButton{
						{{Text: "1"}, {Text: "2"}, {Text: "3"}},
					}
					likes = 0
					autho = 1
					ctx.Send("произошла ошибка в боте:(")
					return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
				}

				Answer := &telebot.Photo{
					File:    telebot.FromReader(bytes.NewReader(image)),
					Caption: fmt.Sprintf("%s, %d, %s - %s", outUser.Name, outUser.Age, outUser.City, outUser.Description),
				}
				likes = 2
				key := [][]telebot.ReplyButton{
					{{Text: "❤"}, {Text: "👎"}, {Text: "💤"}},
				}
				return ctx.Send(Answer, &telebot.ReplyMarkup{ReplyKeyboard: key, ResizeKeyboard: true})
			}

			if ctx.Text() == "💤" {
				log.Println("dasdasd")
				likes = 0
				autho = 1
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}
		}

		if autho > 0 {
			user, err := uc.userService.GetUserByID(ctx.Sender().ID)
			autho, err = strconv.Atoi(ctx.Text())
			if err != nil {
				return ctx.Send("Нет такого варианта ответа")
			}
			autho += 1

			if autho == 2 {
				gender := "Парень"
				if user.Gender == "Парень" {
					gender = "Девушка"
				}
				usersGet, err := uc.userService.SearchUser(user.Age-3, user.Age+3, user.City, gender)
				if err != nil {
					log.Println("Лоооооохвхыхвхы", err)
					ctx.Send("Произашла ошибка! попробуй еще раз")
					profileKeys := [][]telebot.ReplyButton{
						{{Text: "1"}, {Text: "2"}, {Text: "3"}},
					}
					return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
				}
				if len(usersGet) > 0 {
					usersLike = append(usersLike, usersGet...)
					Answer := [][]telebot.ReplyButton{
						{{Text: "Начать"}},
					}
					autho = 0
					likes = 1
					return ctx.Send("Смогли подобрать идеальную пару для тебя нажми \"Начать\"", &telebot.ReplyMarkup{ReplyKeyboard: Answer, ResizeKeyboard: true})
				} else {
					return ctx.Send("Не смогли подобрать тебе пару :(")
				}
			}

			if autho == 3 {
				if err != nil {
					ctx.Send("Произашла ошибка! попробуй еще раз")
					profileKeys := [][]telebot.ReplyButton{
						{{Text: "1"}, {Text: "2"}, {Text: "3"}},
					}
					return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
				}
				fmt.Println(user)
				imageBytes, err := utilites.DownloadImageAsBytes(user.Photo)
				if err != nil {
					log.Println(err)
					return ctx.Send("Ошибка загрузки фотографии из бд")
				}
				Answer := &telebot.Photo{
					File:    telebot.FromReader(bytes.NewReader(imageBytes)),
					Caption: fmt.Sprintf("%s, %d, %s - %s", user.Name, user.Age, user.City, user.Description),
				}
				ctx.Send(Answer)
				profileKeys := [][]telebot.ReplyButton{
					{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				}
				return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
			}

			if autho == 4 {
				err := uc.userService.Delete(ctx.Sender().ID)
				if err != nil {
					log.Fatal(err)
				}
				state = 1
				autho = 0
				return ctx.Send("Как тебя зовут?", &telebot.ReplyMarkup{RemoveKeyboard: true})
			}
		}

		if state > 0 {
			user := users[ctx.Sender().ID]
			if state == 1 {
				user.Name = ctx.Text()
				users[ctx.Sender().ID] = user
				state = 2
				return ctx.Send("Теперь укажи свой возраст:")
			}

			if state == 2 {
				age, err := strconv.Atoi(ctx.Text())
				if err != nil {
					return ctx.Send("Возраст должен быть числом!")
				}
				user.Age = age
				users[ctx.Sender().ID] = user
				state = 3
				return ctx.Send("В каком городе ты живешь?")
			}

			if state == 3 {
				user.City = ctx.Text()
				users[ctx.Sender().ID] = user
				state = 4
				genderKeys := [][]telebot.ReplyButton{
					{{Text: "Парень"}, {Text: "Девушка"}},
				}
				return ctx.Send("Выбери свой пол:", &telebot.ReplyMarkup{ReplyKeyboard: genderKeys, ResizeKeyboard: true})
			}

			if state == 4 {
				if ctx.Text() != "Парень" && ctx.Text() != "Девушка" {
					ctx.Send("Такого пола нет!")
					genderKeys := [][]telebot.ReplyButton{
						{{Text: "Парень"}, {Text: "Девушка"}},
					}
					return ctx.Send("Выбери свой пол:", &telebot.ReplyMarkup{ReplyKeyboard: genderKeys, ResizeKeyboard: true})
				}
				user.Gender = ctx.Text()
				users[ctx.Sender().ID] = user
				state = 5
				return ctx.Send("Напиши описание к своей анкете:", &telebot.ReplyMarkup{RemoveKeyboard: true})
			}

			if state == 5 {
				user.Description = ctx.Text()
				state = 6
				users[ctx.Sender().ID] = user
				return ctx.Send("Пришли фото для анкеты:")
			}
		}
		return nil
	})

	b.Handle(telebot.OnPhoto, func(ctx telebot.Context) error {
		if state == 6 {
			photo := ctx.Message().Photo
			if photo == nil {
				return ctx.Send("Ошибка при получении фотографии. Отправь фото еще раз")
			}

			file, err := b.FileByID(photo.FileID)
			if err != nil {
				return ctx.Send("Ошибка при получении файла. Попробуйте еще раз.")
			}

			if file.FilePath == "" {
				return ctx.Send("Ошибка: путь к файлу отсутствует.")
			}

			log.Println("Фото получено: ", file.FilePath)

			filePatch := fmt.Sprintf("./%s", file.FileID)

			fileReader, err := b.File(&file)
			if err != nil {
				return ctx.Send("Ошибка при чтении файла.")
			}

			fileData, err := io.ReadAll(fileReader)
			if err != nil {
				return ctx.Send("Ошибка при чтении файла. Попробуйте еще раз.")
			}

			user := users[ctx.Sender().ID]
			users[ctx.Sender().ID] = user
			state = 0

			err = uc.userService.CreateUser(user.Name, user.City, user.Gender, user.Description, user.Age, ctx.Sender().ID, fileData, filePatch)
			if err != nil {
				log.Println(err)
				return ctx.Send("Ошибка при отправке в базу. Попробуйте еще раз.")
			}

			ctx.Send("Анкета Успешно создана! 🎉")
			time.Sleep(50 * time.Millisecond)

			ctx.Send("Твоя анкета:")
			Answer := &telebot.Photo{
				File:    telebot.FromReader(bytes.NewReader(fileData)),
				Caption: fmt.Sprintf("%s, %d, %s - %s", user.Name, user.Age, user.City, user.Description),
			}
			delete(users, ctx.Sender().ID)
			autho = 1
			ctx.Send(Answer)
			profileKeys := [][]telebot.ReplyButton{
				{{Text: "1"}, {Text: "2"}, {Text: "3"}},
			}
			return ctx.Send("1. Смотреть анкеты 🚀. \n2. Моя анкета 📱.\n3. Изменить анкету.", &telebot.ReplyMarkup{ReplyKeyboard: profileKeys, ResizeKeyboard: true})
		}
		return nil
	})

	log.Println("Бот запущен...")
	b.Start()
}
