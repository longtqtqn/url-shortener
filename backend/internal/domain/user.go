package domain

import "time"

type User struct {
	ID        int64
	Email     string
	Password  *string // Optional password for user authentication
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}
