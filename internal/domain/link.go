package domain

import (
	"time"
)

type Link struct {
	ID            int64
	UserID        int64
	ShortCode     string
	LongURL       string
	ClickCount    int64
	LastClickedAt *time.Time
	DeletedAt     *time.Time
	CreatedAt     time.Time
}
