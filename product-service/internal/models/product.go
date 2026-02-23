package models

import "time"

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       int64
	Stock       int64
	Created_at  time.Duration
}
