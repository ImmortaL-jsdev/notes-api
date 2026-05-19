package repository

import (
	"context"
	"fmt"

	myerrors "github.com/ImmortaL-jsdev/notes-api/internal/errors"
	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(connString string) (*UserStore, error) {
	pool, err := pgxpool.New(context.Background(), connString)

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &UserStore{pool: pool}, nil
}

func (s *UserStore) Close() {
	s.pool.Close()
}

func (s *UserStore) CreateUser(ctx context.Context, email, passwordHash string) (models.User, error) {
	var created models.User

	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at`

	err := s.pool.QueryRow(ctx, query, email, passwordHash).Scan(&created.ID, &created.Email, &created.PasswordHash, &created.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert user: %w", err)
	}
	return created, nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`

	err := s.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return models.User{}, &myerrors.NotFoundError{Entity: "user", ID: email}
	}
	return user, nil
}
