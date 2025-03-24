package repository

import (
	"context"
	"fmt"
	"service1/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	Pool   *pgxpool.Pool
	Logger *logrus.Logger
}

func NewUserRepository(pool *pgxpool.Pool, logger *logrus.Logger) *UserRepository {
	return &UserRepository{Pool: pool, Logger: logger}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	query := `INSERT INTO users (name, age, description, photo, telegram_id, city, gender) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int

	r.Logger.WithFields(logrus.Fields{
		"name":        user.Name,
		"age":         user.Age,
		"description": user.Description,
		"photo":       user.Photo,
		"telegram_id": user.TelegramID,
		"city":        user.City,
		"gender":      user.Gender,
	}).Info("Executing CreateUser query")

	err := r.Pool.QueryRow(ctx, query, user.Name, user.Age, user.Description, user.Photo, user.TelegramID, user.City, user.Gender).Scan(&id)
	if err != nil {
		r.Logger.WithFields(logrus.Fields{
			"user": user.Name,
		}).Error("Error creating user: ", err)
		return 0, err
	}

	r.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Info("User created successfully")
	return id, nil
}

func (r *UserRepository) SearchUsers(ctx context.Context, filter entity.UserFilter) ([]entity.User, error) {
	query := `SELECT id, telegram_id, name, age, city, gender, description, photo FROM users WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.MinAge != nil {
		query += fmt.Sprintf(" AND age >= $%d", argIndex)
		args = append(args, *filter.MinAge)
		argIndex++
	}
	if filter.MaxAge != nil {
		query += fmt.Sprintf(" AND age <= $%d", argIndex)
		args = append(args, *filter.MaxAge)
		argIndex++
	}
	if filter.City != "" {
		query += fmt.Sprintf(" AND city = $%d", argIndex)
		args = append(args, filter.City)
		argIndex++
	}
	if filter.Gender != "" {
		query += fmt.Sprintf(" AND gender = $%d", argIndex)
		args = append(args, filter.Gender)
		argIndex++
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.TelegramID, &user.Name, &user.Age, &user.City, &user.Gender, &user.Description, &user.Photo); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, telegram_id int64) (*entity.User, error) {
	query := `SELECT id, name, age, description, photo, telegram_id, city, gender FROM users WHERE telegram_id = $1`
	user := &entity.User{}

	r.Logger.WithFields(logrus.Fields{
		"user_telegram_ID": telegram_id,
	}).Info("Executing GetUserByID query")

	err := r.Pool.QueryRow(ctx, query, telegram_id).Scan(&user.ID, &user.Name, &user.Age, &user.Description, &user.Photo, &user.TelegramID, &user.City, &user.Gender)
	if err != nil {
		r.Logger.WithFields(logrus.Fields{
			"user_telegram_ID": telegram_id,
		}).Error("Error getting user by ID: ", err)
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"user_telegram_ID": telegram_id,
	}).Info("User retrieved successfully")
	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, user *entity.User) error {
	query := `UPDATE users SET name = $1, age = $2, description = $3, photo = $4, telegram_id = $5, city = $6, gender = $7 WHERE id = $8`

	r.Logger.WithFields(logrus.Fields{
		"userID":      id,
		"name":        user.Name,
		"age":         user.Age,
		"description": user.Description,
		"photo":       user.Photo,
		"telegram_id": user.TelegramID,
		"city":        user.City,
		"gender":      user.Gender,
	}).Info("Executing UpdateUser query")

	_, err := r.Pool.Exec(ctx, query, user.Name, user.Age, user.Description, user.Photo, user.TelegramID, user.City, user.Gender, id)
	if err != nil {
		r.Logger.WithFields(logrus.Fields{
			"userID": id,
		}).Error("Error updating user: ", err)
		return err
	}

	r.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Info("User updated successfully")
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE telegram_id = $1`

	r.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Info("Executing DeleteUser query")

	_, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		r.Logger.WithFields(logrus.Fields{
			"userID": id,
		}).Error("Error deleting user: ", err)
		return err
	}

	r.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Info("User deleted successfully")
	return nil
}
