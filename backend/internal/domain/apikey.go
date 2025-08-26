package domain

import "time"

type ApiKey struct {
	ID        int64
	UserID    int64
	Key       string
	DeletedAt *time.Time
	CreatedAt time.Time
}
