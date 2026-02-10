package jwt

import (
	"auth-service/internal/domens/models"
	"fmt"
	"testing"
	"time"

	JWT "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const (
	secret = "test-secret"
)

func TestNewToken_Success(t *testing.T) {
	user := models.User{
		ID:       1,
		Email:    "test",
		PassHash: []byte("fsdf"),
		FullName: "test",
	}

	dur := time.Minute
	token, err := NewToken(user, dur, secret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	parsed, err := JWT.Parse(token, func(t *JWT.Token) (any, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(JWT.MapClaims)
	require.True(t, ok)

	require.Equal(t, fmt.Sprintf("%v", user.ID), fmt.Sprintf("%v", claims["uid"]))
	require.Equal(t, user.Email, claims["email"])
}

func TestNewToken_Expired(t *testing.T) {
	user := models.User{
		ID:       1,
		Email:    "test",
		PassHash: []byte("fsdf"),
		FullName: "test",
	}

	secret := "test-secret"
	// отрицательная длительность — токен сразу просрочен
	tokenStr, err := NewToken(user, -time.Minute, secret)
	require.NoError(t, err)

	_, err = JWT.Parse(tokenStr, func(t *JWT.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.Error(t, err)
}

func TestNewToken_WrongSecret(t *testing.T) {
	user := models.User{
		ID:       1,
		Email:    "test",
		PassHash: []byte("fsdf"),
		FullName: "test",
	}

	token, err := NewToken(user, time.Minute, secret)
	require.NoError(t, err)

	_, err = JWT.Parse(token, func(t *JWT.Token) (any, error) {
		return []byte("secret2"), nil
	})
	require.Error(t, err)
}
