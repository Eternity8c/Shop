package postgres

import (
	"auth-service/internal/domens/models"
	"auth-service/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "Storage.Postgres.New"

	pool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte, full_name string) (int64, error) {
	const op = "Storage.Postgres.SaveUser"

	stmt := `INSERT INTO users(full_name, email, pass_hash)
	VALUES($1, $2, $3)
	RETURNING id`

	var id int64

	err := s.pool.QueryRow(ctx, stmt, full_name, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "Storage.Postgres.User"

	stmt := `SELECT id, full_name, email, pass_hash
	FROM users
	WHERE email = $1`

	var user models.User
	slog.Info(email)
	err := s.pool.QueryRow(ctx, stmt, email).Scan(&user.ID, &user.FullName, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, uid int64) (bool, error) {
	const op = "Storage.Postgres.IsAdmin"

	stmt := `SELECT is_admin FROM users WHERE id = $1`

	var isAdmin bool

	err := s.pool.QueryRow(ctx, stmt, uid).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
