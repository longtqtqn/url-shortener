package domain

import (
	"context"
)

// LinkRepository defines the interface for link data operations
type LinkRepository interface {
	CreateLink(ctx context.Context, link *Link) error

	GetLinkByShortCode(ctx context.Context, shortCode string) (*Link, error)
	GetListLinksByAPIKeyID(ctx context.Context, apiKeyID int64) ([]*Link, error)
	GetLinkNoAPIKey(ctx context.Context, shortCode string) ([]*Link, error)

	SoftDeleteByShortCode(ctx context.Context, apiKeyID int64, shortCode string) error
	TrackClick(ctx context.Context, shortCode string) error
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	CreateAPIKey(ctx context.Context, userID int64, key string) error

	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetAPIKeysByUserID(ctx context.Context, userID int64, limit, offset int) ([]*ApiKey, error)

	GetAPIKeyIDByAPIKey(ctx context.Context, apiKey string) (int64, error)

	SoftDeleteByID(ctx context.Context, userID int64) error
}
