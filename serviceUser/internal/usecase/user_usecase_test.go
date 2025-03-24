package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"service1/internal/entity"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, id int, user *entity.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *MockRepository) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) SearchUsers(ctx context.Context, filter entity.UserFilter) ([]entity.User, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entity.User), args.Error(1)
}

type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) UploadFile(ctx context.Context, file multipart.File, fileName string, filesize int64) (string, error) {
	args := m.Called(ctx, file, fileName, filesize)
	return args.String(0), args.Error(1)
}

type MockRedisStorage struct {
	mock.Mock
}

func (m *MockRedisStorage) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisStorage) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisStorage) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func TestUserUsecase_Create(t *testing.T) {
	repo := new(MockRepository)
	fileStorage := new(MockFileStorage)
	redisStorage := new(MockRedisStorage)
	usecase := NewUserUsecase(repo, fileStorage, redisStorage)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		name := "test name"
		description := "test description"
		fileName := "photo.jpg"
		telegramID := int64(321312312)
		age := 25
		file := new(multipart.File)
		fileSize := int64(1024)
		gender := "men"
		city := "moscow"

		fileStorage.On("UploadFile", ctx, file, fileName, fileSize).Return("http://example.com/photo.jpg", nil)
		repo.On("CreateUser", ctx, mock.AnythingOfType("*entity.User")).Return(1, nil)

		userID, err := usecase.Create(ctx, name, description, fileName, gender, city, age, *file, fileSize, telegramID)

		assert.NoError(t, err)
		assert.Equal(t, 1, userID)

		fileStorage.AssertExpectations(t)
		repo.AssertExpectations(t)
		redisStorage.AssertExpectations(t)
	})

	t.Run("Fail on UploadFIle", func(t *testing.T) {
		ctx := context.Background()

		name := "test name"
		description := "test description"
		fileName := "photo.jpg"
		age := 25
		file := new(multipart.File)
		fileSize := int64(1024)
		telegramID := int64(321312312)
		gender := "men"
		city := "moscow"
		fileStorage.On("UploadFile", ctx, file, fileName, fileSize).Return("", errors.New("upload error"))
		userID, err := usecase.Create(ctx, name, description, fileName, gender, city, age, *file, fileSize, telegramID)

		// Проверяем, что произошла ошибка
		assert.Error(t, err)
		assert.Equal(t, 0, userID)

		fileStorage.AssertExpectations(t)

	})
}

func TestUserUsecase_GetByID(t *testing.T) {
	repo := new(MockRepository)
	fileStorage := new(MockFileStorage)
	redisStorage := new(MockRedisStorage)
	usecase := NewUserUsecase(repo, fileStorage, redisStorage)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		id := int64(1)
		expectedUser := entity.User{
			ID:          int(id),
			Name:        "test name",
			Age:         25,
			Description: "test description",
			Photo:       "http://example.com/photo.jpg",
		}

		redisStorage.On("Get", ctx, fmt.Sprintf("user:%d", id)).Return(redis.NewStringResult("", redis.Nil))
		redisStorage.On("Set", ctx, fmt.Sprintf("user:%d", id), mock.Anything, time.Minute*10).Return(nil)

		repo.On("GetUserByID", ctx, id).Return(&expectedUser, nil)

		user, err := usecase.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, *user)

		repo.AssertExpectations(t)
		redisStorage.AssertExpectations(t)
	})

	t.Run("Fail on Get", func(t *testing.T) {
		ctx := context.Background()
		id := 1
		expectedUser := entity.User{
			ID:          id,
			Name:        "test name",
			Age:         25,
			Description: "test description",
			Photo:       "http://example.com/photo.jpg",
		}

		redisStorage.On("Get", ctx, fmt.Sprintf("user:%d", id)).Return(redis.NewStringResult("", redis.ErrClosed))
		redisStorage.On("Set", ctx, fmt.Sprintf("user:%d", id), mock.Anything, time.Minute*10).Return(nil)

		repo.On("GetUserByID", ctx, id).Return(&expectedUser, nil)

		user, err := usecase.GetByID(ctx, int64(id))

		assert.Error(t, err)
		assert.Equal(t, "", *user)

		repo.AssertExpectations(t)
		redisStorage.AssertExpectations(t)
	})
}

