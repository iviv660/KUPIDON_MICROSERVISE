# Телеграм-бот бот для знакомств "KupiDon"

## Содержание
- [Краткое описание](#краткое-описание)
- [Стек технологий](#стек-технологий)
- [Структура проекта](#структура-проекта)
- [Установка и настройка](#установка-и-настройка)

## Краткое-описание
Телеграм-бот для знакомств "KupiDon" позволяет пользователям находить интересных собеседников, обмениваться лайками и начинать общение.
Основные возможности:
- 📌 Регистрация и создание профиля
- 🔍 Подбор потенциальных пар 
- ❤️ Лайки и дизлайки для поиска совпадени
- 🔔 Уведомления о новых лайках и мэтчах 

Бот: @ExchangeRateRussian111_bot

Команда для запуска бота /start

## Стек технологий
- Язык программирования: Go
- База данных: PostgreSQL
- База данных для храннения фото: S3 (minio)
- Кэш: Redis
- API: REST 
- Очереди сообщений: Kafka
- Docker 

## Структура проекта

![Скриншот 23-03-2025 222751](https://github.com/user-attachments/assets/f10a2311-3578-49dc-8af1-88ff66bb53c4)

## Установка и настройка
Эти инструкции помогут развернуть проект на локальном компьютере для использования.

### Перед тем, как начать
- Проверьте, что у вас установленa актуальная версия GO

### Клонирование репозитория
- Перейдите в папку проекта и выполните команду git clone, чтобы скопировать файлы репозитория

```bash
git clone  https://github.com/iviv660/KUPIDON_MICROSERVISE. git 
```

### Создание бота в телеграм (BotFather)
- Найдите бота BotFather в телеграм и создайте нового бота командой /newbot
- Следуя инструкции, укажите название и username бота. Скопируйте отправленный вам токен бота.

### Docker-compose
- Откройте ServiceBot в файле Docker-compose вставте в TELEGRAM_BOT_TOKEN = "Ваше токен"
- Откройте ServiceNotification в файле Docker-compose вставте в TELEGRAM_BOT_TOKEN = "Ваше токен"

### Запуск бота
- Создайте образы каждого Dokecrfile:
  ```bash
   docker build -t service-bot:latest .
   docker build -t service-user:latest .
   docker build -t service-match:latest .
   docker build -t service-notification:latest .
  

- Запустите все docker-compose (для каждого сервиса)
  ```bash
	docker-compose up --build 

Команда для запуска бота /start
![Скриншот 23-03-2025 230429](https://github.com/user-attachments/assets/afeaf510-13b0-4679-8ad2-ca943f682c97)

## Используемые библиотеки

### Для работы с Telegram:
- **github.com/go-telegram-bot-api/telegram-bot-api** — библиотека для взаимодействия с Telegram Bot API.

### Для работы с PostgreSQL:
- **github.com/jackc/pgxpool/v4** — PostgreSQL клиент для Go, используется для работы с базой данных.

### Для работы с Kafka:
- **github.com/segmentio/kafka-go** — библиотека для работы с Kafka.

### Для работы с Redis:
- **github.com/go-redis/redis/v8** — библиотека для работы с Redis.

### Для загрузки и работы с изображениями:
- **github.com/minio/minio-go/v7** — клиент для работы с MinIO, используется для хранения изображений в S3 совместимом хранилище.

### Для тестирования:
- **github.com/stretchr/testify** — библиотека для упрощения юнит-тестирования.
