package repo

import (
	"context"
	"url-shortener/internal/domain"
	"url-shortener/internal/repo/model"
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
	userModel := model.ToUserBunModel(user)
	_, err := r.db.NewInsert().Model(userModel).Exec(ctx)
	return err
}

// FindByAPIKey implements usecase.UserRepository.
func (r *UserPGRepository) FindByAPIKey(ctx context.Context, apiKey string) (*domain.User, error) {
	userModel := new(model.UserBunModel)
	err := r.db.NewSelect().Model(userModel).Where("apikey = ?", apiKey).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return userModel.ToDomain(), nil
}

func (r *UserPGRepository) SoftDeleteByID(ctx context.Context, userID int64) error {
	_, err := r.db.NewDelete().
		Model((*model.UserBunModel)(nil)).
		Where("id = ?", userID).
		Exec(ctx)
	return err
}