func TestUserUsecase_Update(t *testing.T) {
	repo := new(MockRepository)
	fileStorage := new(MockFileStorage)
	redisStorage := new(MockRedisStorage)
	usecase := NewUserUsecase(repo, fileStorage, redisStorage)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		id := 1
		age := 25
		name := "test name"
		description := "test description"
		filename := "photo.jpg"
		file := new(multipart.File)
		filesize := int64(1024)
		telegramID := int64(321312312)
		gender := "men"
		city := "moscow"

		fileStorage.On("UploadFile", ctx, file, filename, filesize).Return("http://example.com/photo.jpg", nil)
		redisStorage.On("Del", fmt.Sprintf("user:%d", id)).Return(nil)
		repo.On("UpdateUser", ctx, id, mock.AnythingOfType("*entity.User")).Return(nil)

		err := usecase.Update(ctx, name, description, filename, gender, city, age, id, *file, filesize, telegramID)

		assert.NoError(t, err)

		fileStorage.AssertExpectations(t)
		redisStorage.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Fail on Upload", func(t *testing.T) {
		ctx := context.Background()

		id := 1
		age := 25
		name := "test name"
		description := "test description"
		filename := "photo.jpg"
		file := new(multipart.File)
		filesize := int64(1024)
		telegramID := int64(321312312)
		gender := "men"
		city := "moscow"

		fileStorage.On("UploadFile", ctx, file, filename, filesize).Return("http://example.com/photo.jpg", nil)

		err := usecase.Update(ctx, name, description, filename, gender, city, age, id, *file, filesize, telegramID)

		assert.Error(t, err)

		fileStorage.AssertExpectations(t)
	})
}

func TestUserUsecase_Delete(t *testing.T) {
	repo := new(MockRepository)
	fileStorage := new(MockFileStorage)
	redisStorage := new(MockRedisStorage)
	usecase := NewUserUsecase(repo, fileStorage, redisStorage)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		id := 1

		repo.On("DeleteUser", ctx, id).Return(nil)
		redisStorage.On("Del", ctx, fmt.Sprintf("user:%d", id))

		err := usecase.Delete(ctx, int64(id))

		assert.NoError(t, err)

		repo.AssertExpectations(t)
		redisStorage.AssertExpectations(t)
	})

	t.Run("Fail on DeleteUser", func(t *testing.T) {
		ctx := context.Background()

		id := 1

		repo.On("DeleteUser", ctx, id).Return(errors.New("Fail DeleteUser"))

		err := usecase.Delete(ctx, int64(id))

		assert.Error(t, err)

		repo.AssertExpectations(t)
	})
}

func TestUserUsecase_Search(t *testing.T) {
	repo := new(MockRepository)
	fileStorage := new(MockFileStorage)
	redisStorage := new(MockRedisStorage)
	usecase := NewUserUsecase(repo, fileStorage, redisStorage)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		repo.On("SearchUsers", ctx, entity.UserFilter{}).Return([]entity.User{}, nil)
		redisStorage.On("Del", fmt.Sprintf("search:%+v", entity.UserFilter{})).Return(nil)
		redisStorage.On("Set", ctx, fmt.Sprintf("search:%+v", entity.UserFilter{}), mock.Anything, time.Minute*10).Return(nil)

		users, err := usecase.Search(ctx, entity.UserFilter{})

		assert.NoError(t, err)
		assert.Equal(t, nil, users)

		repo.AssertExpectations(t)
		redisStorage.AssertExpectations(t)
	})
	t.Run("Error", func(t *testing.T) {
		ctx := context.Background()

		expectedErr := errors.New("search failed")
		repo.On("SearchUsers", ctx, entity.UserFilter{}).Return(nil, expectedErr)

		users, err := usecase.Search(ctx, entity.UserFilter{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, users)

		repo.AssertExpectations(t)
	})
}
