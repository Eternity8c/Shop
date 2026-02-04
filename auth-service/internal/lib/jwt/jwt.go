package jwt

import (
	"auth-service/internal/domens/models"
	"time"

	JWT "github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, duration time.Duration, secret string) (string, error) {
	token := JWT.New(JWT.SigningMethodHS256)

	claims := token.Claims.(JWT.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
