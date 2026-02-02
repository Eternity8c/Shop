package auth

import "context"

type Auth interface {
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string) (token string, err error)
	Validate(ctx context.Context, userID int64) (bool, error)
}
