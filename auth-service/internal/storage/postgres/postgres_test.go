package postgres_test

import (
	"auth-service/internal/domens/models"
	"auth-service/internal/storage"
	"auth-service/internal/storage/postgres"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pool *pgxpool.Pool
var container testcontainers.Container
var connStr string

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "postgres:16",
		Env: map[string]string{
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_DB":       "auth_test",
		},
		ExposedPorts: []string{"5432"},
		WaitingFor:   wait.ForListeningPort("5432").WithStartupTimeout(30 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Println("failed to start container:", err)
		os.Exit(1)
	}
	container = c
	defer container.Terminate(ctx)

	host, err := container.Host(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mp, err := container.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	connStr = fmt.Sprintf("postgres://postgres:pass@%s:%s/auth_test?sslmode=disable", host, mp.Port())

	pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Println("pgxpool.New:", err)
		os.Exit(1)
	}
	defer pool.Close()

	_, _ = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
	)`)

	if _, err := pool.Exec(ctx,
		`INSERT INTO users(full_name, email, pass_hash, is_admin)
         VALUES($1, $2, $3, $4)
         ON CONFLICT (email) DO NOTHING`,
		"isAdminTrue", "tc_isAdmin@gmail.com", []byte("admin-pass-hash"), true); err != nil {
		fmt.Println("failed to insert admin user:", err)
		pool.Close()
		_ = container.Terminate(ctx)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestSaveUser_Success(t *testing.T) {
	if pool == nil {
		t.Skip("pg pool not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	id, err := s.SaveUser(ctx, "tc@example.com", []byte("hash"), "TC User")
	require.NoError(t, err)
	require.NotEqualValues(t, id, 0)
}

func TestSaveUser_DuplicationRegister(t *testing.T) {
	if pool == nil {
		t.Skip("pg pool not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	id, err := s.SaveUser(ctx, "tc@example.com", []byte("hash"), "TC User")
	require.ErrorIs(t, err, storage.ErrUserExists)
	require.EqualValues(t, id, 0)
}

func TestLogin_Success(t *testing.T) {
	if pool == nil {
		t.Skip("pg pool not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	var u models.User
	u, err = s.User(ctx, "tc@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, u)
}

func TestLogin_Failed(t *testing.T) {
	if pool == nil {
		t.Skip("pg pool not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	var u models.User
	u, err = s.User(ctx, "tc111111@example.com")
	require.ErrorIs(t, err, storage.ErrUserNotFound)
	require.Empty(t, u)
}

func TestIsAdmin_Success(t *testing.T) {
	if pool == nil {
		t.Skip("pg pool not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	u, err := s.IsAdmin(ctx, 1)
	require.NoError(t, err)
	require.True(t, u)
}

func TestIsAdmin_Failed(t *testing.T) {
	if pool == nil {
		t.Skip("pg pool not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	u, err := s.IsAdmin(ctx, 2)
	require.NoError(t, err)
	require.False(t, u)
}
func TestIsAdmin_UnknownUser(t *testing.T) {
	if pool == nil {
		t.Skip("pg pool not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s, err := postgres.New(ctx, connStr)
	require.NoError(t, err)

	u, err := s.IsAdmin(ctx, 1222)
	require.ErrorIs(t, err, storage.ErrUserNotFound)
	require.False(t, u)
}
