package authHandlers_test

import (
	authHandlers "api-geteway/internal/handlers/auth"
	"api-geteway/internal/handlers/auth/mocks"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http/httptest"
	authproto "shop/auth-service/api/gen/go/api/proto"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name     string
		email    string
		password string
		err      error
		wantCode int
		wantBody string
	}{
		{
			name:     "Success",
			email:    "test1@gmail.com",
			password: "test",
			wantBody: `{"token":"test"}`,
			wantCode: 200,
		},
		{
			name:     "Empty email",
			email:    "",
			password: "test",
			err:      status.Error(codes.InvalidArgument, "invalid email or password"),
			wantCode: 400,
			wantBody: "invalid email or password\n",
		},
		{
			name:     "Empty password",
			email:    "test1@gmail.com",
			password: "",
			err:      status.Error(codes.InvalidArgument, "invalid email or password"),
			wantCode: 400,
			wantBody: "invalid email or password\n",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ua := mocks.NewUserAuth(t)

			if tc.err != nil {
				ua.On("Login", mock.Anything, tc.email, tc.password).
					Return(nil, tc.err)
			} else {
				ua.On("Login", mock.Anything, tc.email, tc.password).
					Return(&authproto.LoginResponse{
						Token: "test",
					}, nil)
			}

			ah := authHandlers.NewAuthHandler(ua, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))
			body := `{ "email":"` + tc.email + `","password":"` + tc.password + `"}`
			r := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte(body)))
			w := httptest.NewRecorder()
			ah.Login(w, r)

			require.Equal(t, tc.wantCode, w.Code)
			if tc.wantBody != "" {
				require.Equal(t, tc.wantBody, w.Body.String())
			}

			ua.AssertExpectations(t)
		})
	}
}

func TestRegister(t *testing.T) {
	cases := []struct {
		nameCs   string
		email    string
		password string
		fullName string
		err      error
		wantCode int
		wantBody string
	}{
		{
			nameCs:   "Success",
			email:    "test@gmail.com",
			password: "test",
			fullName: "test",
			wantCode: 201,
		},
		{
			nameCs:   "Dublication Register",
			email:    "test@gmail.com",
			password: "test",
			fullName: "test",
			err:      status.Error(codes.AlreadyExists, "user already exists"),
			wantCode: 409,
			wantBody: "user already exists\n",
		},
		{
			nameCs:   "Empty email",
			email:    "",
			password: "test",
			fullName: "test",
			err:      errors.New("empty email"),
			wantCode: 500,
			wantBody: "internal server error\n",
		},
		{
			nameCs:   "Empty password",
			email:    "test@gmail.com",
			password: "",
			fullName: "test",
			err:      errors.New("empty password"),
			wantCode: 500,
			wantBody: "internal server error\n",
		},
		{
			nameCs:   "Empty fullName",
			email:    "test@gmail.com",
			password: "test",
			fullName: "",
			err:      errors.New("empty fullName"),
			wantCode: 500,
			wantBody: "internal server error\n",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.nameCs, func(t *testing.T) {
			t.Parallel()

			ua := mocks.NewUserAuth(t)

			if tc.err != nil {
				ua.On("Register", mock.Anything, tc.email, tc.password, tc.fullName).
					Return(nil, tc.err)
			} else {
				ua.On("Register", mock.Anything, tc.email, tc.password, tc.fullName).
					Return(&authproto.RegisterResponse{
						UserId: 1,
					}, nil)
			}

			ah := authHandlers.NewAuthHandler(ua, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})))
			body := `{ "email":"` + tc.email + `","password":"` + tc.password + `","full_name":"` + tc.fullName + `"}`
			r := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte(body)))
			w := httptest.NewRecorder()
			ah.Register(w, r)

			require.EqualValues(t, tc.wantCode, w.Code)
			if tc.wantBody != "" {
				require.EqualValues(t, tc.wantBody, w.Body.String())
			}

			ua.AssertExpectations(t)
		})
	}
}
