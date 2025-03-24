package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) SaveLike(fromUserID, toUserID int64) error {
	query := `INSERT INTO likes(from_user_id, to_user_id) VALUES($1, $2)`
	_, err := r.pool.Exec(context.Background(), query, fromUserID, toUserID)
	return err
}

func (r *Repository) CheckMatch(fromUserID, toUserID int64) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM likes
		WHERE from_user_id = $1 AND to_user_id = $2
	`
	err := r.pool.QueryRow(context.Background(), query, fromUserID, toUserID).Scan(&count)
	return count > 0, err
}
