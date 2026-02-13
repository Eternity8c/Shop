package auth_test

import (
	"auth-service/internal/domens/models"
	"auth-service/internal/services/auth"
	"auth-service/internal/services/auth/mocks"
	"auth-service/internal/storage"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name     string
		email    string
		password string
		provErr  error
	}{
		{
			name:     "Success",
			email:    "test1@gmail.com",
			password: "test1",
		},
		{
			name:     "Empty email",
			email:    "",
			password: "test",
			provErr:  auth.ErrInvalidCredentials,
		},
		{
			name:     "Empty pass",
			email:    "test1",
			password: "",
			provErr:  auth.ErrInvalidCredentials,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

			usrProv := mocks.NewUserProvider(t)
			usrSaver := mocks.NewUserSaver(t)

			if tc.provErr != nil {
				usrProv.On("User", mock.Anything, tc.email).Return(models.User{}, tc.provErr)

				svc := auth.New(logger, usrSaver, usrProv, time.Minute, "test_secret")

				token, err := svc.Login(ctx, tc.email, tc.password)
				require.ErrorIs(t, err, auth.ErrInvalidCredentials)
				require.Empty(t, token)
			} else {
				hash, err := bcrypt.GenerateFromPassword([]byte(tc.password), bcrypt.DefaultCost)
				require.NoError(t, err)

				user := models.User{
					ID:       1,
					Email:    tc.email,
					PassHash: hash,
				}

				usrProv.On("User", mock.Anything, tc.email).Return(user, nil)

				svc := auth.New(logger, usrSaver, usrProv, time.Minute, "test_secret")

				token, err := svc.Login(ctx, tc.email, tc.password)
				require.NoError(t, err)
				require.NotEmpty(t, token)
			}

			usrProv.AssertExpectations(t)
			usrSaver.AssertExpectations(t)
		})
	}
}

func TestRegisterNewUser(t *testing.T) {
	cases := []struct {
		nameCases  string
		email      string
		password   string
		fullName   string
		errSaver   error
		errStorage error
	}{
		{
			nameCases: "Success",
			email:     "test1gmail.com",
			password:  "test",
			fullName:  "test",
		},
		{
			nameCases: "Empty email",
			email:     "",
			password:  "test",
			fullName:  "test",
			errSaver:  errors.New("empty value"),
		},
		{
			nameCases: "Empty password",
			email:     "test1gmail.com",
			password:  "",
			fullName:  "test",
			errSaver:  errors.New("empty password"),
		},
		{
			nameCases: "Empty fullName",
			email:     "test1gmail.com",
			password:  "test",
			fullName:  "",
			errSaver:  errors.New("empty fullName"),
		},
		{
			nameCases:  "User exists",
			email:      "test1gmail.com",
			password:   "test",
			fullName:   "test",
			errStorage: auth.ErrUserExists,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.nameCases, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

			usrSaver := mocks.NewUserSaver(t)

			if tc.errSaver != nil {
				usrSaver.On("SaveUser", mock.Anything, tc.email, mock.MatchedBy(func(b []byte) bool { return len(b) > 0 }), tc.fullName).Return(int64(0), tc.errSaver)

				svc := auth.New(logger, usrSaver, mocks.NewUserProvider(t), time.Minute, "test-secret")

				id, err := svc.RegisterNewUser(ctx, tc.email, tc.password, tc.fullName)
				require.Error(t, err)
				require.EqualValues(t, id, int64(0))
			} else if tc.errStorage != nil {
				usrSaver.On("SaveUser", mock.Anything, tc.email, mock.MatchedBy(func(b []byte) bool { return len(b) > 0 }), tc.fullName).Return(int64(0), auth.ErrUserExists)

				svc := auth.New(logger, usrSaver, mocks.NewUserProvider(t), time.Minute, "test-secret")

				id, err := svc.RegisterNewUser(ctx, tc.email, tc.password, tc.fullName)
				require.Error(t, err)
				require.ErrorIs(t, err, auth.ErrUserExists)
				require.EqualValues(t, id, int64(0))
			} else {
				usrSaver.On("SaveUser", mock.Anything, tc.email, mock.MatchedBy(func(b []byte) bool { return len(b) > 0 }), tc.fullName).Return(int64(1), nil)

				svc := auth.New(logger, usrSaver, mocks.NewUserProvider(t), time.Minute, "test-secret")

				id, err := svc.RegisterNewUser(ctx, tc.email, tc.password, tc.fullName)
				require.NoError(t, err)
				require.NotEmpty(t, id)
			}
			usrSaver.AssertExpectations(t)
		})
	}
}

func TestIsAdmin(t *testing.T) {
	cases := []struct {
		name    string
		id      int64
		isAdmin bool
		err     error
	}{
		{
			name:    "True IsAdmin",
			isAdmin: true,
			id:      1,
		},
		{
			name:    "False IsAdmin",
			isAdmin: false,
			id:      1,
		},
		{
			name:    "User not found",
			isAdmin: false,
			id:      1,
			err:     storage.ErrUserNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

			usrProvider := mocks.NewUserProvider(t)

			if tc.err != nil {
				usrProvider.On("IsAdmin", mock.Anything, tc.id).Return(false, tc.err)

				svc := auth.New(logger, mocks.NewUserSaver(t), usrProvider, time.Minute, "test-sectet")

				isAdmin, err := svc.IsAdmin(ctx, tc.id)
				require.ErrorIs(t, err, storage.ErrUserNotFound)
				require.False(t, isAdmin)
			} else if tc.isAdmin {
				usrProvider.On("IsAdmin", mock.Anything, tc.id).Return(tc.isAdmin, nil)

				svc := auth.New(logger, mocks.NewUserSaver(t), usrProvider, time.Minute, "test-sectet")

				isAdmin, err := svc.IsAdmin(ctx, tc.id)
				require.NoError(t, err)
				require.True(t, isAdmin)
			} else {
				usrProvider.On("IsAdmin", mock.Anything, tc.id).Return(tc.isAdmin, nil)

				svc := auth.New(logger, mocks.NewUserSaver(t), usrProvider, time.Minute, "test-sectet")

				isAdmin, err := svc.IsAdmin(ctx, tc.id)
				require.NoError(t, err)
				require.False(t, isAdmin)
			}
			usrProvider.AssertExpectations(t)
		})
	}
}
