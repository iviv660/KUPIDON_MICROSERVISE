package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"service1/internal/entity"
	"service1/internal/storage"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	UpdateUser(ctx context.Context, id int, user *entity.User) error
	DeleteUser(ctx context.Context, id int64) error
	SearchUsers(ctx context.Context, filter entity.UserFilter) ([]entity.User, error)
}

type UserUsecase struct {
	repo         UserRepository
	fileStorage  storage.FileStorage
	redisStorage storage.RedisStorage
}

func NewUserUsecase(repo UserRepository, fileStorage storage.FileStorage, redisStorage storage.RedisStorage) *UserUsecase {
	if repo == nil {
		panic("UserRepository cannot be nil")
	}
	if fileStorage == nil {
		panic("FileStorage cannot be nil")
	}
	if redisStorage == nil {
		panic("RedisStorage cannot be nil")
	}

	return &UserUsecase{repo: repo, fileStorage: fileStorage, redisStorage: redisStorage}
}

func (u *UserUsecase) Create(ctx context.Context, name, description, fileName, gender, city string, age int, file multipart.File, filesize, telegramId int64) (int, error) {
	if name == "" {
		return 0, errors.New("name is required")
	}
	if age <= 0 {
		return 0, errors.New("age is required")
	}
	if description == "" {
		return 0, errors.New("description is required")
	}
	url, err := u.fileStorage.UploadFile(ctx, file, fileName, filesize)
	if err != nil {
		return 0, err
	}
	user := &entity.User{
		Name:        name,
		Age:         age,
		Description: description,
		Photo:       url,
		TelegramID:  telegramId,
		Gender:      gender,
		City:        city,
	}
	return u.repo.CreateUser(ctx, user)
}

func (u *UserUsecase) Search(ctx context.Context, filter entity.UserFilter) ([]entity.User, error) {
	cacheKey := fmt.Sprintf("search:%+v", filter)

	cachedData, err := u.redisStorage.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var users []entity.User
		if err := json.Unmarshal([]byte(cachedData), &users); err != nil {
			return users, nil
		}
	}

	users, err := u.repo.SearchUsers(ctx, filter)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	u.redisStorage.Set(ctx, cacheKey, data, 24*time.Hour)

	return users, nil
}

func (u *UserUsecase) GetByID(ctx context.Context, telegram_id int64) (*entity.User, error) {
	if telegram_id <= 0 {
		return nil, errors.New("invalid id")
	}
	cacheKey := fmt.Sprintf("user:%d", telegram_id)
	cachedUser, err := u.redisStorage.Get(ctx, cacheKey).Result()
	if err == nil && cachedUser != "" {
		user := &entity.User{}
		if err := json.Unmarshal([]byte(cachedUser), user); err != nil {
			return nil, err
		}
	}

	user, err := u.repo.GetUserByID(ctx, telegram_id)
	if err != nil {
		return nil, err
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	_ = u.redisStorage.Set(ctx, cacheKey, userJSON, 24*time.Hour)

	return user, nil
}

func (u *UserUsecase) Update(ctx context.Context, name, description, fileName, gender, city string, age, id int, file multipart.File, filesize, telegramId int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if age <= 0 {
		return errors.New("age is required")
	}
	if description == "" {
		return errors.New("description is required")
	}
	if city == "" {
		return errors.New("city is required")
	}
	if gender == "" {
		return errors.New("gender is required")
	}

	url, err := u.fileStorage.UploadFile(ctx, file, fileName, filesize)
	if err != nil {
		return err
	}
	user := &entity.User{
		Name:        name,
		Age:         age,
		Description: description,
		Photo:       url,
		TelegramID:  telegramId,
		Gender:      gender,
		City:        city,
	}

	cacheKey := fmt.Sprintf("user:%d", id)
	_ = u.redisStorage.Del(ctx, cacheKey)

	return u.repo.UpdateUser(ctx, id, user)
}

func (u *UserUsecase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user:%d", id)
	_ = u.redisStorage.Del(ctx, cacheKey)

	return nil
}
