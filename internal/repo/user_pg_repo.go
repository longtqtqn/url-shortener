package repo

import (
	"context"
	"url-shortener/internal/domain"
	"url-shortener/internal/usecase"

	"github.com/uptrace/bun"
)

type UserPGRepository struct {
	db *bun.DB
}

func NewUserPGRepository(db *bun.DB) usecase.UserRepository {
	return &UserPGRepository{db: db}
}

// Create implements usecase.UserRepository.
func (r *UserPGRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

// GetByAPIKey implements usecase.UserRepository.
func (r *UserPGRepository) GetByAPIKey(ctx context.Context, apiKey string) (*domain.User, error) {
	user := new(domain.User)
	err := r.db.NewSelect().Model(user).Where("apikey = ?", apiKey).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
