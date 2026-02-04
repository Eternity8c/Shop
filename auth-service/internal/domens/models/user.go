package models

type User struct {
	ID        int64
	Email     string
	PassHash  []byte
	Full_Name string
}
