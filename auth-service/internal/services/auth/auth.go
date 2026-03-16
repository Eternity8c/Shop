package auth

import (
	"auth-service/internal/domens/models"
	"auth-service/internal/lib/jwt"
	"auth-service/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
	secret      string
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserSaver
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte, full_name string) (int64, error)
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserProvider
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, uid int64) (bool, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProcider UserProvider,
	tokenTTL time.Duration,
	secret string,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProcider,
		tokenTTL:    tokenTTL,
		secret:      secret,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "Auth.Login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("useremail", email),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", "error", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", "error", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	isAdmin, err := a.usrProvider.IsAdmin(ctx, user.ID)
	if err != nil {
		a.log.Warn("failed to check is_admin, defaulting to false", "error", err)
		isAdmin = false
	}
	user.IsAdmin = isAdmin

	a.log.Info("user logged is successfully")

	token, err := jwt.NewToken(user, a.tokenTTL, a.secret)
	if err != nil {
		a.log.Error("failed to generation token", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context,
	email string,
	password string,
	fullName string,
) (int64, error) {
	const op = "Services.Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "error", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash, fullName)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", "error", err)

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", "error", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "Services.Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", "error", err)

			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("chacked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
