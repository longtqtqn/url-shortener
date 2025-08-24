package domain

import (
	"time"
)

type Link struct {
	ID            int64
	APIKeyID      *int64 // Optional: 0 or 1 apikey per link
	ShortCode     string
	LongURL       string
	Password      *string // Optional password for link protection
	ClickCount    int64
	LastClickedAt *time.Time
	DeletedAt     *time.Time
	CreatedAt     time.Time
}
