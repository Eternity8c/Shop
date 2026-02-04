package models

type User struct {
	ID       int64
	Email    string
	PassHash []byte
	FullName string `db:"full_name"`
}
