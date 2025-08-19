package domain

import "time"

type User struct {
	ID            int64
	Email         string
	DeletedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	Plan          string
	Role          string
	PlanExpiresAt *time.Time
}
