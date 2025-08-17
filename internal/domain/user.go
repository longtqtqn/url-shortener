package domain

import "time"

type User struct {
	ID        int64
	Email     string
	APIKEY    string
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}
